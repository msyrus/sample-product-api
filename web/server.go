package web

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"

	"github.com/msyrus/simple-product-inv/log"
	"github.com/msyrus/simple-product-inv/web/middleware"
	"github.com/msyrus/simple-product-inv/web/resp"
)

// NewRouter returns a http.Handler with all API registered
func NewRouter(pdtCtrl *ProductController, sysCtl *SystemController) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recover)
	router.Use(middleware.Logger(log.DefaultOutputLogger))

	router.NotFound(NotFoundHandler)
	router.MethodNotAllowed(MethodNotAllowed)

	router.Route("/", func(r chi.Router) {
		r.Mount("/products", productHandlers(pdtCtrl))
		r.Mount("/system", systemHandlers(sysCtl))
		r.Mount("/debug", debugHandlers())
	})

	return router
}

// ServeData serves data and meta with http status code 2xx
func ServeData(w http.ResponseWriter, r *http.Request, code int, data interface{}, meta *resp.Pager) {
	if code < 200 || code > 299 {
		panic(fmt.Errorf("serve data with %d", code))
	}
	re := resp.Response{
		Code: code,
		Data: data,
		Meta: meta,
	}
	resp.Render(w, r, re)
}

// NotFoundHandler handles when no routes match
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// MethodNotAllowed handles when no routes match
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func productHandlers(ctrl *ProductController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/", ctrl.List)
		r.With(middleware.Auth).Post("/", ctrl.Create)
		r.Get("/{id}", ctrl.Get)
		r.With(middleware.Auth).Put("/{id}", ctrl.Update)
		r.With(middleware.Auth).Patch("/{id}", ctrl.UpdatePartial)
		r.With(middleware.Auth).Delete("/{id}", ctrl.Delete)
		r.Post("/{id}/rating", ctrl.Rate)
	})
	return h
}

// svc := service.NewProduct()
// 	ctrl := NewProductController(svc)

func systemHandlers(ctrl *SystemController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/health", ctrl.Health)
		r.Get("/ready", ctrl.Ready)
	})
	return h
}

func debugHandlers() http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/pprof", pprof.Index)
		r.Get("/pprof/cmdline", pprof.Cmdline)
		r.Get("/pprof/profile", pprof.Profile)
		r.Get("/pprof/profile", pprof.Profile)
		r.Get("/pprof/symbol", pprof.Symbol)
		r.Get("/pprof/trace", pprof.Trace)
	})
	return h
}
