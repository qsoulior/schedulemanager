package moodle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/1asagne/schedulemanager/internal/schedule"
)

func getEnvVars() (string, string, string, int, error) {
	username := os.Getenv("MOODLE_USERNAME")
	if username == "" {
		return "", "", "", 0, errors.New("MOODLE_USERNAME is missing")
	}
	password := os.Getenv("MOODLE_PASSWORD")
	if password == "" {
		return "", "", "", 0, errors.New("MOODLE_PASSWORD is missing")
	}
	rootUrl := os.Getenv("MOODLE_ROOT_URL")
	if rootUrl == "" {
		return "", "", "", 0, errors.New("MOODLE_ROOT_URL is missing")
	}
	courseId := os.Getenv("MOODLE_COURSE_ID")
	if courseId == "" {
		return "", "", "", 0, errors.New("MOODLE_COURSE_ID is missing")
	}
	courseIdInt, err := strconv.Atoi(courseId)
	if err != nil {
		return "", "", "", 0, err
	}
	return username, password, rootUrl, courseIdInt, nil
}

type FileInfo struct {
	Name string
	Url  string
}

func getFilesInfo(sections []Section) ([]FileInfo, error) {
	filesInfo := make([]FileInfo, 0)
	for _, section := range sections {
		if section.Name != "Общее" {
			for _, module := range section.Modules {
				if module.Name == "Расписание экзаменов" {
					break
				} else if module.ModName == "folder" {
					for _, content := range module.Contents {
						if content.Type == "file" {
							filesInfo = append(filesInfo, FileInfo{content.FileName, content.FileUrl})
						}
					}
				}
			}
		}
	}
	return filesInfo, nil
}

func downloadFile(fileInfo FileInfo, accessToken string, fileCh chan schedule.File, errorCh chan error) {
	resp, err := http.Get(fileInfo.Url + "&token=" + accessToken)
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
	fileCh <- schedule.File{Name: fileInfo.Name, Data: body}
}

func DownloadFiles() ([]schedule.File, error) {
	username, password, rootUrl, courseId, err := getEnvVars()
	if err != nil {
		return nil, err
	}

	client, err := NewClient(username, password, rootUrl)
	if err != nil {
		return nil, err
	}

	sections, err := client.GetCourseSections(courseId)
	filesInfo, err := getFilesInfo(sections)
	if err != nil {
		return nil, err
	}

	files := make([]schedule.File, 0)

	fileCh := make(chan schedule.File)
	errorCh := make(chan error)

	for _, fileInfo := range filesInfo {
		go downloadFile(fileInfo, client.Token, fileCh, errorCh)
	}

	for i := 0; i < len(filesInfo); i++ {
		select {
		case file := <-fileCh:
			files = append(files, file)
		case err := <-errorCh:
			return nil, err
		}
	}
	return files, nil
}
