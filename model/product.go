package model

import (
	"time"
)

// Product holds the data of a product
type Product struct {
	ID string

	Name      string
	Price     int
	Weight    int
	Available bool

	Deleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Validate checks if the product is valid to store
// it returns nil if there is no error
// otherwise it will return ValidationError
func (r *Product) Validate() error {
	err := ValidationError{}
	if r.ID == "" {
		err.Add("ID", "is required")
	}
	if r.Name == "" {
		err.Add("Name", "is empty")
	}
	if r.Price < 1 {
		err.Add("Price", "is required")
	}
	if r.Weight < 1 {
		err.Add("Weight", "is invalid")
	}

	if len(err) == 0 {
		return nil
	}
	return err
}
