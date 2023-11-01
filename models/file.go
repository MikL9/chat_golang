package models

import (
	"time"

	"github.com/guregu/null"
)

type FileType int

// TODO: определиться, какие есть типы файлов на самом деле
const (
	_ FileType = iota
	ProjectLogo
	ProjectPhoto
	ReportPhoto
	TaskPhoto
	VisitPhoto
)

type File struct {
	ID        int       `json:"id"`
	Guid      string    `json:"guid"`
	Type      string    `json:"type"`
	ParentID  int       `json:"parent_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Fullname  string    `json:"fullname"`
	Extension string    `json:"extension"`
	MimeType  string    `json:"mime_type"`
	Size      int       `json:"size"`
	Mtype     string    `json:"mtype"`
	Index     int       `json:"index"`
	Created   time.Time `json:"created"`
	Updated   null.Time `json:"updated"`
}
