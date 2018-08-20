package web

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/service"
	"github.com/msyrus/simple-product-inv/web/resp"
)

var randSrc = rand.NewSource(time.Now().UnixNano())

const alpha = `abcdefghijklmnopqrstuvwxyz` +
	`ABCDEFGHIJKLMNOPQRSTUVWXYZ`
const num = `0123456789`
const charset = alpha + num

func generateErrorID(l int) string {
	r := rand.New(randSrc)
	b := make([]byte, l)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// ServeBadRequest serves http BadRequest
func ServeBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	re := resp.Response{
		Code: http.StatusBadRequest,
		Errors: []resp.Error{
			{
				ID:      generateErrorID(10),
				Message: err.Error(),
			},
		},
	}
	resp.Render(w, r, re)
}

// ServeNotFound serves http NotFound
func ServeNotFound(w http.ResponseWriter, r *http.Request, err error) {
	re := resp.Response{
		Code: http.StatusNotFound,
		Errors: []resp.Error{
			{
				ID:      generateErrorID(10),
				Message: err.Error(),
			},
		},
	}
	resp.Render(w, r, re)
}

// ServeUnprocessableEntity serves http UnprocessableEntity
func ServeUnprocessableEntity(w http.ResponseWriter, r *http.Request, err error, dtl map[string]interface{}) {
	re := resp.Response{
		Code: http.StatusUnprocessableEntity,
		Errors: []resp.Error{
			{
				ID:      generateErrorID(10),
				Message: err.Error(),
				Details: dtl,
			},
		},
	}
	resp.Render(w, r, re)
}

// ServeInternalServerError serves http InternalServerError
func ServeInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	re := resp.Response{
		Code: http.StatusInternalServerError,
		Errors: []resp.Error{
			{
				ID:      generateErrorID(10),
				Message: err.Error(),
			},
		},
	}
	resp.Render(w, r, re)
}

// ServeError serves error with appropriate http status code determined from error type
func ServeError(w http.ResponseWriter, r *http.Request, err error) {
	switch err := err.(type) {
	case model.ValidationError:
		dtl := map[string]interface{}{}
		for k, v := range err {
			dtl[k] = v
		}
		ServeUnprocessableEntity(w, r, err, dtl)
	case service.NotFoundError:
		ServeNotFound(w, r, err)
	default:
		ServeInternalServerError(w, r, err)
	}
}
