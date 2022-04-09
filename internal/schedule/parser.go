package schedule

import (
	"encoding/json"
	"strings"

	"github.com/1asagne/scheduleparser"
)

func parseFile(file File, planCh chan Plan, errorCh chan error) {
	fileDataParsed, err := scheduleparser.ParseBytes(file.Data)
	if err != nil {
		errorCh <- err
		return
	}
	plan := Plan{}
	plan.Group = strings.Split(file.Name, ".")[0]
	schedule := Schedule{}
	if err := json.Unmarshal(fileDataParsed, &schedule.Events); err != nil {
		errorCh <- err
		return
	}
	schedule.Modified = file.Modified
	plan.Schedules = append(plan.Schedules, schedule)
	planCh <- plan
}

func ParseFiles(files []File) ([]Plan, error) {
	planCh := make(chan Plan)
	defer close(planCh)
	errorCh := make(chan error)
	defer close(errorCh)

	for _, file := range files {
		go parseFile(file, planCh, errorCh)
	}

	plans := make([]Plan, 0)
	for i := 0; i < len(files); i++ {
		select {
		case plan := <-planCh:
			plans = append(plans, plan)
		case err := <-errorCh:
			return nil, err
		}
	}
	return plans, nil
}
