package web

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/service"
	"github.com/msyrus/simple-product-inv/web/resp"
)

// ProductController holds necessary fields to serve product handlers
type ProductController struct {
	pdtSvc *service.Product
}

// NewProductController returns a new ProductController with the svc
func NewProductController(svc *service.Product) *ProductController {
	return &ProductController{
		pdtSvc: svc,
	}
}

func parseJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

type createProductBody struct {
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Weight    int    `json:"weight"`
	Available bool   `json:"available"`
}

// Create is the product create handler
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	body := createProductBody{}
	if err := parseJSON(r.Body, &body); err != nil {
		ServeBadRequest(w, r, err)
		return
	}

	pdt := model.Product{
		Name:      body.Name,
		Price:     body.Price,
		Weight:    body.Weight,
		Available: body.Available,
	}
	rID, err := c.pdtSvc.Add(pdt)
	if err != nil {
		ServeError(w, r, err)
		return
	}

	ServeData(w, r, http.StatusCreated, rID, nil)
	return
}

// Get serves a product with its id from url param {id}
func (c *ProductController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pdt, err := c.pdtSvc.Get(id)
	if err != nil {
		ServeError(w, r, err)
		return
	}
	rt, err := c.pdtSvc.AvgRating(pdt.ID)
	if err != nil {
		ServeError(w, r, err)
		return
	}
	ServeData(w, r, http.StatusOK, toRespProduct(*pdt, rt), nil)
	return
}

// List serves a list of products
// it also filters with query params
func (c *ProductController) List(w http.ResponseWriter, r *http.Request) {
	skip, limit := getSkipLimit(r, 20)
	prms := r.URL.Query()

	n, err := c.pdtSvc.Count(prms)
	if err != nil {
		ServeError(w, r, err)
		return
	}

	pgr := resp.NewPager(n, skip, limit)

	if n <= skip {
		ServeData(w, r, http.StatusOK, []struct{}{}, pgr)
		return
	}

	pdts, err := c.pdtSvc.Find(prms, skip, limit)
	if err != nil {
		ServeError(w, r, err)
		return
	}

	rs := []resp.Product{}
	for _, pdt := range pdts {
		rt, err := c.pdtSvc.AvgRating(pdt.ID)
		if err != nil {
			ServeError(w, r, err)
			return
		}
		rs = append(rs, toRespProduct(pdt, rt))
	}

	ServeData(w, r, http.StatusOK, rs, pgr)
	return
}

type updateProductBody struct {
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Weight    int    `json:"weight"`
	Available bool   `json:"available"`
}

// Update updates a product finding it with its id from url param {id}
func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	body := updateProductBody{}
	if err := parseJSON(r.Body, &body); err != nil {
		ServeBadRequest(w, r, err)
		return
	}

	id := chi.URLParam(r, "id")
	pdt, err := c.pdtSvc.Get(id)
	if err != nil {
		ServeError(w, r, err)
		return
	}

	pdt.Name = body.Name
	pdt.Price = body.Price
	pdt.Weight = body.Weight
	pdt.Available = body.Available
	if err := c.pdtSvc.Update(id, *pdt); err != nil {
		ServeError(w, r, err)
		return
	}

	ServeData(w, r, http.StatusOK, pdt.ID, nil)
}

type updatePartProductBody struct {
	Name      *string `json:"name"`
	Price     *int    `json:"price"`
	Weight    *int    `json:"weight"`
	Available *bool   `json:"available"`
}

// UpdatePartial updates a product partially with request body finding it with its id from url param {id}
func (c *ProductController) UpdatePartial(w http.ResponseWriter, r *http.Request) {
	body := updatePartProductBody{}
	if err := parseJSON(r.Body, &body); err != nil {
		ServeBadRequest(w, r, err)
		return
	}

	id := chi.URLParam(r, "id")
	pdt, err := c.pdtSvc.Get(id)
	if err != nil {
		ServeError(w, r, err)
		return
	}
	if body.Name != nil {
		pdt.Name = *body.Name
	}
	if body.Price != nil {
		pdt.Price = *body.Price
	}
	if body.Weight != nil {
		pdt.Weight = *body.Weight
	}
	if body.Available != nil {
		pdt.Available = *body.Available
	}

	if err := c.pdtSvc.Update(id, *pdt); err != nil {
		ServeError(w, r, err)
		return
	}

	ServeData(w, r, http.StatusOK, pdt.ID, nil)
}

// Delete deletes a product with its id from url param {id}
func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := c.pdtSvc.Remove(id)
	if err != nil {
		if _, ok := err.(service.NotFoundError); !ok {
			ServeError(w, r, err)
			return
		}
	}
	ServeData(w, r, http.StatusOK, true, nil)
	return
}

type rateProductBody struct {
	Value int `json:"value"`
}

// Rate rates a product with its id from url param {id}
func (c *ProductController) Rate(w http.ResponseWriter, r *http.Request) {
	body := rateProductBody{}
	if err := parseJSON(r.Body, &body); err != nil {
		ServeBadRequest(w, r, err)
		return
	}

	id := chi.URLParam(r, "id")
	rID, err := c.pdtSvc.Rate(id, body.Value)
	if err != nil {
		ServeError(w, r, err)
		return
	}
	ServeData(w, r, http.StatusCreated, rID, nil)
	return
}

func getSkipLimit(r *http.Request, dlimit int) (int, int) {
	q := r.URL.Query()
	skip, _ := strconv.Atoi(q.Get("skip"))
	if skip < 0 {
		skip = 0
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 0 {
		limit = 0
	}
	if limit == 0 {
		limit = dlimit
	}
	return skip, limit
}

func toRespProduct(pdt model.Product, rating float64) resp.Product {
	return resp.Product{
		ID:        pdt.ID,
		Name:      pdt.Name,
		Price:     pdt.Price,
		Weight:    pdt.Weight,
		Available: pdt.Available,
		AvgRating: rating,
	}
}
