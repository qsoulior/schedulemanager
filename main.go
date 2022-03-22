package main

import (
	"fmt"

	"github.com/1asagne/ScheduleManager/internal/moodle"
	parser "github.com/1asagne/ScheduleParser"
	"github.com/joho/godotenv"
)

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

	scheduleFilesParsed := make([]moodle.ScheduleFile, 0)

	for _, scheduleFile := range scheduleFiles {
		fileDataParsed, err := parser.ParseScheduleBytes(scheduleFile.Data)
		if err != nil {
			fmt.Println(err)
			return
		}
		scheduleFilesParsed = append(scheduleFilesParsed, moodle.ScheduleFile{Name: scheduleFile.Name, Data: fileDataParsed})
	}

	for _, i := range scheduleFilesParsed {
		fmt.Println(i.Name)
		fmt.Println(string(i.Data))
		fmt.Println()
	}

	fmt.Println("Main completed")
}
