package moodle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func getMoodleEnv() (string, string, string, int, error) {
	moodleUsername := os.Getenv("MOODLE_USERNAME")
	if moodleUsername == "" {
		return "", "", "", 0, errors.New("MOODLE_USERNAME is missing")
	}
	moodlePassword := os.Getenv("MOODLE_PASSWORD")
	if moodlePassword == "" {
		return "", "", "", 0, errors.New("MOODLE_PASSWORD is missing")
	}
	moodleRootUrl := os.Getenv("MOODLE_ROOT_URL")
	if moodleRootUrl == "" {
		return "", "", "", 0, errors.New("MOODLE_ROOT_URL is missing")
	}
	moodleCourseId := os.Getenv("MOODLE_COURSE_ID")
	if moodleCourseId == "" {
		return "", "", "", 0, errors.New("MOODLE_COURSE_ID is missing")
	}
	moodleCourseIdInteger, err := strconv.Atoi(moodleCourseId)
	if err != nil {
		return "", "", "", 0, err
	}
	return moodleUsername, moodlePassword, moodleRootUrl, moodleCourseIdInteger, nil
}

type ScheduleFileInfo struct {
	Name string
	Url  string
}

func getScheduleFilesInfo(sections []Section) ([]ScheduleFileInfo, error) {
	filesInfo := make([]ScheduleFileInfo, 0)
	for _, section := range sections {
		if section.Name != "Общее" {
			for _, module := range section.Modules {
				if module.Name == "Расписание экзаменов" {
					break
				} else if module.ModName == "folder" {
					for _, content := range module.Contents {
						if content.Type == "file" {
							filesInfo = append(filesInfo, ScheduleFileInfo{content.FileName, content.FileUrl})
						}
					}
				}
			}
		}
	}
	return filesInfo, nil
}

type ScheduleFile struct {
	Name string
	Data []byte
}

func downloadScheduleFile(scheduleFileInfo ScheduleFileInfo, token string, scheduleCh chan ScheduleFile, errorCh chan error) {
	resp, err := http.Get(scheduleFileInfo.Url + "&token=" + token)
	if err != nil {
		errorCh <- err
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errorCh <- errors.New(fmt.Sprintf("Response status code: %d", resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorCh <- err
		return
	}

	var bodyJson map[string]string
	_ = json.Unmarshal(body, &bodyJson)

	if errorCode, ok := bodyJson["errorcode"]; ok {
		errorCh <- errors.New(errorCode)
		return
	}
	scheduleCh <- ScheduleFile{scheduleFileInfo.Name, body}
}

func GetScheduleFiles() ([]ScheduleFile, error) {
	moodleUsername, moodlePassword, moodleRootUrl, moodleCourseId, err := getMoodleEnv()
	if err != nil {
		return nil, err
	}

	moodleClient, err := NewMoodleClient(moodleUsername, moodlePassword, moodleRootUrl)
	if err != nil {
		return nil, err
	}

	sections, err := moodleClient.GetCourseSections(moodleCourseId)
	scheduleFilesInfo, err := getScheduleFilesInfo(sections)
	if err != nil {
		return nil, err
	}

	scheduleFiles := make([]ScheduleFile, 0)

	scheduleCh := make(chan ScheduleFile)
	errorCh := make(chan error)

	for _, fileInfo := range scheduleFilesInfo {
		go downloadScheduleFile(fileInfo, moodleClient.Token, scheduleCh, errorCh)
	}

	for i := 0; i < len(scheduleFilesInfo); i++ {
		select {
		case scheduleFile := <-scheduleCh:
			scheduleFiles = append(scheduleFiles, scheduleFile)
		case err := <-errorCh:
			fmt.Println(err)
			return nil, err
		}
	}
	return scheduleFiles, nil
}
