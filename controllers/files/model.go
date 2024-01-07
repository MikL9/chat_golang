package files

import (
	"mime/multipart"
)

type Form struct {
	Files []*multipart.FileHeader `form:"attachment[]" binding:"required"`
	Type  string                  `form:"type"`
	Path  string                  `form:"path"`
	ID    string                  `form:"id"`
}

type ResponseFile struct {
	file      Form
	FileID    int `form:"file_id"`
	FileName  string
	Extension string
	Type      string
	Guid      string
}

type ResponseFiles []struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Extension string
	Type      string
	Guid      string
	Mtype     string
}
