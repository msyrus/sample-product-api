package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/msyrus/simple-product-inv/mock_repo"
	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/repo"
	"github.com/msyrus/simple-product-inv/service"
	"github.com/msyrus/simple-product-inv/web/resp"
	uuid "github.com/satori/go.uuid"
)

func injectChiURLParam(req *http.Request, key, val string) {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)

	*req = *req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestNewProductController(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := service.NewRating(rateRepo)
	pdtSvc := service.NewProduct(pdtRepo, rateSvc)

	type args struct {
		svc *service.Product
	}
	tests := []struct {
		name string
		args args
		want *ProductController
	}{
		{
			args: args{
				svc: pdtSvc,
			},
			want: &ProductController{pdtSvc: pdtSvc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProductController(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProductController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductController_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := service.NewRating(rateRepo)
	pdtSvc := service.NewProduct(pdtRepo, rateSvc)

	body1 := bytes.NewBufferString(`{"name": "", "price": 100, "weight": 0, "available": false}`)
	req1, err := http.NewRequest("POST", "/", body1)
	if err != nil {
		t.Fatal(err)
	}

	body2 := bytes.NewBufferString(`"name": "Test1", "price": 100, "weight": 1, "available": false`)
	req2, err := http.NewRequest("POST", "/", body2)
	if err != nil {
		t.Fatal(err)
	}

	body3 := bytes.NewBufferString(`{"name": "Test1", "price": 100, "weight": 1, "available": false}`)
	req3, err := http.NewRequest("POST", "/", body3)
	if err != nil {
		t.Fatal(err)
	}

	body4 := bytes.NewBufferString(`{"name": "Test1", "price": 100, "weight": 1, "available": false}`)
	req4, err := http.NewRequest("POST", "/", body4)
	if err != nil {
		t.Fatal(err)
	}

	uid := uuid.NewV4().String()
	gomock.InOrder(
		pdtRepo.EXPECT().Create(gomock.Any()).Return("", model.ValidationError{}),
		pdtRepo.EXPECT().Create(gomock.Any()).Return("", errors.New("db failed")),
		pdtRepo.EXPECT().Create(gomock.Any()).Return(uid, nil),
	)

	type fields struct {
		pdtSvc *service.Product
	}
	tests := []struct {
		name     string
		r        *http.Request
		fields   fields
		wantCode int
	}{
		{
			r: req1,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			r: req2,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			r: req3,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			r: req4,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			rr := httptest.NewRecorder()
			c.Create(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("ProductController.Create() Code = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestProductController_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := service.NewRating(rateRepo)
	pdtSvc := service.NewProduct(pdtRepo, rateSvc)

	req1, err := http.NewRequest("GET", "/unavailable_id", nil)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req1, "id", "unavailable_id")

	req2, err := http.NewRequest("GET", "/valid_id", nil)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req2, "id", "valid_id")

	gomock.InOrder(
		pdtRepo.EXPECT().Fetch("unavailable_id").Return(nil, nil),
		pdtRepo.EXPECT().Fetch("valid_id").Return(model.Product{ID: "valid_id"}, nil),
		rateRepo.EXPECT().Avg(repo.Query{"product_id": []interface{}{"valid_id"}}, "value").Return(1.5, nil),
		pdtRepo.EXPECT().Fetch("valid_id").Return(nil, errors.New("db failed")),
	)

	type fields struct {
		pdtSvc *service.Product
	}
	tests := []struct {
		name     string
		r        *http.Request
		fields   fields
		wantCode int
	}{
		{
			r: req1,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusNotFound,
		},
		{
			r: req2,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusOK,
		},
		{
			r: req2,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			rr := httptest.NewRecorder()
			c.Get(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("ProductController.Get() Code = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestProductController_List(t *testing.T) {
	type fields struct {
		pdtSvc *service.Product
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			c.List(tt.args.w, tt.args.r)
		})
	}
}

func TestProductController_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := service.NewRating(rateRepo)
	pdtSvc := service.NewProduct(pdtRepo, rateSvc)

	body0 := bytes.NewBufferString(`"name": "Test1", "price": 100, "weight": 1, "available": false`)
	req0, err := http.NewRequest("PUT", "/valid_id", body0)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req0, "id", "valid_id")

	body1 := bytes.NewBufferString(`{"name": "Test2", "price": 200, "weight": 2, "available": true}`)
	req1, err := http.NewRequest("PUT", "/unavailable_id", body1)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req1, "id", "unavailable_id")

	body2 := bytes.NewBufferString(`{"name": "Test2", "price": 200, "weight": 2, "available": true}`)
	req2, err := http.NewRequest("PUT", "/valid_id", body2)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req2, "id", "valid_id")

	body3 := bytes.NewBufferString(`{"price": 100, "weight": 1, "available": true}`)
	req3, err := http.NewRequest("PUT", "/valid_id", body3)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req3, "id", "valid_id")

	body4 := bytes.NewBufferString(`{"id": "no_effect", "price": 100, "weight": 1, "available": true}`)
	req4, err := http.NewRequest("PUT", "/valid_id", body4)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req4, "id", "valid_id")

	gomock.InOrder(
		pdtRepo.EXPECT().Fetch("unavailable_id").Return(nil, nil),
		pdtRepo.EXPECT().Fetch("valid_id").Return(model.Product{ID: "valid_id", Name: "Test1", Price: 100, Weight: 1, Available: false}, nil),
		pdtRepo.EXPECT().Update("valid_id", model.Product{ID: "valid_id", Name: "Test2", Price: 200, Weight: 2, Available: true}).Return(nil),
		pdtRepo.EXPECT().Fetch("valid_id").Return(nil, errors.New("db failed")),
		pdtRepo.EXPECT().Fetch("valid_id").Return(model.Product{ID: "valid_id", Name: "Test2", Price: 200, Weight: 2, Available: true}, nil),
		pdtRepo.EXPECT().Update("valid_id", model.Product{ID: "valid_id", Name: "", Price: 100, Weight: 1, Available: true}).Return(model.ValidationError{}),
	)

	type fields struct {
		pdtSvc *service.Product
	}
	tests := []struct {
		name     string
		r        *http.Request
		fields   fields
		wantCode int
	}{
		{
			r: req0,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			r: req1,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusNotFound,
		},
		{
			r: req2,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusOK,
		},
		{
			r: req3,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			r: req4,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			rr := httptest.NewRecorder()
			c.Update(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("ProductController.Update() Code = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestProductController_UpdatePartial(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := service.NewRating(rateRepo)
	pdtSvc := service.NewProduct(pdtRepo, rateSvc)

	body0 := bytes.NewBufferString(`"name": "Test1", "price": 100, "weight": 1, "available": false`)
	req0, err := http.NewRequest("PATCH", "/valid_id", body0)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req0, "id", "valid_id")

	body1 := bytes.NewBufferString(`{"name": "Test2", "price": 200, "weight": 2, "available": true}`)
	req1, err := http.NewRequest("PATCH", "/unavailable_id", body1)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req1, "id", "unavailable_id")

	body2 := bytes.NewBufferString(`{"name": "Test2", "price": 200, "weight": 2, "available": true}`)
	req2, err := http.NewRequest("PATCH", "/valid_id", body2)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req2, "id", "valid_id")

	body3 := bytes.NewBufferString(`{"price": 100, "weight": 1, "available": true}`)
	req3, err := http.NewRequest("PATCH", "/valid_id", body3)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req3, "id", "valid_id")

	body4 := bytes.NewBufferString(`{"id": "no_effect", "price": 100, "weight": 1, "available": true}`)
	req4, err := http.NewRequest("PATCH", "/valid_id", body4)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req4, "id", "valid_id")

	gomock.InOrder(
		pdtRepo.EXPECT().Fetch("unavailable_id").Return(nil, nil),
		pdtRepo.EXPECT().Fetch("valid_id").Return(model.Product{ID: "valid_id", Name: "Test1", Price: 100, Weight: 1, Available: false}, nil),
		pdtRepo.EXPECT().Update("valid_id", model.Product{ID: "valid_id", Name: "Test2", Price: 200, Weight: 2, Available: true}).Return(nil),
		pdtRepo.EXPECT().Fetch("valid_id").Return(nil, errors.New("db failed")),
		pdtRepo.EXPECT().Fetch("valid_id").Return(model.Product{ID: "valid_id", Name: "Test2", Price: 200, Weight: 2, Available: true}, nil),
		pdtRepo.EXPECT().Update("valid_id", model.Product{ID: "valid_id", Name: "Test2", Price: 100, Weight: 1, Available: true}).Return(nil),
	)

	type fields struct {
		pdtSvc *service.Product
	}
	tests := []struct {
		name     string
		r        *http.Request
		fields   fields
		wantCode int
	}{
		{
			r: req0,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			r: req1,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusNotFound,
		},
		{
			r: req2,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusOK,
		},
		{
			r: req3,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			r: req4,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			rr := httptest.NewRecorder()
			c.UpdatePartial(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("ProductController.UpdatePartial() Code = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestProductController_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := service.NewRating(rateRepo)
	pdtSvc := service.NewProduct(pdtRepo, rateSvc)

	req1, err := http.NewRequest("DELETE", "/unavailable_id", nil)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req1, "id", "unavailable_id")

	req2, err := http.NewRequest("DELETE", "/valid_id", nil)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req2, "id", "valid_id")

	req3, err := http.NewRequest("DELETE", "/valid_id", nil)
	if err != nil {
		t.Fatal(err)
	}
	injectChiURLParam(req3, "id", "valid_id")

	gomock.InOrder(
		pdtRepo.EXPECT().Delete("unavailable_id").Return(nil),
		pdtRepo.EXPECT().Delete("valid_id").Return(errors.New("db failed")),
		pdtRepo.EXPECT().Delete("valid_id").Return(nil),
	)

	type fields struct {
		pdtSvc *service.Product
	}
	tests := []struct {
		name     string
		r        *http.Request
		fields   fields
		wantCode int
	}{
		{
			r: req1,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusOK,
		},
		{
			r: req2,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			r: req3,
			fields: fields{
				pdtSvc: pdtSvc,
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			rr := httptest.NewRecorder()
			c.Delete(rr, tt.r)
			if got := rr.Code; got != tt.wantCode {
				t.Errorf("ProductController.Delete() Code = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestProductController_Rate(t *testing.T) {
	type fields struct {
		pdtSvc *service.Product
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProductController{
				pdtSvc: tt.fields.pdtSvc,
			}
			c.Rate(tt.args.w, tt.args.r)
		})
	}
}

func Test_getSkipLimit(t *testing.T) {
	type args struct {
		r      *http.Request
		dlimit int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getSkipLimit(tt.args.r, tt.args.dlimit)
			if got != tt.want {
				t.Errorf("getSkipLimit() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getSkipLimit() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_toRespProduct(t *testing.T) {
	type args struct {
		pdt    model.Product
		rating float64
	}
	tests := []struct {
		name string
		args args
		want resp.Product
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRespProduct(tt.args.pdt, tt.args.rating); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toRespProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
