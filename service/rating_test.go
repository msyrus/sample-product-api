package service

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/msyrus/simple-product-inv/log"
	"github.com/msyrus/simple-product-inv/mock_repo"
	"github.com/msyrus/simple-product-inv/model"
	"github.com/msyrus/simple-product-inv/repo"
)

func TestNewRating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rateRepo := mock_repo.NewMockRating(mockCtrl)

	type args struct {
		rep  repo.Rating
		opts []RatingOpt
	}
	tests := []struct {
		name string
		args args
		want *Rating
	}{
		{
			args: args{
				rep: rateRepo,
			},
			want: &Rating{
				rateRepo: rateRepo,
				olgr:     log.DefaultOutputLogger,
				elgr:     log.DefaultErrorLogger,
			},
		},
		{
			args: args{
				rep: rateRepo,
				opts: []RatingOpt{
					SetRatingOutputLogger(nil),
					SetRatingErrorLogger(log.DefaultOutputLogger),
				},
			},
			want: &Rating{
				rateRepo: rateRepo,
				olgr:     &noOpLogger{},
				elgr:     log.DefaultOutputLogger,
			},
		},
		{
			args: args{
				rep: rateRepo,
				opts: []RatingOpt{
					SetRatingErrorLogger(nil),
					SetRatingOutputLogger(log.DefaultErrorLogger),
				},
			},
			want: &Rating{
				rateRepo: rateRepo,
				olgr:     log.DefaultErrorLogger,
				elgr:     &noOpLogger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRating(tt.args.rep, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRating_Add(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rateRepo := mock_repo.NewMockRating(mockCtrl)

	rate1 := model.Rating{}
	rate2 := model.Rating{ProductID: "1234", Value: 1}

	gomock.InOrder(
		rateRepo.EXPECT().Create(rate1).Return("", model.ValidationError{}),
		rateRepo.EXPECT().Create(rate2).Return("1234", nil),
	)

	type fields struct {
		rateRepo repo.Rating
		olgr     log.Logger
		elgr     log.Logger
	}

	flds := fields{
		rateRepo: rateRepo,
		olgr:     log.DefaultOutputLogger,
		elgr:     log.DefaultErrorLogger,
	}

	type args struct {
		rat model.Rating
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: flds,
			args: args{
				rat: rate1,
			},
			want:    "",
			wantErr: true,
		},
		{
			fields: flds,
			args: args{
				rat: rate2,
			},
			want:    "1234",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rating{
				rateRepo: tt.fields.rateRepo,
				olgr:     tt.fields.olgr,
				elgr:     tt.fields.elgr,
			}
			got, err := r.Add(tt.args.rat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rating.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rating.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRating_AvgRating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rateRepo := mock_repo.NewMockRating(mockCtrl)

	gomock.InOrder(
		rateRepo.EXPECT().Avg(repo.Query{"product_id": []interface{}{"1234"}}, "value").Return(1.5, nil),
	// 	rateRepo.EXPECT().Create(rate2).Return("1234", nil),
	)

	type fields struct {
		rateRepo repo.Rating
		olgr     log.Logger
		elgr     log.Logger
	}
	flds := fields{
		rateRepo: rateRepo,
		olgr:     log.DefaultOutputLogger,
		elgr:     log.DefaultErrorLogger,
	}
	type args struct {
		pdtID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			fields: flds,
			args: args{
				pdtID: "1234",
			},
			want:    1.5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rating{
				rateRepo: tt.fields.rateRepo,
				olgr:     tt.fields.olgr,
				elgr:     tt.fields.elgr,
			}
			got, err := r.AvgRating(tt.args.pdtID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rating.AvgRating() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rating.AvgRating() = %v, want %v", got, tt.want)
			}
		})
	}
}
