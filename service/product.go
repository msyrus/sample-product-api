package service

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/msyrus/simple-product-inv/log"
	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/repo"
)

// Product holds fields and dependencies to serve product
type Product struct {
	pdtRepo repo.Product
	olgr    log.Logger
	elgr    log.Logger
	ratSvc  *Rating
}

// ProductOpt represents options for NewProduct
type ProductOpt interface {
	Apply(p *Product)
}

// ProductOptFunc is an implementation of ProductOpt
type ProductOptFunc func(p *Product)

// Apply calls f
func (f ProductOptFunc) Apply(p *Product) {
	f(p)
}

// SetProductOutputLogger sets Product service output logger
func SetProductOutputLogger(l log.Logger) ProductOpt {
	return ProductOptFunc(func(p *Product) {
		if l == nil {
			l = &noOpLogger{}
		}
		p.olgr = l
	})
}

// SetProductErrorLogger sets Product service error logger
func SetProductErrorLogger(l log.Logger) ProductOpt {
	return ProductOptFunc(func(p *Product) {
		if l == nil {
			l = &noOpLogger{}
		}
		p.elgr = l
	})
}

// NewProduct returns a new Product service
func NewProduct(rep repo.Product, rat *Rating, opts ...ProductOpt) *Product {
	r := &Product{
		pdtRepo: rep,
		ratSvc:  rat,
		olgr:    log.DefaultOutputLogger,
		elgr:    log.DefaultErrorLogger,
	}
	for _, opt := range opts {
		opt.Apply(r)
	}
	return r
}

// Add creates a new product
func (p *Product) Add(pdt model.Product) (string, error) {
	p.olgr.Println("creating product", pdt)
	nPdt, err := p.pdtRepo.Create(pdt)
	if err != nil {
		p.elgr.Println("failed to create product", pdt)
		return "", err
	}
	p.olgr.Println("created product", pdt)
	return nPdt, nil
}

// Get returns a model.Product finding by its id
func (p *Product) Get(id string) (*model.Product, error) {
	p.olgr.Println("fetching product by id", id)
	pdtI, err := p.pdtRepo.Fetch(id)
	if err != nil {
		p.elgr.Println("failed to fetch product by id", id, err)
		return nil, err
	}
	if pdtI == nil {
		p.elgr.Println("failed product not found id", id)
		return nil, ErrProductNotFound
	}
	pdt, ok := pdtI.(model.Product)
	if !ok {
		p.elgr.Printf("failed to assert model.Product %#v\n", pdtI)
		return nil, ErrFailedToAssert
	}
	p.olgr.Println("fetched product by id", id)
	return &pdt, nil
}

// Update updates a product finding it with id
func (p *Product) Update(id string, rec model.Product) error {
	p.olgr.Println("updating product by id", id)
	fmt.Printf("%#v\n", rec)
	err := p.pdtRepo.Update(id, rec)
	if err != nil {
		p.elgr.Println("failed to update product by id", id, err)
		return err
	}
	p.olgr.Println("updated product by id", id)
	return nil
}

// Remove deletes a product by its id
func (p *Product) Remove(id string) error {
	p.olgr.Println("deleting product by id", id)
	err := p.pdtRepo.Delete(id)
	if err != nil {
		p.elgr.Println("failed to delete product by id", id, err)
		return err
	}
	p.olgr.Println("deleteed product by id", id)
	return nil
}

// Find returns Products that matches query q with skip and limit
func (p *Product) Find(prms url.Values, skip, limit int) ([]model.Product, error) {
	p.olgr.Println("listing products", prms, skip, limit)
	var res []interface{}
	var err error
	q := buildProductQuery(prms)
	if q == nil {
		res, err = p.pdtRepo.List(skip, limit)
	} else {
		res, err = p.pdtRepo.Search(q, skip, limit)
	}
	if err != nil {
		p.elgr.Println("failed to list products", prms, skip, limit)
		return nil, err
	}
	pdts := []model.Product{}
	for _, re := range res {
		pdt, ok := re.(model.Product)
		if !ok {
			p.elgr.Printf("failed to assert model.Product %#v\n", re)
			return nil, ErrFailedToAssert
		}
		pdts = append(pdts, pdt)
	}
	p.olgr.Println("listed products", prms, skip, limit)
	return pdts, nil
}

// Count returns number of products that matches query q
func (p *Product) Count(prms url.Values) (int, error) {
	p.olgr.Println("counting products", prms)
	var n int
	var err error
	q := buildProductQuery(prms)
	if q == nil {
		n, err = p.pdtRepo.Count()
	} else {
		n, err = p.pdtRepo.SearchCount(q)
	}
	if err != nil {
		p.elgr.Println("failed to count products", prms)
		return 0, err
	}
	p.olgr.Println("counted products", prms)
	return n, nil
}

// Rate rates a product
func (p *Product) Rate(id string, rate int) (string, error) {
	pdt, err := p.Get(id)
	if err != nil {
		return "", err
	}
	rat := model.Rating{
		ProductID: pdt.ID,
		Value:     rate,
	}
	return p.ratSvc.Add(rat)
}

// AvgRating returns average rating of a product by its id
func (p *Product) AvgRating(id string) (float64, error) {
	return p.ratSvc.AvgRating(id)
}

func buildProductQuery(prms url.Values) repo.Query {
	q := repo.Query{}
	for k := range prms {
		if k == "name" {
			if v := prms.Get(k); v != "" {
				q.Add(k, v)
			}
		}
		if k == "weight" {
			if d, err := strconv.Atoi(prms.Get(k)); err == nil && d > 0 {
				q.Add(k, d)
			}
		}
		if k == "available" {
			if b, err := strconv.ParseBool(prms.Get(k)); err == nil {
				q.Add(k, b)
			}
		}
		if k == "price" {
			if d, err := strconv.Atoi(prms.Get(k)); err == nil && d > 0 {
				q.Add(k, d)
			}
		}
	}
	if len(q) == 0 {
		return nil
	}
	return q
}
