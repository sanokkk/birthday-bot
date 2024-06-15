package models

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Username string
	Birthday *time.Time
	ChatIds  pq.Int64Array `gorm:"type:bigint[]"`
}

type Chat struct {
	Id     int64 `gorm:"primaryKey;autoIncrement"`
	ChatId int64
	UserId int64
}
