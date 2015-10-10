package main

import (
	"encoding/json"
	"time"
)

type User struct {
	Nick string `json: "nick"`
}

type Message struct {
	User      string          `json:"user"`
	Edited    int             `json:"edited"`
	Content   json.RawMessage `string:"content,omitempty"`
	Id        int             `json:"id"`
	CreatedAt CustomTime      `json:"created_at"`
}


type CustomTime struct {
	time.Time
}

var nilTime = (time.Time{}).UnixNano()

const jsonFormatDate = "2006-01-02T15:04:05.000Z"

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

func (d *CustomTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	d.Time, err = time.Parse(jsonFormatDate, string(b))
	return err
}