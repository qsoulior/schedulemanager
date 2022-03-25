package main

import (
	"fmt"

	"github.com/1asagne/schedulemanager/internal/moodle"
	"github.com/1asagne/schedulemanager/internal/schedule"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("moodle.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	scheduleFiles, err := moodle.DownloadFiles()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Getting files completed")

	scheduleFilesParsed, err := schedule.ParseFiles(scheduleFiles)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Parsing files completed")

	for _, i := range scheduleFilesParsed {
		fmt.Println(i.Name)
		fmt.Println(string(i.Data))
	}
}
