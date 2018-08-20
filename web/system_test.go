package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/msyrus/simple-product-inv/service"
)

func TestNewSystemController(t *testing.T) {
	sysSvc := service.NewSystem()

	type args struct {
		svc *service.System
	}
	tests := []struct {
		name string
		args args
		want *SystemController
	}{
		{
			args: args{
				svc: sysSvc,
			},
			want: &SystemController{sysSvc: sysSvc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystemController(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemController_Health(t *testing.T) {
	sysSvc := service.NewSystem()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		sysSvc *service.System
	}
	tests := []struct {
		name     string
		fields   fields
		r        *http.Request
		wantCode int
		wantBody []byte
	}{
		{
			fields: fields{
				sysSvc: sysSvc,
			},
			r:        req,
			wantCode: http.StatusOK,
			wantBody: []byte(`{"data":true}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SystemController{
				sysSvc: tt.fields.sysSvc,
			}
			rr := httptest.NewRecorder()
			c.Health(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("SystemController.Health() Code = %v, want %v", got, tt.wantCode)
			}
			if got := rr.Body.Bytes(); bytes.Compare(got, tt.wantBody) != 0 {
				t.Errorf("SystemController.Health() Body = %v, want %v", string(got), string(tt.wantBody))
			}
		})
	}
}

func TestSystemController_Ready(t *testing.T) {
	sysSvc := service.NewSystem()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		sysSvc *service.System
	}
	tests := []struct {
		name     string
		fields   fields
		r        *http.Request
		wantCode int
		wantBody []byte
	}{
		{
			fields: fields{
				sysSvc: sysSvc,
			},
			r:        req,
			wantCode: http.StatusOK,
			wantBody: []byte(`{"data":true}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SystemController{
				sysSvc: tt.fields.sysSvc,
			}
			rr := httptest.NewRecorder()
			c.Ready(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("SystemController.Ready() Code = %v, want %v", got, tt.wantCode)
			}
			if got := rr.Body.Bytes(); bytes.Compare(got, tt.wantBody) != 0 {
				t.Errorf("SystemController.Ready() Body = %v, want %v", string(got), string(tt.wantBody))
			}
		})
	}
}
