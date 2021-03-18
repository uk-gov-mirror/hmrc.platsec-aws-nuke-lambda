package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockNuke struct {
	filepath string
	dryrun   bool
}

func (m MockNuke) fileExists() bool {
	return true
}

func (m MockNuke) nuke() bool {
	return false
}

type MockNukeAllSuccess struct {
	filepath string
	dryrun   bool
}

func (m MockNukeAllSuccess) fileExists() bool {
	return true
}
func (m MockNukeAllSuccess) nuke() bool {
	return true
}

var mockedExitStatus int

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "EXIT_STATUS=" + es}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	es, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(es)
}

func mockRunNuke(nuker Nuker) error {
	return nil
}

func mockRunNukeError(nuker Nuker) error {
	return errors.New("Nuke did not complete")
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func Test_runNuke(t *testing.T) {
	type args struct {
		nuker Nuker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{name: "FileNotFoundCheck", args: args{nuker: nukeObject{filepath: "dsfasd"}}, wantErr: true, errMsg: "File not found"},
		{name: "FileFound", args: args{nuker: MockNuke{filepath: "/tmp/meh"}}, wantErr: true, errMsg: "Nuke did not complete"},
		{name: "NukeSuccess", args: args{nuker: MockNukeAllSuccess{}}, wantErr: false, errMsg: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runNuke(tt.args.nuker)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" {
				assert.Equal(t, tt.errMsg, err.Error(), "Error msg needs to match")
			}
		})
	}
}

func Test_validateDryRun(t *testing.T) {
	type args struct {
		dryrun string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "DryRunFalse", args: args{dryrun: "false"}, want: false},
		{name: "DryRunTrue", args: args{dryrun: "true"}, want: true},
		{name: "DryRunBlank", args: args{dryrun: ""}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateDryRun(tt.args.dryrun); got != tt.want {
				t.Errorf("validateDryRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nukeObject_nuke(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	type fields struct {
		filepath string
		dryrun   bool
	}
	tests := []struct {
		name       string
		fields     fields
		exitStatus int
		want       bool
	}{
		{name: "NukeFailed", fields: fields{filepath: "/blah", dryrun: true}, exitStatus: 1, want: false},
		{name: "NukeSuccess", fields: fields{filepath: "/blah", dryrun: true}, exitStatus: 0, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedExitStatus = tt.exitStatus
			no := nukeObject{
				filepath: tt.fields.filepath,
				dryrun:   tt.fields.dryrun,
			}
			if got := no.nuke(); got != tt.want {
				t.Errorf("nukeObject.nuke() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nukeObject_fileExists(t *testing.T) {
	os.Create("/tmp/mah2")
	defer os.Remove("/tmp/mah2")

	type fields struct {
		filepath string
		dryrun   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "FileExistsFail", fields: fields{filepath: "mah", dryrun: true}, want: false},
		{name: "FileExistsSuccess", fields: fields{filepath: "/tmp/mah2", dryrun: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := nukeObject{
				filepath: tt.fields.filepath,
				dryrun:   tt.fields.dryrun,
			}
			if got := no.fileExists(); got != tt.want {
				t.Errorf("nukeObject.fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleLambdaEvent(t *testing.T) {
	//runNukeFunction = mockRunNuke
	//mockRunNukeError
	defer func() { runNukeFunction = runNuke }()

	type args struct {
		event MyEvent
	}
	tests := []struct {
		name    string
		args    args
		want    MyResponse
		wantErr bool
		errMsg  string
	}{
		{name: "Success", args: args{MyEvent{ConfigFilename: "meh", DryRun: "true"}}, want: MyResponse{Message: "ConfigFilename is meh and DryRun is true, the nuke ran"}, wantErr: false, errMsg: ""},
		{name: "Failure", args: args{MyEvent{ConfigFilename: "meh", DryRun: "true"}}, want: MyResponse{}, wantErr: true, errMsg: "Nuke failed: Nuke did not complete"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runNukeFunction = mockRunNuke
			if tt.wantErr == true {
				runNukeFunction = mockRunNukeError
			}

			got, err := HandleLambdaEvent(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLambdaEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errMsg != "" {
				assert.Equal(t, tt.errMsg, err.Error(), "Error msg needs to match")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLambdaEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
