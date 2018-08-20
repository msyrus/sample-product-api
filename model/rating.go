package model

import (
	"time"
)

// Rating holds the data of a product rating
type Rating struct {
	ID string

	ProductID string
	Value     int

	CreatedAt time.Time
}

// Validate checks if the rating is valid to store
// it returns nil if there is no error
// otherwise it will return ValidationError
func (r *Rating) Validate() error {
	err := ValidationError{}
	if r.ID == "" {
		err.Add("ID", "is required")
	}
	if r.ProductID == "" {
		err.Add("ProductID", "is empty")
	}
	if r.Value < 1 || r.Value > 5 {
		err.Add("Value", "is invalid")
	}

	if len(err) == 0 {
		return nil
	}
	return err
}
