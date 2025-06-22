package data

import "time"

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Working   bool      `json:"working"`
	Duration  int       `json:"duration"`
}
