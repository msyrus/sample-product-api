package repo

import (
	"reflect"
	"testing"
)

func TestQuery_Add(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name string
		e    Query
		args args
		want Query
	}{
		{
			e: Query{},
			args: args{
				key: "a",
				val: "test",
			},
			want: Query{"a": []interface{}{"test"}},
		},
		{
			e: Query{"a": []interface{}{"test"}},
			args: args{
				key: "a",
				val: "test",
			},
			want: Query{"a": []interface{}{"test", "test"}},
		},
		{
			e: Query{"b": []interface{}{"test"}},
			args: args{
				key: "a",
				val: "test",
			},
			want: Query{"a": []interface{}{"test"}, "b": []interface{}{"test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Add(tt.args.key, tt.args.val)
			if !reflect.DeepEqual(tt.e, tt.want) {
				t.Errorf("Query.Add() Query got = %#v, want = %#v", tt.e, tt.want)
			}
		})
	}
}
