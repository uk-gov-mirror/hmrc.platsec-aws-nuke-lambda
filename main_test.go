package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MyMockedNukeObject struct {
	mock.Mock
}

func (m *MyMockedNukeObject) fileExists() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MyMockedNukeObject) nuke() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

/*func TestPerformNuke(t *testing.T) {
	testObj := new(MyMockedNukeObject)
	testObj.On("fileExists").Return(true)
	testObj.On("nuke").Return(true)

	got := run(testObj)
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
*/

func Test_run(t *testing.T) {
	// Setup test files
	tmpfile, _ := ioutil.TempFile("/tmp", "meh")
	defer os.Remove(tmpfile.Name())

	testObj := new(MyMockedNukeObject)

	type args struct {
		nuker nukeObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{name: "FileNotFoundCheck", args: args{}, wantErr: true, errMsg: "File not found"},
		{name: "FileFound", args: args{nuker: nukeObject{filepath: "/tmp/meh"}}, wantErr: true, errMsg: ""},
		{name: "NukeSuccess", args: args{nuker: testObj}, wantErr: true, errMsg: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args.nuker)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" {
				assert.Equal(t, tt.errMsg, err.Error(), "Error msg needs to match")
			}
		})
	}
}
