package models

import (
	"reflect"
	"testing"
)

func TestGetStringToPtrStringValue(t *testing.T) {
	type args struct {
		origValue any
		newValue  any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "no new return original",
			args: args{
				origValue: "something",
				newValue:  nil,
			},
			want: "something",
		},
		{
			name: "new value",
			args: args{
				origValue: "something",
				newValue: func() *string {
					var s = "something2"

					return &s
				}(),
			},
			want: "something2",
		},
		{
			name: "new value empty, new value",
			args: args{
				origValue: "something",
				newValue: func() *string {
					var s = ""

					return &s
				}(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringToPtrStringValue(tt.args.origValue, tt.args.newValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStringToPtrStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
