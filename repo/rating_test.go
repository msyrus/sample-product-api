package repo

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/msyrus/simple-product-inv/infra"
	"github.com/msyrus/simple-product-inv/mock_infra"
	"github.com/msyrus/simple-product-inv/model"
)

func TestNewCritic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)

	type args struct {
		tab string
		db  infra.DB
	}
	tests := []struct {
		name string
		args args
		want *Critic
	}{
		{
			args: args{
				db:  db,
				tab: "test",
			},
			want: &Critic{table: "test", db: db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCritic(tt.args.tab, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCritic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCritic_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)
	ctc := NewCritic("test", db)

	rat := model.Rating{ID: "1", ProductID: "111", Value: 3, CreatedAt: time.Time{}}

	gomock.InOrder(
		db.EXPECT().Exec(gomock.Any()).Return(nil),
		db.EXPECT().Exec(gomock.Any()).Return(sql.ErrConnDone),
	)

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *Critic
		args    args
		want    bool
		wantErr bool
	}{
		{
			c: ctc,
			args: args{
				v: struct{}{},
			},
			want:    false,
			wantErr: true,
		},
		{
			c: ctc,
			args: args{
				v: model.Rating{},
			},
			want:    false,
			wantErr: true,
		},
		{
			c: ctc,
			args: args{
				v: rat,
			},
			want:    rat.ID != "",
			wantErr: false,
		},
		{
			c: ctc,
			args: args{
				v: rat,
			},
			want:    rat.ID == "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Create(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Critic.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != "") != tt.want {
				t.Errorf("Critic.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCritic_Avg(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)
	ctc := NewCritic("test", db)

	rat := model.Rating{ID: "1", ProductID: "111", Value: 3, CreatedAt: time.Time{}}

	var f sql.NullFloat64
	row := mock_infra.NewMockRow(mockCtrl)
	row.EXPECT().Next().Return(true)
	row.EXPECT().Scan(gomock.AssignableToTypeOf(&f)).Return(nil)
	row.EXPECT().Close().Return(nil)

	db.EXPECT().Query(fmt.Sprintf(`SELECT AVG("value") FROM %s WHERE "product_id" = $1`, ctc.table), rat.ID).Return(row, nil)

	type args struct {
		q     Query
		field string
	}
	tests := []struct {
		name    string
		c       *Critic
		args    args
		want    float64
		wantErr bool
	}{
		{
			c: ctc,
			args: args{
				q:     Query{"product_id": []interface{}{rat.ID}},
				field: "value",
			},
			want:    f.Float64,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Avg(tt.args.q, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Critic.Avg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Critic.Avg() = %v, want %v", got, tt.want)
			}
		})
	}
}
