package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/urfave/cli/v3"
)

type CreateUserStruct struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserStruct struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

var BASE = "http://localhost:8000/api/v1"

type ResponseStruct struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Key     string `json:"key"`
}

func CreateUser(_ context.Context, cmd *cli.Command) error {
	u := CreateUserStruct{Email: cmd.String("email"), Username: cmd.String("username"), Password: cmd.String("password")}
	b, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("Failed to Serialize to JSON from native Go struct type: %v", err)
		return err
	}
	body := bytes.NewBuffer(b)
	url := BASE + "/auth/sign-up/"
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Fatalf("Failed to create resource at: %s and the error is: %v\n", url, err)
		return err
	}
	defer resp.Body.Close()
	var response ResponseStruct
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Print(err)
	}

	if resp.StatusCode == 400 {
		log.Printf(response.Message)
		return nil
	}

	if resp.StatusCode == 201 {
		log.Print(response.Message)
		log.Print("Please login with the password you have used to create the account")
	}

	return nil
}

func LoginUser(_ context.Context, cmd *cli.Command) error {
	u := LoginUserStruct{Username: cmd.String("username"), Password: cmd.String("password")}
	// serialize the struct
	b, err := json.Marshal(u)
	if err != nil {
		log.Print("Error with json serializing", err)
		return err
	}

	body := bytes.NewBuffer(b)
	url := BASE + "/auth/sign-in/"
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Print("Error sending request to the server")
	}
	defer resp.Body.Close()
	var response ResponseStruct
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Print("Error in deserializing the response from the server")
	}
	print("Success")
	print(response.Key)
	return nil
}
