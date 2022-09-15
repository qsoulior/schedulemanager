package repository

import (
	"time"

	"github.com/qsoulior/schedulemanager/internal/entity"
	"github.com/qsoulior/schedulemanager/pkg/moodle"
)

type Moodle struct {
	client *moodle.Client
}

func (repo *Moodle) GetFilesInfo(courseId int) ([]moodle.Content, error) {
	sections, err := repo.client.GetCourseSections(courseId)
	if err != nil {
		return nil, err
	}

	contents := make([]moodle.Content, 0)
	for _, section := range sections {
		if section.Name != "Общее" && section.Name != "Списки групп" {
			for _, module := range section.Modules {
				if module.Name == "Расписание экзаменов" {
					break
				} else if module.ModName == "folder" {
					for _, content := range module.Contents {
						if content.Type == "file" {
							contents = append(contents, content)
						}
					}
				}
			}
		}
	}
	return contents, nil
}

func (repo *Moodle) GetFiles(filesInfo []moodle.Content) ([]entity.File, error) {
	fileCh := make(chan *entity.File)
	defer close(fileCh)
	errCh := make(chan error)
	defer close(errCh)

	for _, fileInfo := range filesInfo {
		go func(content moodle.Content) {
			fileBytes, err := repo.client.GetFileBytes(&content)
			if err != nil {
				errCh <- err
				return
			}

			fileCh <- &entity.File{Name: content.FileName, Modified: time.Unix(content.TimeModified, 0), Data: fileBytes}
		}(fileInfo)
	}

	files := make([]entity.File, 0, len(filesInfo))

	for i := 0; i < len(filesInfo); i++ {
		select {
		case file := <-fileCh:
			files = append(files, *file)
		case err := <-errCh:
			return nil, err
		}
	}
	return files, nil
}

func NewMoodle(client *moodle.Client) PlanWeb {
	return &Moodle{client}
}
