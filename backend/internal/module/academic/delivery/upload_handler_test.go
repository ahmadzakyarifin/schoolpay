package delivery

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewUploadHandler(t *testing.T) {
	tests := []struct {
		name string
		want *UploadHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUploadHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUploadHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUploadHandler_UploadStudentPhoto(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *UploadHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UploadHandler{}
			h.UploadStudentPhoto(tt.args.c)
		})
	}
}
