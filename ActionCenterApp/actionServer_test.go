package actioncenterapp

import (
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name string
		want *ActionServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunServer()
		})
	}
}

func TestActionServer_init(t *testing.T) {
	tests := []struct {
		name string
		s    *ActionServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.init()
		})
	}
}

func TestActionServer_Close(t *testing.T) {
	tests := []struct {
		name string
		s    *ActionServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Close()
		})
	}
}
