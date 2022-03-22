package moodle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	parser "github.com/1asagne/ScheduleParser"
	"github.com/joho/godotenv"
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

type MoodleFileInfo struct {
	Name string
	Url  string
}

func getMoodleFilesInfo(sections []Section) ([]MoodleFileInfo, error) {
	filesInfo := make([]MoodleFileInfo, 0)
	for _, section := range sections {
		if section.Name != "Общее" {
			for _, module := range section.Modules {
				if module.Name == "Расписание экзаменов" {
					break
				} else if module.ModName == "folder" {
					for _, content := range module.Contents {
						if content.Type == "file" {
							filesInfo = append(filesInfo, MoodleFileInfo{content.FileName, content.FileUrl})
						}
					}
				}
			}
		}
	}
	return filesInfo, nil
}

func downloadMoodleFile(moodleFile MoodleFileInfo, token string) ([]byte, error) {
	resp, err := http.Get(moodleFile.Url + "&token=" + token)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Response status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bodyJson map[string]string
	_ = json.Unmarshal(body, &bodyJson)

	if errorCode, ok := bodyJson["errorcode"]; ok {
		return nil, errors.New(errorCode)
	}
	return body, err
}

func Scrap() error {
	err := godotenv.Load("moodle.env")
	if err != nil {
		return err
	}

	moodleUsername, moodlePassword, moodleRootUrl, moodleCourseId, err := getMoodleEnv()
	if err != nil {
		return err
	}

	moodleClient, err := NewMoodleClient(moodleUsername, moodlePassword, moodleRootUrl)
	if err != nil {
		return err
	}

	sections, err := moodleClient.GetCourseSections(moodleCourseId)
	filesInfo, err := getMoodleFilesInfo(sections)
	if err != nil {
		return err
	}

	for _, fileInfo := range filesInfo {
		pdfBytes, err := downloadMoodleFile(fileInfo, moodleClient.Token)
		if err != nil {
			return err
		}
		jsonBytes, err := parser.ParseScheduleBytes(pdfBytes)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n\n", string(jsonBytes))
	}
	return nil
}
