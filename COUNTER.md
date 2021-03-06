# Микросервис счетчиков

## Содержание
1. [ Задание ](#task)
    - [ Цель ](#task-goal)
    - [ Приобретенные навыки ](#task-skills)
    - [ Постановка задачи ](#task-statement)
2. [ Сведения ](#information)
    - [ Используемые инструменты ](#information-tools)
    - [ Характеристики железа ](#information-computer)
3. [ Ход работы ](#work)
    - [ Предложенный вариант решения ](#work-solution)
    - [ Оценка предложенного варианта решения ](#work-solution-estimate)
    - [ Предназначение микросервиса destrib-syncer ](#work-destrib-syncer-microservice)
    - [ Выполнение ](#work-execute)
        - [ Подготовка ](#work-execute-preparation)
        - [ Начало обмена сообщениями ](#work-execute-start-chatting)
            - [ Терминальное окно Боба ](#work-execute-start-chatting-bob-term)
            - [ Терминальное окно Алисы ](#work-execute-start-chatting-alice-term)
            - [ Отправка сообщений ](#work-execute-start-chatting-send-msg)
        - [ Окончание переписки ](#work-execute-stop-chatting)
        - [ Историческая выгрузка сообщений ](#work-execute-history-dump)
4. [ Итоги ](#results)

<img align="right" width="600" src="static/counter/preview.jpg">

<a name="task"></a>
## Задание
Внедрение микросервиса счетчиков в существующую инфраструктуру.

<a name="task-goal"></a>
### Цель
Реализовать и внедрить микросервис счетчиков в существующую инфраструктуру системы.

<a name="task-skills"></a>
### Приобретенные навыки
В результате выполненного задания необходимо приобрести следующие навыки:
- проектирование масштабируемых архитектур;
- разработка отказоустойчивых микросервисов;
- использование механизма кеширования.

<a name="task-statement"></a>
### Постановка задачи
В процессе достижения цели необходимо:
- разработать и внедрить микросервис счетчиков;
- обеспечить консистентность между величиной счетчика и реальным числом непрочитанных сообщений.

<a name="information"></a>
## Сведения

<a name="information-tools"></a>
### Используемые инструменты
Для выполнения задания понадобятся следующие инструменты:
- [docker](https://docs.docker.com/get-docker/) (>= version 19.03.8) & [docker compose](https://docs.docker.com/compose/install/) (>= version 1.25.5);
- [curl](https://curl.haxx.se/download.html) (>= version 7.68.0);
- [websocat](https://github.com/vi/websocat/releases) (>= version 1.6.0);
- [jq](https://stedolan.github.io/jq/download/) (>= version 1.5).

<a name="information-computer"></a>
### Характеристики железа
Задание выполнялось на железе со следующими характеристиками:
- CPU - AMD Ryzen 9: 12 ядер 24 потока;
- RAM - 2xHyperX Fury Black: DDR4 DIMM 3000MHz 8GB;
- SSD - Intel® SSD 540s Series: 480GB, 2.5in SATA 6Gb/s, 16nm, TLC


<a name="work"></a>
## Ход работы

<a name="work-solution"></a>
### Предложенный вариант решения
При решении таких задач, как:
- разработка микросервивиса счетчиков непрочитанных сообщений пользователем;
- обеспечения консистентности величины счетчика с реальным числом непрочитанных сообщений; 

предлагается решить их в простом и "прямом" виде, не прибегая к сложным шаблонам синхронизации (SAGA, механизм 
двухфазного коммита и тд). Предлагается реализовать помимо самого микросервиса счетчиков **counter** микросервис 
**destrib-syncer**, задача которого состоит в том, чтобы обеспечить консистентность количества непрочитанных сообщений,
передаваемых пользователю, между микросервисами **messenger** и **counter** соответственно.  

Разработанная инфраструктура имеет следующее представление: </br>
<p align="center">
    <img src="static/counter/infrastructure-schema.png">
</p>

<a name="work-solution-estimate"></a>
### Оценка предложенного варианта решения
Как видно из описания выше есть определенные ограничения и недостатки данного подхода. В целом это связано с 
масштабируемостью экземпляров микросервисов **messenger** и **counter**, если, безусловно, такая проблема
возникнет. Сейчас это проблема решается лишь путем корректировки инфраструктурно-транспортного слоя на стороне 
микросервиса **destrib-syncer**. В дальнейшем, само себе полагающимся, необходимо использовать более серьезные шаблоны 
проектирования в виде, например, **Хореографической саги**, которая с помощью очередей той же шины данных NATS-Streaming
или Kafka даст ту самую окраску распределенности и масштабируемости без изменения того же микросервиса 
**destrib-syncer**. 

Подчеркну, что данный подход к решению поставленной задачи не оптимален и не универсален в общем виде, однако для 
решения конкретно нашей задачи - вполне может быть приемлемым.

<a name="work-destrib-syncer-microservice"></a>
### Предназначение микросервиса destrib-syncer
Как уже было сказано выше микросервис **destrib-syncer** обеспечивает консистентность количества непрочитанных
сообщений, однако на этом его участь далеко не заканчивается. Теперь WS соединение держит именно он, а общается уже с
микросервисом **messenger** по REST API, для таких операций как:
- создание чатов между собеседниками;
- создание сообщений пользователями в чатах;
- выгрузка исторических сообщений.

Безусловно, что маршрутизация сообщений между конечными пользователями также лежит на микросервисе **destrib-syncer**,
но это как раз ключевой момент, с помощью которого можно отслеживать счетчик непрочитанных сообщений каждому из 
пользователей и инкриминировать или декриминировать его в микросервисе **counter**. 

Ниже представлены схемы алгоритмов **трех основных сценариев** действий микросервиса **destrib-syncer**, а именно:
- создание чата;
- отправка и маршрутизация сообщений;
- выгрузка сообщений с offset'ом и limit'ом.


Схема алгоритма создание чатов между собеседниками имеет следующее представление: </br>
<p align="center">
    <img src="static/counter/create_chat_algorithm.png">
</p>

Важно отметить, что при создании чата, на каждого пользователя создается **уникальный счетчик** в виде (**userID** + 
**chatID**), в котором ведутся значения непрочитанных сообщений.


Схема алгоритма отправки сообщений от пользователей к пользователям представлена следующим образом:</br>
<p align="center">
    <img src="static/counter/send_message_algorithm.png">
</p>

На схеме выпукло показана при каких условиях **счетчик инкриминируется** у собеседника, а именно в случае, если **нет 
активных ws-соединений собеседника**.

И наконец схема алгоритма выгрузки сообщений в виде offset и limit'а:</br>
<p align="center">
    <img src="static/counter/get_messages_algorithm.png">
</p>
В данной схеме сосредаточена самая сложна и интересная часть - а именно актуализация счетчика непрочитанных сообщений.

<a name="work-execute"></a>
### Выполнение
Склонируем наш проект:
```shell
git clone https://github.com/teploff/otus-highload.git && cd otus-highload
```

Поднимаем инфраструктуру и применяем миграции:
```shell
make init && make migrate && make app
```

<a name="work-execute-preparation"></a>
#### Подготовка
Для демонстрации работы разработанной системы, на примере осуществления диалогов между пользователями, создадим двух 
пользователей в системе: **Боба** и **Алису**.
```shell script
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890", "name": "Bob", "surname": "Tallor", "birthday": "1994-04-10T20:21:25+00:00", "sex": "male", "city": "New Yourk", "interests": "programming"}' \
    http://localhost:10000/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "alice@email.com", "password": "1234567890", "name": "Alice", "surname": "Swift", "birthday": "1995-10-10T20:21:25+00:00", "sex": "female", "city": "California", "interests": "running"}' \
    http://localhost:10000/auth/sign-up
```

<a name="work-execute-start-chatting"></a>
#### Начало обмена сообщениями
Для того, чтобы начать обмениваться сообщениями, необходимо открыть два терминальных окна. Одно будет принадлежать Бобу,
другое Алисе.

<a name="work-execute-start-chatting-bob-term"></a>
##### Терминальное окно Боба
В первом терминальном окне, предназначенном для Боба, получим access token командой:
```shell script
export BOB_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
``` 

Проверим наличие access token-а:
```shell script
echo $BOB_ACCESS_TOKEN
```

Для того, чтобы от лица Боба создать чат с Алисой, необходимо получить ее ID, так как он понадобится для указания 
собеседника. Воспользуемся endpoint-ом на получения ID пользователя, зная его email:
```shell script
export ALICE_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    http://localhost:10000/auth/user/get-by-email?email=alice@email.com | jq -r '.user_id')
```

Проверим, что запрос успешно выполнился, применив команду:
```shell script
echo $ALICE_ID
```

Создадим чат от лица Боба с Алисой:
```shell script
export CHAT_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    -d '{"companion_id": "'"$ALICE_ID"'"}' \
    http://localhost:10000/messenger/create-chat | jq -r '.chat_id')
```

Проверим, что в переменной окружения находится UUID созданного чата:
```shell script
echo $CHAT_ID
```

<a name="work-execute-start-chatting-alice-term"></a>
##### Терминальное окно Алисы
Во втором терминальном окне, предназначенном для Алисы, получим access token командой:
```shell script
export ALICE_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "alice@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
``` 

Для установления websocket-ного соединения со стороны Алисы к серверу необходимо ввести команду:
```shell script
websocat ws://localhost:10000/messenger/ws\?token=${ALICE_ACCESS_TOKEN}
```

<a name="work-execute-start-chatting-send-msg"></a>
##### Отправка сообщений
Теперь, когда Алиса установила websocket-ное соединение с сервером, она готова принимать сообщения от Боба.

Для этого необходимо вернуться в первое терминальное окно, принадлежащее Бобу. Установить также websocket-ное соединение
командой:
```shell script
websocat ws://localhost:10000/messenger/ws\?token=${BOB_ACCESS_TOKEN}
```
И зная ChatID, например, если он имеет значение **e8d3dc26-a218-4ca1-ae4b-da38b27ed9b3**, отправить сообщения следующего
вида:
```shell script
{"topic":"messenger", "action": "send-message", "payload":"{\"chat_id\":\"e8d3dc26-a218-4ca1-ae4b-da38b27ed9b3\", \"messages\":[{\"text\": \"Hello, Alice!\", \"status\": \"created\"}]}"}
{"topic":"messenger", "action": "send-message", "payload":"{\"chat_id\":\"e8d3dc26-a218-4ca1-ae4b-da38b27ed9b3\", \"messages\":[{\"text\": \"What is up?\", \"status\": \"created\"}]}"}
{"topic":"messenger", "action": "send-message", "payload":"{\"chat_id\":\"e8d3dc26-a218-4ca1-ae4b-da38b27ed9b3\", \"messages\":[{\"text\": \"I miss you!\", \"status\": \"created\"}]}"}
```

Теперь перейдем в терминал Алисы и удостоверимся, что получили все три сообщения от Боба. В терминале должны увидеть 
следующее:</br>
<p align="center">
  <img src="static/counter/websocket-get-messages.png">
</p>

<a name="work-execute-stop-chatting"></a>
#### Окончание переписки
Для того, чтобы закрыть websocket-ные соединения от каждого из пользователей с сервером, необходимо в каждом 
терминальном окне ввести сочетание клавиш **Ctrl+C**.

<a name="work-execute-history-dump"></a>
#### Историческая выгрузка сообщений
Для того чтобы выгрузить сообщения по конкретному чату, необходимо, находясь в терминальном окне Алисы и зная id чата 
(в данном контексте он имеет значение **e8d3dc26-a218-4ca1-ae4b-da38b27ed9b3**), отправить запрос, с помощью которого
получим сообщения, которые ей отослал Боб:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${ALICE_ACCESS_TOKEN}" \
http://localhost:10000/messenger/messages?chat_id=e8d3dc26-a218-4ca1-ae4b-da38b27ed9b3
```

Если все прошло успешно, должны увидеть нечто похожее: </br>
<p align="center">
    <img src="static/counter/http-get-messages.png">
</p>

<a name="results"></a>
## Итоги
В ходе выполнения задания:
- был описан процесс сборки и конфигурирования программного комплекса;
- были разработан и внедрен микросервис счетчиков в существующую инфраструктуру;
- была обеспечена консистентность между величиной счетчика и реальным числом непрочитанных сообщений.