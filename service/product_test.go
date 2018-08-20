package service

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/golang/mock/gomock"
	"github.com/msyrus/simple-product-inv/log"
	"github.com/msyrus/simple-product-inv/mock_repo"
	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/repo"
)

func TestNewProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := NewRating(rateRepo)

	type args struct {
		rep  repo.Product
		rat  *Rating
		opts []ProductOpt
	}
	tests := []struct {
		name string
		args args
		want *Product
	}{
		{
			args: args{
				rep: pdtRepo,
				rat: rateSvc,
			},
			want: &Product{
				pdtRepo: pdtRepo,
				ratSvc:  rateSvc,
				olgr:    log.DefaultOutputLogger,
				elgr:    log.DefaultErrorLogger,
			},
		},
		{
			args: args{
				rep: pdtRepo,
				rat: rateSvc,
				opts: []ProductOpt{
					SetProductOutputLogger(nil),
					SetProductErrorLogger(log.DefaultOutputLogger),
				},
			},
			want: &Product{
				pdtRepo: pdtRepo,
				ratSvc:  rateSvc,
				olgr:    &noOpLogger{},
				elgr:    log.DefaultOutputLogger,
			},
		},
		{
			args: args{
				rep: pdtRepo,
				rat: rateSvc,
				opts: []ProductOpt{
					SetProductOutputLogger(log.DefaultErrorLogger),
					SetProductErrorLogger(nil),
				},
			},
			want: &Product{
				pdtRepo: pdtRepo,
				ratSvc:  rateSvc,
				olgr:    log.DefaultErrorLogger,
				elgr:    &noOpLogger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProduct(tt.args.rep, tt.args.rat, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_Add(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := NewRating(rateRepo)
	pdtSvc := NewProduct(pdtRepo, rateSvc)

	rec1 := model.Product{}
	rec2 := model.Product{Name: "Test1", Price: 100, Weight: 1, Available: true}

	uid := uuid.NewV4().String()

	gomock.InOrder(
		pdtRepo.EXPECT().Create(rec1).Return("", model.ValidationError{}),
		pdtRepo.EXPECT().Create(rec2).Return(uid, nil),
	)

	type args struct {
		rec model.Product
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		want    string
		wantErr bool
	}{
		{
			r: pdtSvc,
			args: args{
				rec: rec1,
			},
			want:    "",
			wantErr: true,
		},
		{
			r: pdtSvc,
			args: args{
				rec: rec2,
			},
			want:    uid,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Add(tt.args.rec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Product.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := NewRating(rateRepo)
	pdtSvc := NewProduct(pdtRepo, rateSvc)

	uid := uuid.NewV4().String()
	rec1 := model.Product{ID: uid, Name: "Test1", Price: 100, Weight: 1, Available: true}

	gomock.InOrder(
		pdtRepo.EXPECT().Fetch("not_available_id").Return(nil, nil),
		pdtRepo.EXPECT().Fetch(uid).Return(rec1, nil),
	)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		want    *model.Product
		wantErr bool
	}{
		{
			r: pdtSvc,
			args: args{
				id: "not_available_id",
			},
			want:    nil,
			wantErr: true,
		},
		{
			r: pdtSvc,
			args: args{
				id: uid,
			},
			want:    &rec1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Product.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := NewRating(rateRepo)
	pdtSvc := NewProduct(pdtRepo, rateSvc)

	uid := uuid.NewV4().String()
	rec0 := model.Product{}
	rec1 := model.Product{Name: "Test1", Price: 100, Weight: 1, Available: true}

	gomock.InOrder(
		pdtRepo.EXPECT().Update("not_available_id", rec1).Return(nil),
		pdtRepo.EXPECT().Update(uid, rec0).Return(model.ValidationError{}),
		pdtRepo.EXPECT().Update(uid, rec1).Return(nil),
	)

	type args struct {
		id  string
		rec model.Product
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		wantErr bool
	}{
		{
			r: pdtSvc,
			args: args{
				id:  "not_available_id",
				rec: rec1,
			},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				id:  uid,
				rec: rec0,
			},
			wantErr: true,
		},
		{
			r: pdtSvc,
			args: args{
				id:  uid,
				rec: rec1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Update(tt.args.id, tt.args.rec); (err != nil) != tt.wantErr {
				t.Errorf("Product.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProduct_Remove(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := NewRating(rateRepo)
	pdtSvc := NewProduct(pdtRepo, rateSvc)

	uid := uuid.NewV4().String()

	gomock.InOrder(
		pdtRepo.EXPECT().Delete("not_available_id").Return(nil),
		pdtRepo.EXPECT().Delete(uid).Return(nil),
	)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		wantErr bool
	}{
		{
			r: pdtSvc,
			args: args{
				id: "not_available_id",
			},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				id: uid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Remove(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Product.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProduct_Find(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pdtRepo := mock_repo.NewMockProduct(mockCtrl)
	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rateSvc := NewRating(rateRepo)
	pdtSvc := NewProduct(pdtRepo, rateSvc)

	pdt1 := model.Product{ID: "1", Name: "Test1", Weight: 1, Price: 100, Available: false}
	pdt2 := model.Product{ID: "2", Name: "Test2", Weight: 2, Price: 200, Available: false}
	pdt3 := model.Product{ID: "3", Name: "Test3", Weight: 3, Price: 300, Available: false}

	gomock.InOrder(
		pdtRepo.EXPECT().List(0, 1).Return([]interface{}{pdt1}, nil),
		pdtRepo.EXPECT().Search(repo.Query{"available": []interface{}{false}}, 0, 0).Return([]interface{}{}, nil),
		pdtRepo.EXPECT().Search(repo.Query{"available": []interface{}{false}}, 0, 1).Return([]interface{}{pdt1}, nil),
		pdtRepo.EXPECT().Search(repo.Query{"available": []interface{}{false}}, 1, 2).Return([]interface{}{pdt2, pdt3}, nil),
		pdtRepo.EXPECT().Search(repo.Query{"available": []interface{}{true}}, 0, 1).Return([]interface{}{}, nil),
		pdtRepo.EXPECT().List(0, 5).Return([]interface{}{pdt1, pdt2, pdt3}, nil),
	)

	type args struct {
		prms  url.Values
		skip  int
		limit int
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		want    []model.Product
		wantErr bool
	}{
		{
			r: pdtSvc,
			args: args{
				skip:  0,
				limit: 1,
			},
			want:    []model.Product{pdt1},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				prms:  url.Values{"available": []string{"false"}},
				skip:  0,
				limit: 0,
			},
			want:    []model.Product{},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				prms:  url.Values{"available": []string{"false"}, "invalid_key": []string{"true"}},
				skip:  0,
				limit: 1,
			},
			want:    []model.Product{pdt1},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				prms:  url.Values{"available": []string{"false"}},
				skip:  1,
				limit: 2,
			},
			want:    []model.Product{pdt2, pdt3},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				prms:  url.Values{"available": []string{"true"}},
				skip:  0,
				limit: 1,
			},
			want:    []model.Product{},
			wantErr: false,
		},
		{
			r: pdtSvc,
			args: args{
				prms:  url.Values{"available": []string{"invalid"}},
				skip:  0,
				limit: 5,
			},
			want:    []model.Product{pdt1, pdt2, pdt3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Find(tt.args.prms, tt.args.skip, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Product.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_Count(t *testing.T) {

	type args struct {
		prms url.Values
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Count(tt.args.prms)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Product.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_Rate(t *testing.T) {
	type args struct {
		id   string
		rate int
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Rate(tt.args.id, tt.args.rate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Rate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Product.Rate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_AvgRating(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		r       *Product
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.AvgRating(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.AvgRating() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Product.AvgRating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildProductQuery(t *testing.T) {
	type args struct {
		prms url.Values
	}
	tests := []struct {
		name string
		args args
		want repo.Query
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildProductQuery(tt.args.prms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildProductQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
