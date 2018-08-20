package repo

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/msyrus/simple-product-inv/infra"
	"github.com/msyrus/simple-product-inv/mock_infra"
	"github.com/msyrus/simple-product-inv/model"
)

func TestNewChef(t *testing.T) {
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
		want *Chef
	}{
		{
			args: args{
				tab: "test",
				db:  db,
			},
			want: &Chef{
				table: "test",
				db:    db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChef(tt.args.tab, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChef_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)
	chf := NewChef("test", db)

	pdt := model.Product{ID: "1", Name: "Test", Price: 100, Weight: 1, Available: false}

	gomock.InOrder(
		db.EXPECT().Exec(gomock.Any()).Return(nil),
		db.EXPECT().Exec(gomock.Any()).Return(sql.ErrConnDone),
	)

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		want    bool
		wantErr bool
	}{
		{
			c: chf,
			args: args{
				v: struct{}{},
			},
			want:    false,
			wantErr: true,
		},
		{
			c: chf,
			args: args{
				v: model.Product{},
			},
			want:    false,
			wantErr: true,
		},
		{
			c: chf,
			args: args{
				v: pdt,
			},
			want:    pdt.ID != "",
			wantErr: false,
		},
		{
			c: chf,
			args: args{
				v: pdt,
			},
			want:    pdt.ID == "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Create(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chef.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != "") != tt.want {
				t.Errorf("Chef.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChef_Fetch(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)
	row := mock_infra.NewMockRow(mockCtrl)
	chf := NewChef("test", db)

	pdt := model.Product{ID: "1", Name: "Test", Price: 100, Weight: 1, Available: false}

	db.EXPECT().Query(fmt.Sprintf(`SELECT * FROM %s WHERE "id"='%s' AND "deleted"=FALSE`, chf.table, pdt.ID)).Return(row, nil)
	row.EXPECT().Next().Return(true)
	row.EXPECT().Scan(gomock.Any()).Return(nil)
	row.EXPECT().Close().Return(nil)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			c: chf,
			args: args{
				id: "1",
			},
			want:    model.Product{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Fetch(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chef.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chef.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChef_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)
	chf := NewChef("test", db)

	pdt := model.Product{ID: "1", Name: "Test", Price: 100, Weight: 1, Available: false}

	gomock.InOrder(
		db.EXPECT().Exec(fmt.Sprintf(`UPDATE %s SET ("name", "price", "weight", "available", "updated_at") = ('%s', %d, %d, %t, CURRENT_TIMESTAMP)
		WHERE "id"='unavailable_id' AND "deleted"=FALSE`, chf.table, pdt.Name, pdt.Price, pdt.Weight, pdt.Available)).Return(nil),
		db.EXPECT().Exec(fmt.Sprintf(`UPDATE %s SET ("name", "price", "weight", "available", "updated_at") = ('%s', %d, %d, %t, CURRENT_TIMESTAMP)
		WHERE "id"='%s' AND "deleted"=FALSE`, chf.table, pdt.Name, pdt.Price, pdt.Weight, pdt.Available, pdt.ID)).Return(nil),
	)

	type args struct {
		id string
		v  interface{}
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		wantErr bool
	}{
		{
			c: chf,
			args: args{
				id: pdt.ID,
				v:  struct{}{},
			},
			wantErr: true,
		},
		{
			c: chf,
			args: args{
				id: pdt.ID,
				v:  model.Product{},
			},
			wantErr: true,
		},
		{
			c: chf,
			args: args{
				id: "unavailable_id",
				v:  pdt,
			},
			wantErr: false,
		},
		{
			c: chf,
			args: args{
				id: pdt.ID,
				v:  pdt,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Update(tt.args.id, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Chef.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChef_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := mock_infra.NewMockDB(mockCtrl)
	chf := NewChef("test", db)

	pdt := model.Product{ID: "1", Name: "Test", Price: 100, Weight: 1, Available: false}

	gomock.InOrder(
		db.EXPECT().Exec(fmt.Sprintf(`UPDATE %s SET ("deleted", "deleted_at") = (TRUE, CURRENT_TIMESTAMP) WHERE "id"='%s' AND "deleted"=FALSE`, chf.table, "unavailable_id")).Return(nil),
		db.EXPECT().Exec(fmt.Sprintf(`UPDATE %s SET ("deleted", "deleted_at") = (TRUE, CURRENT_TIMESTAMP) WHERE "id"='%s' AND "deleted"=FALSE`, chf.table, pdt.ID)).Return(nil),
	)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		wantErr bool
	}{
		{
			c: chf,
			args: args{
				id: "unavailable_id",
			},
			wantErr: false,
		},
		{
			c: chf,
			args: args{
				id: pdt.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Chef.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChef_List(t *testing.T) {
	type args struct {
		skip  int
		limit int
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.List(tt.args.skip, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chef.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chef.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChef_Count(t *testing.T) {
	tests := []struct {
		name    string
		c       *Chef
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Count()
			if (err != nil) != tt.wantErr {
				t.Errorf("Chef.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Chef.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChef_Search(t *testing.T) {
	type args struct {
		q     Query
		skip  int
		limit int
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Search(tt.args.q, tt.args.skip, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chef.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chef.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChef_SearchCount(t *testing.T) {
	type args struct {
		q Query
	}
	tests := []struct {
		name    string
		c       *Chef
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.SearchCount(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chef.SearchCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Chef.SearchCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildProductQuery(t *testing.T) {
	type args struct {
		q Query
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := buildProductQuery(tt.args.q)
			if got != tt.want {
				t.Errorf("buildProductQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildProductQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
