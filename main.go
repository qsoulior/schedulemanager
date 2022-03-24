package main

import (
	"fmt"

	"github.com/1asagne/schedulemanager/internal/moodle"
	sp "github.com/1asagne/scheduleparser"
	"github.com/joho/godotenv"
)

func parseSchedule(scheduleFile moodle.ScheduleFile, scheduleCh chan moodle.ScheduleFile, errorCh chan error) {
	fileDataParsed, err := sp.ParseScheduleBytes(scheduleFile.Data)
	if err != nil {
		errorCh <- err
		return
	}
	scheduleCh <- moodle.ScheduleFile{Name: scheduleFile.Name, Data: fileDataParsed}
}

func main() {
	err := godotenv.Load("moodle.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	scheduleFiles, err := moodle.GetScheduleFiles()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Getting files completed")

	scheduleFilesParsed := make([]moodle.ScheduleFile, 0)

	scheduleCh := make(chan moodle.ScheduleFile)
	defer close(scheduleCh)
	errorCh := make(chan error)
	defer close(errorCh)

	for _, scheduleFile := range scheduleFiles {
		go parseSchedule(scheduleFile, scheduleCh, errorCh)
	}

	for i := 0; i < len(scheduleFiles); i++ {
		select {
		case scheduleFileParsed := <-scheduleCh:
			scheduleFilesParsed = append(scheduleFilesParsed, scheduleFileParsed)
		case err := <-errorCh:
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Parsing files completed")

	for _, i := range scheduleFilesParsed {
		fmt.Println(i.Name)
		fmt.Println(string(i.Data))
		fmt.Println()
	}

	fmt.Println("Main completed")
}
