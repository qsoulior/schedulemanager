package moodle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	*http.Client
	token string
	host  string
}

func NewClient(host, username, password string) (*Client, error) {
	client := &Client{&http.Client{}, "", host}
	err := client.UpdateToken(username, password)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) UpdateToken(username, password string) error {
	loginUrl := client.host + "/login/token.php?username=%s&password=%s&service=moodle_mobile_app"
	resp, err := client.Get(fmt.Sprintf(loginUrl, username, password))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	result := new(struct {
		Token        string
		PrivateToken string
		Error        string
		ErrorCode    string
	})

	if err := decoder.Decode(result); err != nil {
		return err
	}

	if errorCode := result.ErrorCode; errorCode != "" {
		return errors.New(errorCode)
	}

	if token := result.Token; token == "" {
		return errors.New("result doesn't contain token")
	} else {
		client.token = token
		return nil
	}
}

func (client *Client) GetCourseSections(courseId int) ([]Section, error) {
	courseGetContentsUrl := client.host +
		"/webservice/rest/server.php?moodlewsrestformat=json&wstoken=%s&wsfunction=core_course_get_contents&courseid=%d"
	resp, err := client.Get(fmt.Sprintf(courseGetContentsUrl, client.token, courseId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	sections := make([]Section, 0)
	decoder := json.NewDecoder(resp.Body)

	if _, err := decoder.Token(); err != nil {
		log.Fatal(err)
	}

	for decoder.More() {
		section := new(Section)
		if err := decoder.Decode(section); err == nil {
			sections = append(sections, *section)
		}
	}

	if _, err := decoder.Token(); err != nil {
		log.Fatal(err)
	}
	return sections, nil
}

func (client *Client) GetFileBytes(content *Content) ([]byte, error) {
	if content.Type != "file" {
		return nil, errors.New("content is not a file")
	}
	resp, err := client.Get(fmt.Sprintf("%s&token=%s", content.FileUrl, client.token))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	defer resp.Body.Close()

	if contentType == "application/pdf" {
		return io.ReadAll(resp.Body)
	}

	var result *struct {
		Error     string
		ErrorCode string
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(result); err != nil {
		return nil, err
	}

	if errorCode := result.ErrorCode; errorCode != "" {
		return nil, errors.New(errorCode)
	}

	return nil, nil
}
