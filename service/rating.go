package service

import (
	"github.com/msyrus/simple-product-inv/log"
	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/repo"
)

// Rating is a basic implementation of ProductRating service
type Rating struct {
	rateRepo repo.Rating
	olgr     log.Logger
	elgr     log.Logger
}

// RatingOpt represents options for NewRating
type RatingOpt interface {
	Apply(r *Rating)
}

// RatingOptFunc is an implementation of RatingOpt
type RatingOptFunc func(r *Rating)

// Apply calls the f
func (f RatingOptFunc) Apply(r *Rating) {
	f(r)
}

// SetRatingOutputLogger sets Rating service output logger
func SetRatingOutputLogger(l log.Logger) RatingOpt {
	return RatingOptFunc(func(r *Rating) {
		if l == nil {
			l = &noOpLogger{}
		}
		r.olgr = l
	})
}

// SetRatingErrorLogger sets Rating service error logger
func SetRatingErrorLogger(l log.Logger) RatingOpt {
	return RatingOptFunc(func(r *Rating) {
		if l == nil {
			l = &noOpLogger{}
		}
		r.elgr = l
	})
}

// NewRating returns a new Rating service
func NewRating(rep repo.Rating, opts ...RatingOpt) *Rating {
	r := &Rating{
		rateRepo: rep,
		olgr:     log.DefaultOutputLogger,
		elgr:     log.DefaultErrorLogger,
	}
	for _, opt := range opts {
		opt.Apply(r)
	}
	return r
}

// Add creates a new rating
func (r *Rating) Add(rat model.Rating) (string, error) {
	r.olgr.Println("creating rating", rat)
	nRat, err := r.rateRepo.Create(rat)
	if err != nil {
		r.elgr.Println("failed to create rating", rat)
		return "", err
	}
	r.olgr.Println("created rating", rat)
	return nRat, nil
}

// AvgRating returns the average rating of a Product
func (r *Rating) AvgRating(pdtID string) (float64, error) {
	r.olgr.Println("getting avg rating of", pdtID)
	val, err := r.rateRepo.Avg(repo.Query{"product_id": []interface{}{pdtID}}, "value")
	if err != nil {
		r.elgr.Println("failed to get avg rating of", pdtID)
		return 0, err
	}
	r.olgr.Println("got avg rating of", pdtID)
	return val, nil
}
