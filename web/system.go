package web

import (
	"net/http"

	"github.com/msyrus/simple-product-inv/service"
)

// SystemController holds necessary fields to serve system handlers
type SystemController struct {
	sysSvc *service.System
}

// NewSystemController returns new SystemController
func NewSystemController(svc *service.System) *SystemController {
	return &SystemController{
		sysSvc: svc,
	}
}

// Health is the system health check handler
func (c *SystemController) Health(w http.ResponseWriter, r *http.Request) {
	ok, err := c.sysSvc.Health()
	if err != nil {
		ServeError(w, r, err)
		return
	}
	ServeData(w, r, http.StatusOK, ok, nil)
}

// Ready is the system ready check handler
func (c *SystemController) Ready(w http.ResponseWriter, r *http.Request) {
	ok, err := c.sysSvc.Ready()
	if err != nil {
		ServeError(w, r, err)
		return
	}
	ServeData(w, r, http.StatusOK, ok, nil)
}
