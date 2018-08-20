package model

import (
	"reflect"
	"testing"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name string
		r    *Product
		err  error
	}{
		{
			r: &Product{},
			err: ValidationError{
				"ID":     []string{"is required"},
				"Name":   []string{"is empty"},
				"Weight": []string{"is invalid"},
				"Price":  []string{"is required"},
			},
		},
		{
			r: &Product{
				ID:     "123",
				Name:   "Test1",
				Weight: 0,
				Price:  0,
			},
			err: ValidationError{
				"Weight": []string{"is invalid"},
				"Price":  []string{"is required"},
			},
		},
		{
			r: &Product{
				ID:     "123",
				Name:   "Test1",
				Weight: -1,
				Price:  100,
			},
			err: ValidationError{
				"Weight": []string{"is invalid"},
			},
		},
		{
			r: &Product{
				ID:     "123",
				Name:   "Test1",
				Weight: 3,
				Price:  100,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Product.Validate() error = %#v, err %v", err, tt.err)
			}
		})
	}
}
