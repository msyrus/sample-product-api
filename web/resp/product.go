package resp

// Product presents the response object of a product
type Product struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     int     `json:"price"`
	Weight    int     `json:"weight"`
	Available bool    `json:"available"`
	AvgRating float64 `json:"avgRating"`
}
