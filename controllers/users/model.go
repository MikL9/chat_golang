package users

import "gochat/models"

type (
	ResponseUsers []struct {
		models.User
		Guid     string `json:"guid"`
		FileName string `json:"file_name"`
	}
)
