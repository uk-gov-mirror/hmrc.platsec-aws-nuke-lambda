package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/lambda"
)

type Nuker interface {
	fileExists() bool
	nuke() bool
}
type nukeObject struct {
	filepath string
	dryrun   bool
}

func (no nukeObject) fileExists() bool {
	if _, err := os.Stat(no.filepath); os.IsNotExist(err) {
		log.Printf("File %s does not exist", no.filepath)
		return false
	}
	log.Printf("File %s is found", no.filepath)
	return true
}

var execCommand = exec.Command

func (no nukeObject) nuke() bool {
	args := nukeCmdArgs(no.filepath, no.dryrun)

	log.Printf("args to nuke are: %v", args)

	cmd := execCommand("aws-nuke", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + string(output))
		return false
	}
	log.Println(string(output))
	log.Printf("Output was %s", output)
	return true
}

func nukeCmdArgs(filepath string, dryrun bool) []string {
	args := []string{"--quiet", "--force", "--force-sleep", "3", "--config", filepath}

	if !dryrun {
		args = append(args, "--no-dry-run")
	}

	return args
}

type MyEvent struct {
	ConfigFilename string `json:"ConfigFilename"`
	DryRun         string `json:"DryRun"`
}

type MyResponse struct {
	Message string `json:"Answers"`
}

var runNukeFunction = runNuke

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	dryrun := validateDryRun(event.DryRun)
	nuker := nukeObject{filepath: "/configs/" + event.ConfigFilename, dryrun: dryrun}
	if err := runNukeFunction(nuker); err != nil {
		return MyResponse{}, fmt.Errorf("nuke failed: %s", err.Error())
	}
	return MyResponse{Message: fmt.Sprintf("ConfigFilename is %s and DryRun is %v, the nuke ran", event.ConfigFilename, event.DryRun)}, nil
}

func runNuke(nuker Nuker) error {
	if nuker.fileExists() {
		if nuker.nuke() {
			return nil
		} else {
			return errors.New("nuke did not complete")
		}
	}
	return errors.New("file not found")
}

func validateDryRun(dryrun string) bool {
	if dryrun == "false" {
		log.Printf("DryRun is off, so nuke for real")
		return false
	}
	log.Print("DryRun is on")
	return true
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
