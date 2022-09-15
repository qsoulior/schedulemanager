package entity

import (
	"time"

	"github.com/qsoulior/scheduleparser"
)

type Event = scheduleparser.Event

type Schedule struct {
	Modified time.Time `json:"modified"`
	Events   []Event   `json:"events"`
}

type Plan struct {
	Group     string     `json:"group"`
	Active    bool       `json:"active"`
	Schedules []Schedule `json:"schedules"`
}

type PlanInfo struct {
	Group    string    `json:"group"`
	Modified time.Time `json:"modified"`
}
