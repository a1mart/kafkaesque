package main

import (
	"fmt"

	validators "github.com/a1mart/kafkaesque/internal/midas"
)

type Inner struct {
	Attribute1 string `json:"attribute" validate:"required,email" sanitize:"lowercase"`
}

type UserRegisterRequest struct {
	Email     string `json:"email" validate:"required,email" sanitize:"lowercase"`
	Password  string `json:"password" validate:"min=2,max=32,regex=^[a-zA-Z0-9]+$"`
	Username  string `json:"username" validate:"min=2,max=32,regex=^[a-zA-Z0-9]+$"`
	IP        string `json:"ip,omitempty"` // Optional IP field
	Attribute Inner
}

func main() {
	req := &UserRegisterRequest{
		Email:     "abc@gmail.com",
		Password:  "123",
		Username:  "aidan",
		Attribute: Inner{Attribute1: "a"},
	}

	validated := validators.Validate(req)
	if validated != nil {
		fmt.Println("Error:", validated.Errors)
	}

}
