package models

import (
	"github.com/ShynaliOtep/paydal-api/services/filesystem"
	"time"
)

type File struct {
	ID         uint
	Filename   string `gorm:"size:50;not null;" json:"filename"`
	Path       string `gorm:"size:100;not null;" json:"path"`
	Size       int    `json:"size"`
	Extension  string `gorm:"size:10;null;" json:"extension"`
	UploadedAt time.Time
}

func SaveFile(path string) (File, error) {
	file := File{
		Filename:   filesystem.GetFilenameFromPath(path),
		Path:       path,
		Size:       0,
		Extension:  filesystem.GetExtension(path),
		UploadedAt: time.Now(),
	}

	err := DB.Create(&file).Error
	if err != nil {
		return File{}, err
	}

	return file, nil
}
