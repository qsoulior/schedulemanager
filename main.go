package main

import (
	"fmt"

	"github.com/1asagne/ScheduleManager/internal/moodle"
)

func main() {
	err := moodle.Scrap()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Main completed")
}
