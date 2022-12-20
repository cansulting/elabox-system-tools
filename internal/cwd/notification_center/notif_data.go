package main

import "time"

const (
	Unread  NotifStatus = 0
	Read    NotifStatus = 1
	Deleted NotifStatus = 2
)

type NotifStatus uint

type SelectionData struct {
	Id      string `json:"id"`
	Caption string `json:"caption"`
}

type NotifData struct {
	Id         string          `json:"id"`
	Icon       string          `json:"icon"`
	Title      string          `json:"title"`
	Message    string          `json:"message"`
	PackageId  string          `json:"packageId"`
	Extra      string          `json:"extra,omitempty"`
	Status     NotifStatus     `json:"status"`
	Selections []SelectionData `json:"selections,omitempty"`
	ReceivedAt time.Time       `json:"receivedAt"`
}
