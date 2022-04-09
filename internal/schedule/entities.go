package schedule

import (
	"time"

	"github.com/1asagne/scheduleparser"
)

type Event = scheduleparser.Event

type Schedule struct {
	Modified time.Time
	Events   []Event
}

type File struct {
	Name     string
	Modified time.Time
	Data     []byte
}

type Plan struct {
	Group     string
	Active    bool
	Schedules []Schedule
}
