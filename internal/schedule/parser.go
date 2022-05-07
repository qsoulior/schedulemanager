package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/qsoulior/scheduleparser"
)

func parseFile(file File, planCh chan Plan, errorCh chan error) {
	initialYear := file.Modified.Year()
	if file.Modified.Month() < 8 {
		initialYear--
	}
	initialDate, err := time.Parse("2006-01-02", fmt.Sprintf("%d-09-01", initialYear))
	fileDataParsed, err := scheduleparser.ParseBytes(file.Data, initialDate)
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
