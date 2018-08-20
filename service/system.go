package service

import (
	"github.com/msyrus/simple-product-inv/version"
)

// System is the basic implementation of System
type System struct {
}

// NewSystem returns a new System service
func NewSystem() *System {
	return &System{}
}

// Ready returns if the system is ready to serve
func (s *System) Ready() (bool, error) {
	return true, nil
}

// Health returns system health
func (s *System) Health() (bool, error) {
	return true, nil
}

// Version returns system version
func (s *System) Version() string {
	return version.Version
}
