package dto

import (
	"time"
)

type Answer struct {
	TimeOfSer time.Duration
	TimeOfDes time.Duration
	Mem       int
}

type Ans struct {
	Ans  Answer
	Name string
	Err  error
}
