package model

import "time"

type Message struct {
	tableName struct{} `sql:"message" pg:",discard_unknown_columns"`

	Id        int64         `sql:"id,pk" json:"id"`
	Text      string        `sql:"text" json:"text"`
	Status    MessageStatus `sql:"status" json:"status"`
	CreatedAt time.Time     `sql:"created_at" json:"-"`
	UpdatedAt time.Time     `sql:"updated_at" json:"-"`
}
