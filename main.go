package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"github.com/spf13/afero"
        "os"
	"os/exec"
)

type nukeInterface interface {
        fileExists()            bool
        nuke()                  bool
}

type nukeObject struct {
        filepath        string
        dryrun          bool
}

func (no nukeObject) fileExists() bool {
        filesystem := afero.NewOsFs()
	return fileExistsOnFilesystem(no.filepath, filesystem)
}

func (no nukeObject) nuke() bool {
        args := []string{"--quiet", "--force", "--force-sleep", "3", "--config", no.filepath}
	// args = append(args, "--std=c++11")
	
	cmd := exec.Command("aws-nuke", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return false
	}
	fmt.Println(string(output))
	fmt.Printf("Output was %s", output)
        return true
}

func fileExistsOnFilesystem(filepath string, filesystem afero.Fs) bool {
	if _, err := filesystem.Stat(filepath); os.IsNotExist(err) {
                log.Printf("File %s does not exist", filepath)
		return false
	}
        log.Printf("File %s is found", filepath)
	return true
}

func validateDryRun(dryrun string) bool {
	if dryrun == "false" {
		log.Printf("DryRun is off, so nuke for real")
		return false
	}
	log.Print("DryRun is on")
	return true
}

type MyEvent struct {
	ConfigFilename string `json:"ConfigFilename"`
	DryRun         string `json:"DryRun"`
}

type MyResponse struct {
	Message string `json:"Answers"`
}

func performNuke(ni nukeInterface) bool {
        if ni.fileExists() {
                return ni.nuke()
        }
	return false
        	
}

var callPerformNuke = performNuke

func (no *nukeObject) HandleLambdaEvent(event MyEvent) (MyResponse, error) {
        dryrun := validateDryRun(event.DryRun)
        no.filepath = "/configs/" + event.ConfigFilename
        no.dryrun = dryrun
        //no := nukeObject{filepath: "/configs/" + event.ConfigFilename, dryrun: dryrun}
        nukeSuccess := "failed"
	if callPerformNuke(no) { 
                nukeSuccess = "was successful" 
        }

        return MyResponse{Message: fmt.Sprintf("ConfigFilename is %s and DryRun is %v, the nuke %s", event.ConfigFilename, event.DryRun, nukeSuccess)}, nil
}



func main() {

        no := nukeObject{
                filepath: "",
                dryrun: true,
        }

	lambda.Start(no.HandleLambdaEvent)
}
