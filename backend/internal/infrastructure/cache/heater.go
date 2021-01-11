package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	staninfrastructure "social-network/internal/infrastructure/stan"
	"time"
)

const (
	maxInFlightMsg = 100
	friendsDBIndex = 1
	NewsDBIndex    = 2
)

type Heater struct {
	redisPool    *Pool
	mysqlConn    *sql.DB
	stanClient   *staninfrastructure.Client
	subscription stan.Subscription
	logger       *zap.Logger
	doneCh       chan struct{}
}

func NewHeater(redisPool *Pool, mysqlConn *sql.DB, stanClient *staninfrastructure.Client, logger *zap.Logger) *Heater {
	return &Heater{
		redisPool:  redisPool,
		mysqlConn:  mysqlConn,
		stanClient: stanClient,
		logger:     logger,
		doneCh:     make(chan struct{}),
	}
}

func (h *Heater) Start() (err error) {
	ctx := context.Background()
	tx, err := h.mysqlConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = h.heatFriends(ctx, tx); err != nil {
		tx.Rollback()

		return err
	}

	if err = h.heatNews(ctx, tx); err != nil {
		tx.Rollback()

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	subscriptionOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.AckWait(time.Second * 1),
		stan.MaxInflight(maxInFlightMsg),
	}

	h.subscription, err = h.stanClient.Subscribe("cache-reload", h.makeCacheReloadSub(),
		append(subscriptionOptions, stan.DurableName("cache-heater"))...)
	if err != nil {
		return err
	}

	<-h.doneCh
	close(h.doneCh)

	h.logger.Info("cache heater stan subscribing is over")

	return nil
}

func (h *Heater) heatFriends(ctx context.Context, tx *sql.Tx) error {
	friends := make(map[string][]string, 1)

	rows, err := tx.Query(`
		SELECT
			master_user_id, slave_user_id
		FROM
			friendship
		WHERE
		    friendship.status = ?`, "accepted")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id1, id2 string

		if err = rows.Scan(&id1, &id2); err != nil {
			return err
		}

		if _, ok := friends[id1]; !ok {
			friends[id1] = make([]string, 0, 1)
		}
		friends[id1] = append(friends[id1], id2)

		if _, ok := friends[id2]; !ok {
			friends[id2] = make([]string, 0, 1)
		}
		friends[id2] = append(friends[id2], id1)
	}

	conn, err := h.redisPool.GetConnByIndexDB(friendsDBIndex)
	if err != nil {
		tx.Rollback()

		return err
	}

	for key, value := range friends {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("cannot marshal friends id, %w", err)
		}

		if err = conn.Set(ctx, key, data, 0).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (h *Heater) heatNews(ctx context.Context, tx *sql.Tx) error {
	_, err := h.redisPool.GetConnByIndexDB(NewsDBIndex)
	if err != nil {
		return err
	}

	// TODO
	//rows, err := tx.Query(`
	//	SELECT
	//		news.id, user.name, user.surname, user.sex, content, news.create_time
	//	FROM
	//		news
	//	JOIN user ON news.owner_id = user.id`)
	//if err != nil {
	//	return err
	//}
	//defer rows.Close()
	//
	//var n []*domain.News
	//result, err := conn.Get(ctx, userID).Result()
	//switch err {
	//case nil:
	//	if err = json.Unmarshal([]byte(result), &n); err != nil {
	//		return fmt.Errorf("cannot unmarshal news, %w", err)
	//	}
	//case redis.Nil:
	//	n = make([]*domain.News, 0, 1)
	//default:
	//	return err
	//}
	//n = append(n, news...)
	//
	//data, err := json.Marshal(n)
	//if err != nil {
	//	return fmt.Errorf("cannot marshal news, %w", err)
	//}
	//
	//return conn.Set(ctx, userID, data, 0).Err()

	return nil
}

func (h *Heater) Stop() {
	if err := h.subscription.Close(); err != nil {
		h.logger.Error("cache heater can't close stan connection: ", zap.Error(err))
	}
}

func (h *Heater) makeCacheReloadSub() func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			h.logger.Error("cache heater fail to acknowledge a message: ", zap.Error(err))
		}

		ctx := context.Background()
		tx, err := h.mysqlConn.BeginTx(ctx, nil)
		if err != nil {
			h.logger.Error("cache heater fail to get transaction", zap.Error(err))
		}

		if err = h.heatFriends(ctx, tx); err != nil {
			tx.Rollback()

			h.logger.Error("cache heater fail to heat friends", zap.Error(err))
		}

		if err = h.heatNews(ctx, tx); err != nil {
			tx.Rollback()

			h.logger.Error("cache heater fail to heat news", zap.Error(err))
		}

		if err = tx.Commit(); err != nil {
			h.logger.Error("cache heater fail to commit transaction", zap.Error(err))
		}
	}
}
