package main

import (
	"os"
	"path/filepath"
	"strings"

	parser "github.com/1asagne/ScheduleParser"
)

const schedulesFolderName = "./assets/schedules"

func parseSchedules() error {
	files, err := os.ReadDir(schedulesFolderName + "/pdf")
	if err != nil {
		return err
	}
	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		if fileExt == ".pdf" {
			filePath := schedulesFolderName + "/pdf/" + file.Name()
			filePathNew := schedulesFolderName + "/json/" + strings.TrimSuffix(file.Name(), fileExt) + ".json"
			err = parser.ParseScheduleFile(filePath, filePathNew)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
