package schedule

import (
	"encoding/json"
	"strings"

	sp "github.com/1asagne/scheduleparser"
)

type File struct {
	Name string
	Data []byte
}

type Schedule struct {
	Name   string
	Events []sp.Event
}

func parseFile(file File, scheduleCh chan Schedule, errorCh chan error) {
	fileDataParsed, err := sp.ParseBytes(file.Data)
	if err != nil {
		errorCh <- err
		return
	}
	schedule := Schedule{}
	schedule.Name = strings.Split(file.Name, ".")[0]
	if err := json.Unmarshal(fileDataParsed, &schedule.Events); err != nil {
		errorCh <- err
		return
	}
	scheduleCh <- schedule
}

func ParseFiles(files []File) ([]Schedule, error) {
	scheduleCh := make(chan Schedule)
	defer close(scheduleCh)
	errorCh := make(chan error)
	defer close(errorCh)

	for _, file := range files {
		go parseFile(file, scheduleCh, errorCh)
	}

	schedules := make([]Schedule, 0)
	for i := 0; i < len(files); i++ {
		select {
		case schedule := <-scheduleCh:
			schedules = append(schedules, schedule)
		case err := <-errorCh:
			return nil, err
		}
	}
	return schedules, nil
}
