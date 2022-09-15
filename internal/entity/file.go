package entity

import "time"

type File struct {
	Name     string
	Modified time.Time
	Data     []byte
}
