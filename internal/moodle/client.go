package moodle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Token   string
	RootUrl string
}

func (client *Client) Init(username, password, rootUrl string) error {
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

func (client *Client) GetCourseSections(courseId int) ([]Section, error) {
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

	return sections, nil
}

func NewClient(username, password, rootUrl string) (*Client, error) {
	client := new(Client)
	err := client.Init(username, password, rootUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type Content struct {
	Type           string
	FileName       string
	FilePath       string
	FileSize       int
	FileUrl        string
	TimeCreated    int64
	TimeModified   int64
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
