package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Data struct {
	Token string `json:"token"`
}

func WriteConfig(pass string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("An error occurred in getting the user home directory")
		return err
	}
	key := Data{Token: pass}
	filepath := filepath.Join(homeDir, "password-manager.json")
	jsonBytes, err := json.Marshal(key)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = os.WriteFile(filepath, jsonBytes, 0600)
	if err != nil {
		fmt.Print("An error occured while trying to create the config file")
		return err
	}
	return nil
}

func ReadConfig() (error, string) {

	homeDir, err := os.UserHomeDir()

	// read the file
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return err, ""
	}
	filePath := filepath.Join(homeDir, "password-manager.json")
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		// Check if the error is because the file doesn't exist yet
		if os.IsNotExist(err) {
			fmt.Println("Error: Password file does not exist yet. Run the write script first.")
			return err, ""
		}
		fmt.Println("Error reading file:", err)
		return err, ""
	}

	var data Data
	err = json.Unmarshal(jsonBytes, &data)
	return nil, data.Token
}

func IsLoggedIn() (bool, string) {
	err, key := ReadConfig()
	if err != nil {
		fmt.Print(err)
	}
	if key != "" {
		return true, key
	}
	return false, ""
}
