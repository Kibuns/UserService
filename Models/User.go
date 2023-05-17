package Models

import (
	"time"
)

type User struct {
	UserID    string    `json:"userid"`
	UserName	string	`json:"username"`
	Password  string    `json:"password"`
	Created time.Time `json:"created"`
}
