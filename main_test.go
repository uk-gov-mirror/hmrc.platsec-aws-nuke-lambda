package main

import (
	"testing"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct{
	mock.Mock
}

func (m *MyMockedObject) fileExists() (bool) {
	args := m.Called()
	return args.Bool(0)
}

func (m *MyMockedObject) nuke() (bool) {
	args := m.Called()
	return args.Bool(0)
}

type MyMockedNukeObject struct{
	mock.Mock
	filepath	string
	dryrun		bool
}

func (no *MyMockedNukeObject) fileExists() (bool) {
	args := no.Called()
	return args.Bool(0)
}

func (no *MyMockedNukeObject) nuke() (bool) {
	args := no.Called()
	return args.Bool(0)
}

func TestPerformNuke(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("fileExists").Return(true)
	testObj.On("nuke").Return(true)

	got := performNuke(testObj)
	want := true

	assert.Equal(t, got, want, "Call should return true")
	testObj.AssertExpectations(t)
}

func TestPerformNukeFailsOnFileNotFound(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("fileExists").Return(false)

	got := performNuke(testObj)
	want := false

	assert.Equal(t, got, want, "Call should return false")
	testObj.AssertExpectations(t)
}

func TestPerformNukeFailsOnNukeFunc(t *testing.T) {
	testObj := new(MyMockedObject)
	testObj.On("fileExists").Return(true)
	testObj.On("nuke").Return(false)

	got := performNuke(testObj)
	want := false

	assert.Equal(t, got, want, "Call should return false")
	testObj.AssertExpectations(t)
}

func TestValidateDryRunTrue(t *testing.T) {
	got := validateDryRun("true")
	want := true
	assert.Equal(t, got, want, "Call should return true")
}

func TestValidateDryRunEmpty(t *testing.T) {
	got := validateDryRun("")
	want := true
	assert.Equal(t, got, want, "Call should return true")
}

func TestValidateDryRunFalse(t *testing.T) {
	got := validateDryRun("false")
	want := false
	assert.Equal(t, got, want, "Call should return false")
}

func TestFileExists(t *testing.T) {
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, "/meh", []byte("a file"), 0755)

	got := fileExistsOnFilesystem("/meh", appFS)
	want := true
	assert.Equal(t, got, want, "Call should return true")
}

func TestFileDoesntExist(t *testing.T) {
	appFS := afero.NewMemMapFs()

	got := fileExistsOnFilesystem("/meh", appFS)
	want := false
	assert.Equal(t, got, want, "Call should return false")
}

func TestHandleLambdaEventFailed(t *testing.T) {
	event := MyEvent{ConfigFilename: "meh", DryRun: "true"}

	no := nukeObject{}

	got, _ := no.HandleLambdaEvent(event)
	want := MyResponse{Message: "ConfigFilename is meh and DryRun is true, the nuke failed"}
	assert.Equal(t, got, want, "String doesn't match")
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
	assert.Equal(t, got, want, "String doesn't match")
}