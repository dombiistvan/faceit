package common

import "testing"

func TestIsEmptyString(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "empty", args: struct{ value string }{value: ""}, want: true},
		{name: "not empty", args: struct{ value string }{value: "something"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyString(tt.args.value); got != tt.want {
				t.Errorf("IsEmptyString() = %v, want %v", got, tt.want)
			}
		})
	}
}
