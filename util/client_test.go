package util

import (
	"net/http"
	"testing"
)

func Test_checkStatusOk(t *testing.T) {
	type args struct {
		res *http.Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := checkStatusOk(tt.args.res); (err != nil) != tt.wantErr {
			t.Errorf("%q. checkStatusOk() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
