// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domain

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson163c17a9DecodeSocialInternalDomain(in *jlexer.Lexer, out *News) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "owner":
			easyjson163c17a9Decode(in, &out.Owner)
		case "content":
			out.Content = string(in.String())
		case "create_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreateTime).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson163c17a9EncodeSocialInternalDomain(out *jwriter.Writer, in News) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"owner\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjson163c17a9Encode(out, in.Owner)
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		out.String(string(in.Content))
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Raw((in.CreateTime).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v News) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson163c17a9EncodeSocialInternalDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v News) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson163c17a9EncodeSocialInternalDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *News) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson163c17a9DecodeSocialInternalDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *News) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson163c17a9DecodeSocialInternalDomain(l, v)
}
func easyjson163c17a9Decode(in *jlexer.Lexer, out *struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Sex     string `json:"sex"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "sex":
			out.Sex = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson163c17a9Encode(out *jwriter.Writer, in struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Sex     string `json:"sex"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"surname\":"
		out.RawString(prefix)
		out.String(string(in.Surname))
	}
	{
		const prefix string = ",\"sex\":"
		out.RawString(prefix)
		out.String(string(in.Sex))
	}
	out.RawByte('}')
}
