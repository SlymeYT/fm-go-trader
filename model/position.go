package model

import "time"

type Position struct {
	LastUpdateTimestamp 	time.Time
	Symbol 					string
	Quantity 				int64

}