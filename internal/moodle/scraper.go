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

func downloadScheduleFile(scheduleFileInfo ScheduleFileInfo, token string) (ScheduleFile, error) {
	resp, err := http.Get(scheduleFileInfo.Url + "&token=" + token)
	if err != nil {
		return ScheduleFile{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ScheduleFile{}, errors.New(fmt.Sprintf("Response status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ScheduleFile{}, err
	}

	var bodyJson map[string]string
	_ = json.Unmarshal(body, &bodyJson)

	if errorCode, ok := bodyJson["errorcode"]; ok {
		return ScheduleFile{}, errors.New(errorCode)
	}
	return ScheduleFile{scheduleFileInfo.Name, body}, err
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

	for _, fileInfo := range scheduleFilesInfo {
		scheduleFile, err := downloadScheduleFile(fileInfo, moodleClient.Token)
		if err != nil {
			return nil, err
		}
		scheduleFiles = append(scheduleFiles, scheduleFile)
	}
	return scheduleFiles, nil
}
