package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	parser "github.com/1asagne/ScheduleParser"
	"github.com/joho/godotenv"
)

type MoodleClient struct {
	Token   string
	RootUrl string
}

func (client *MoodleClient) Init(username, password, rootUrl string) error {
	loginUrl := rootUrl + "/login/token.php?username=%s&password=%s&service=moodle_mobile_app"
	resp, err := http.Get(fmt.Sprintf(loginUrl, username, password))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Response status code: %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var bodyJson map[string]string
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		return err
	}

	if errorCode, ok := bodyJson["errorcode"]; ok {
		return errors.New(errorCode)
	}

	if token, ok := bodyJson["token"]; !ok {
		return errors.New("Response body doesn't contain token")
	} else {
		client.Token = token
		client.RootUrl = rootUrl
		return nil
	}
}

func NewMoodleClient(username, password, rootUrl string) (*MoodleClient, error) {
	moodleClient := new(MoodleClient)
	err := moodleClient.Init(username, password, rootUrl)
	if err != nil {
		return nil, err
	}
	return moodleClient, nil
}

type Content struct {
	Type           string
	FileName       string
	FilePath       string
	FileSize       int
	FileUrl        string
	TimeCreated    int
	TimeModified   int
	SortOrder      int
	Mimetype       string
	IsExternalFile bool
	UserId         int
	Author         string
	License        string
}

type ContentsInfo struct {
	FilesCount     int
	FilesSize      int
	LastModified   int
	MimeTypes      []string
	RepositoryType string
}

type Module struct {
	Id                  int
	Url                 string
	Name                string
	Instance            int
	Description         string
	Visible             int
	UserVisible         bool
	VisibleOnCoursePage int
	ModIcon             string
	ModName             string
	ModPlural           string
	Indent              int
	OnClick             string
	AfterLink           string
	CustomData          string
	NoViewLink          bool
	Completion          int
	Contents            []Content
	ContentsInfo        ContentsInfo
}

type Section struct {
	Id                  int
	Name                string
	Visible             int
	Summary             string
	SummaryFormat       int
	Section             int
	HiddenByNumSections int
	UserVisible         bool
	Modules             []Module
}

type FileInfo struct {
	Name string
	Url  string
}

type SectionInfo struct {
	Name      string
	FilesInfo []FileInfo
}

func (client *MoodleClient) GetFilesInfo(courseId int) ([]FileInfo, error) {
	courseGetContentsUrl := client.RootUrl +
		"/webservice/rest/server.php?moodlewsrestformat=json&wstoken=%s&wsfunction=core_course_get_contents&courseid=%d"
	resp, err := http.Get(fmt.Sprintf(courseGetContentsUrl, client.Token, courseId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sections []Section
	err = json.Unmarshal(body, &sections)
	if err != nil {
		return nil, err
	}

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

func getMoodleEnv(envVars map[string]string) (string, string, string, int, error) {
	moodleUsername := envVars["MOODLE_USERNAME"]
	if moodleUsername == "" {
		return "", "", "", 0, errors.New("MOODLE_USERNAME is missing")
	}
	moodlePassword := envVars["MOODLE_PASSWORD"]
	if moodlePassword == "" {
		return "", "", "", 0, errors.New("MOODLE_PASSWORD is missing")
	}
	moodleRootUrl := envVars["MOODLE_ROOT_URL"]
	if moodleRootUrl == "" {
		return "", "", "", 0, errors.New("MOODLE_ROOT_URL is missing")
	}
	moodleCourseId := envVars["MOODLE_COURSE_ID"]
	if moodleCourseId == "" {
		return "", "", "", 0, errors.New("MOODLE_COURSE_ID is missing")
	}
	moodleCourseIdInteger, err := strconv.Atoi(moodleCourseId)
	if err != nil {
		return "", "", "", 0, err
	}
	return moodleUsername, moodlePassword, moodleRootUrl, moodleCourseIdInteger, nil
}

func downloadMoodleFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
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

func scrap() error {
	envVars, err := godotenv.Read("moodle.env")
	if err != nil {
		return err
	}

	moodleUsername, moodlePassword, moodleRootUrl, moodleCourseId, err := getMoodleEnv(envVars)
	if err != nil {
		return err
	}

	moodleClient, err := NewMoodleClient(moodleUsername, moodlePassword, moodleRootUrl)
	if err != nil {
		return err
	}

	filesInfo, err := moodleClient.GetFilesInfo(moodleCourseId)
	if err != nil {
		return err
	}

	for _, fileInfo := range filesInfo {
		pdfBytes, err := downloadMoodleFile(fileInfo.Url + "&token=" + moodleClient.Token)
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
