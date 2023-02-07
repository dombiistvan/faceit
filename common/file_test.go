package common

import (
	"reflect"
	"testing"
)

func TestGetFileContent(t *testing.T) {
	type args struct {
		fileFullPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "invalid file",
			args: args{
				fileFullPath: "../something.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "file exist",
			args: args{
				fileFullPath: "../config_good_test.yaml",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileContent(tt.args.fileFullPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && tt.want != nil {
				t.Errorf("GetFileContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
