package models

import "errors"

type User struct {
	Id       string `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (req *LoginRequest) Validate() error {
	var err error
	if req.Login == "" {
		err = errors.New("the Login field is required")
	}

	if req.Password == "" {
		err = errors.New("the Password field is required")
	}

	return err
}
