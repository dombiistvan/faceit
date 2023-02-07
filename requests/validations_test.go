package requests

import "testing"

func Test_validatePassword(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "wrong password regexp1",
			args: struct {
				value string
			}{
				value: "testpassword1",
			},
			wantErr: true,
		},
		{
			name: "wrong password regexp2",
			args: struct {
				value string
			}{
				value: "TestPassword1",
			},
			wantErr: true,
		},
		{
			name: "wrong password regexp3",
			args: struct {
				value string
			}{
				value: "TestPassword_",
			},
			wantErr: true,
		},
		{
			name: "wrong password regexp4",
			args: struct {
				value string
			}{
				value: "TESTPASSWORD_1",
			},
			wantErr: true,
		},
		{
			name: "good password regexp",
			args: struct {
				value string
			}{
				value: "Test_Password_1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePassword(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
