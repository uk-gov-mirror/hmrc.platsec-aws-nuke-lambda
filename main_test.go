package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	//"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os/exec"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) fileExists() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MyMockedObject) nuke() bool {
	args := m.Called()
	return args.Bool(0)
}

type MyMockedNukeObject struct {
	mock.Mock
	filepath string
	dryrun   bool
}

func (no *MyMockedNukeObject) fileExists() bool {
	args := no.Called()
	return args.Bool(0)
}

func (no *MyMockedNukeObject) nuke() bool {
	args := no.Called()
	return args.Bool(0)
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

var mockedExitStatus = 0 // Default to return exit code zero

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(mockedExitStatus)
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestPerformNuke(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("fileExists").Return(true)
	testObj.On("nuke").Return(true)

	got := performNuke(testObj)
	want := true

	assert.Equal(t, want, got, "Call should return true")
	testObj.AssertExpectations(t)
}

func TestPerformNukeFailsOnFileNotFound(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("fileExists").Return(false)

	got := performNuke(testObj)
	want := false

	assert.Equal(t, want, got, "Call should return false")
	testObj.AssertExpectations(t)
}

func TestPerformNukeFailsOnNukeFunc(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("fileExists").Return(true)
	testObj.On("nuke").Return(false)

	got := performNuke(testObj)
	want := false

	assert.Equal(t, want, got, "Call should return false")
	testObj.AssertExpectations(t)
}

func TestValidateDryRunTrue(t *testing.T) {
	got := validateDryRun("true")
	want := true
	assert.Equal(t, want, got, "Call should return true")
}

func TestValidateDryRunEmpty(t *testing.T) {
	got := validateDryRun("")
	want := true
	assert.Equal(t, want, got, "Call should return true")
}

func TestValidateDryRunFalse(t *testing.T) {
	got := validateDryRun("false")
	want := false
	assert.Equal(t, want, got, "Call should return false")
}

func TestFileExists(t *testing.T) {
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, "/meh", []byte("a file"), 0755)

	got := fileExistsOnFilesystem("/meh", appFS)
	want := true
	assert.Equal(t, want, got, "Call should return true")
}

func TestFileDoesntExist(t *testing.T) {
	appFS := afero.NewMemMapFs()

	got := fileExistsOnFilesystem("/meh", appFS)
	want := false
	assert.Equal(t, want, got, "Call should return false")
}

func TestHandleLambdaEventFailed(t *testing.T) {
	event := MyEvent{ConfigFilename: "meh", DryRun: "true"}

	no := nukeObject{}

	got, _ := no.HandleLambdaEvent(event)
	want := MyResponse{Message: "ConfigFilename is meh and DryRun is true, the nuke failed"}
	assert.Equal(t, want, got, "String doesn't match")
}

func TestHandleLambdaEventFailedDryRunFalse(t *testing.T) {
	event := MyEvent{ConfigFilename: "meh", DryRun: "false"}

	no := nukeObject{}
	origCallPerformNuke := callPerformNuke
	defer func() { callPerformNuke = origCallPerformNuke }()
	callPerformNuke = func(ni nukeInterface) bool {
		return true
	}

	got, _ := no.HandleLambdaEvent(event)
	want := MyResponse{Message: "ConfigFilename is meh and DryRun is false, the nuke was successful"}
	assert.Equal(t, want, got, "String doesn't match")
}

func TestNukeSuccess(t *testing.T) {
	no := nukeObject{}
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	got := no.nuke()
	want := true
	assert.Equal(t, want, got, "Boolean doesn't match")
}

// TODO this not working
// func TestNukeFailure(t *testing.T) {
// 	no := nukeObject{}
// 	mockedExitStatus = 1 // cannot seem to override this for the fakeexec command
// 	execCommand = fakeExecCommand
//     defer func() { execCommand = exec.Command }()

// 	got := no.nuke()
// 	want := false
// 	assert.Equal(t, want, got, fmt.Sprintf("got is %v", got))
// }
