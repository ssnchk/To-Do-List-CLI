package models

import "time"

type Task struct {
	Id          int
	Description string
	CreateTime  time.Time
	IsCompleted bool
}
