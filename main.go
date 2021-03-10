package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"os/exec"
)

type MyEvent struct {
	ConfigFilename string `json:"ConfigFilename"`
	Profile        string `json:"Profile"`
}

type MyResponse struct {
	Message string `json:"Answers"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	configPath := "/configs/" + event.ConfigFilename
	fileExists := fileExists(configPath)
	// validateProfile(event.Profile)

	nuke(configPath)

	return MyResponse{Message: fmt.Sprintf("ConfigFilename is %s and Profile is %s Does file exist? %t", event.ConfigFilename, event.Profile, fileExists)}, nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File %s in /configs does not exist", filename)
		return false
	}
	log.Printf("File %s is found", filename)
	return true
}

func validateProfile(profile string) bool {
	if profile == "" {
		log.Fatalf("Profile is not set")
		return false
	}
	log.Printf("Profile %s is found", profile)
	return true
}

func nuke(path string) {
        args := []string{"--quiet", "--force", "--force-sleep", "3", "--config", path}
	// args = append(args, "--std=c++11")
	
	cmd := exec.Command("aws-nuke", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	fmt.Printf("Output was %s", output)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
