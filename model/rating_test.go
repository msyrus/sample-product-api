package model

import (
	"reflect"
	"testing"
)

func TestRating_Validate(t *testing.T) {
	tests := []struct {
		name string
		r    *Rating
		err  error
	}{
		{
			r: &Rating{},
			err: ValidationError{
				"ID":        []string{"is required"},
				"ProductID": []string{"is empty"},
				"Value":     []string{"is invalid"},
			},
		},
		{
			r: &Rating{ID: "123"},
			err: ValidationError{
				"ProductID": []string{"is empty"},
				"Value":     []string{"is invalid"},
			},
		},
		{
			r: &Rating{ID: "123", ProductID: "987", Value: 0},
			err: ValidationError{
				"Value": []string{"is invalid"},
			},
		},
		{
			r: &Rating{ID: "123", ProductID: "987", Value: 6},
			err: ValidationError{
				"Value": []string{"is invalid"},
			},
		},
		{
			r:   &Rating{ID: "123", ProductID: "987", Value: 3},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Rating.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
