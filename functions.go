package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	loggedIn, _ := IsLoggedIn()
	if loggedIn {
		fmt.Print("You are already logged in. Please use the logout function to proceed")
		return nil
	}
	u := CreateUserStruct{Email: cmd.String("email"), Username: cmd.String("username"), Password: cmd.String("password")}
	b, err := json.Marshal(u)
	if err != nil {
		fmt.Printf("Failed to Serialize to JSON from native Go struct type: %v", err)
		return err
	}
	body := bytes.NewBuffer(b)
	url := BASE + "/auth/sign-up/"
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		fmt.Printf("Failed to create resource at: %s and the error is: %v\n", url, err)
		return err
	}
	defer resp.Body.Close()
	var response ResponseStruct
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Print(err)
	}

	if resp.StatusCode == 400 {
		fmt.Printf(response.Message)
		return nil
	}

	if resp.StatusCode == 201 {
		fmt.Print(response.Message)
		fmt.Print("Please fmtin with the password you have used to create the account")
	}

	return nil
}

func LoginUser(_ context.Context, cmd *cli.Command) error {
	loggedIn, _ := IsLoggedIn()
	if loggedIn {
		fmt.Print("You are already logged in. Please use the logout function to proceed")
		return nil
	}
	u := LoginUserStruct{Username: cmd.String("username"), Password: cmd.String("password")}
	// serialize the struct
	b, err := json.Marshal(u)
	if err != nil {
		fmt.Print("Error with json serializing", err)
		return err
	}

	body := bytes.NewBuffer(b)
	url := BASE + "/auth/sign-in/"
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		fmt.Print("Error sending request to the server")
	}
	defer resp.Body.Close()
	var response ResponseStruct
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Print("Error in deserializing the response from the server")
	}

	err = WriteConfig(response.Key)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print("User logged in successfully")
	return nil
}

func Logout(_ context.Context, cmd *cli.Command) error {
	err := WriteConfig("")
	if err != nil {
		fmt.Print(err)
		return err
	}
	fmt.Print("Logged out successfully!")
	return nil
}

type CreateKeyStruct struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Password string `json:"password"`
}

type SimpleResponseStruct struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func CreateKey(_ context.Context, cmd *cli.Command) error {
	// get the args
	// there should be two args
	// the name of the key and the key itself
	// upon finding the absence, throw an error reminding them the correct use of the command
	loggedIn, token := IsLoggedIn()
	if !loggedIn {
		fmt.Println("Please create an account or login to manage your keys")
		return errors.New("Create an account to use this command or login ")
	}

	args := cmd.Args()
	if args.Len() != 2 {
		fmt.Print("name of the key and the key itself are required \n")
		fmt.Print("Use the command as mentioned below: create service_name password")
		return errors.New("Not enough arguments provided")
	}

	// get the password from the user
	var password string
	fmt.Print("Enter your password: ")
	_, err := fmt.Scan(&password)
	// having confirmed that the arguments are provided. Let's construct the
	if err != nil {
		fmt.Println("Error reading your password. Please re-run the command")
		return err
	}

	name := args.First()
	key := args.Get(1)
	// proceed with the construction of the struct
	s := CreateKeyStruct{Password: password, Name: name, Key: key}

	// marshall the json

	jsonBytes, err := json.Marshal(s)

	// the golang http library takes an io.Reader instance for the body, therefore convert this
	body := bytes.NewBuffer(jsonBytes)

	url := BASE + "/app/key/"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("An error occurred while trying to create a reqeust")
		return err
	}
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("An error occurred while trying to talk to the server")
		return err
	}

	defer resp.Body.Close()

	// unmarshall the data from the body of the response
	var response SimpleResponseStruct
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("An error occurred while trying to parse response from the server")
		return err
	}

	if resp.StatusCode == 401 {
		fmt.Println("You have entered the wrong password. Please try again.")
		return errors.New("You have probably entered the wrong password. Please try again")
	}
	if response.Success == false {
		fmt.Print(response.Message)
		return errors.New(response.Message)
	}
	if resp.StatusCode == 201 {
		fmt.Print(response.Message)
	}

	return nil
}

type KeyResponse struct {
	Name string `json:"name"`
}

func ListKeys(_ context.Context, cmd *cli.Command) error {
	// get the args
	// there should be two args
	// the name of the key and the key itself
	// upon finding the absence, throw an error reminding them the correct use of the command
	loggedIn, token := IsLoggedIn()
	if !loggedIn {
		fmt.Println("Please create an account or login to manage your keys")
		return errors.New("Create an account to use this command or login ")
	}

	url := BASE + "/app/key/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("An error occurred while trying to create a reqeust")
		return err
	}
	req.Header.Set("Authorization", "Token "+token)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("An error occurred while trying to talk to the server")
		return err
	}

	defer resp.Body.Close()

	// unmarshall the data from the body of the response
	var response []KeyResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("An error occurred while trying to parse response from the server")
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Println("An unknown error occurred")
	}
	if resp.StatusCode == 200 {
		fmt.Print(response.Message)
	}

	return nil
}
