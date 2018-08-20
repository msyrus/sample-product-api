package model

import (
	"reflect"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    ValidationError
		want string
	}{
		{
			e:    ValidationError{},
			want: "invalid data",
		},
		{
			e:    ValidationError{"a": []string{"test"}},
			want: "invalid data",
		},
		{
			e:    ValidationError{"a": []string{"test", "test"}},
			want: "invalid data",
		},
		{
			e:    ValidationError{"a": []string{"test"}, "b": []string{"test"}},
			want: "invalid data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Add(t *testing.T) {
	type args struct {
		key string
		msg string
	}
	tests := []struct {
		name string
		e    ValidationError
		args args
		want ValidationError
	}{
		{
			e: ValidationError{},
			args: args{
				key: "a",
				msg: "test",
			},
			want: ValidationError{"a": []string{"test"}},
		},
		{
			e: ValidationError{"a": []string{"test"}},
			args: args{
				key: "a",
				msg: "test",
			},
			want: ValidationError{"a": []string{"test", "test"}},
		},
		{
			e: ValidationError{"b": []string{"test"}},
			args: args{
				key: "a",
				msg: "test",
			},
			want: ValidationError{"a": []string{"test"}, "b": []string{"test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Add(tt.args.key, tt.args.msg)
			if !reflect.DeepEqual(tt.e, tt.want) {
				t.Errorf("ValidationError.Add() ValidationError got = %#v, want = %#v", tt.e, tt.want)
			}
		})
	}
}
