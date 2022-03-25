package schedule

import (
	"strings"

	"github.com/1asagne/scheduleparser"
)

type File struct {
	Name string
	Data []byte
}

func parseFile(file File, fileCh chan File, errorCh chan error) {
	fileDataParsed, err := scheduleparser.ParseBytes(file.Data)
	if err != nil {
		errorCh <- err
		return
	}
	fileCh <- File{Name: strings.Split(file.Name, ".")[0], Data: fileDataParsed}
}

func ParseFiles(files []File) ([]File, error) {

	fileCh := make(chan File)
	defer close(fileCh)
	errorCh := make(chan error)
	defer close(errorCh)

	for _, file := range files {
		go parseFile(file, fileCh, errorCh)
	}

	filesParsed := make([]File, 0)
	for i := 0; i < len(files); i++ {
		select {
		case fileParsed := <-fileCh:
			filesParsed = append(filesParsed, fileParsed)
		case err := <-errorCh:
			return nil, err
		}
	}
	return filesParsed, nil
}
