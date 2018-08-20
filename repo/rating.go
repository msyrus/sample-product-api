package repo

import (
	"database/sql"
	"fmt"

	"github.com/msyrus/simple-product-inv/infra"
	"github.com/msyrus/simple-product-inv/model"
	uuid "github.com/satori/go.uuid"
)

// Rating interface is the repo wrapper of rating
type Rating interface {
	Creator
	AvgAggrigator
}

// Critic is an implementation of Rating
type Critic struct {
	table string
	db    infra.DB
}

// NewCritic returns a new Critic with table name tab
func NewCritic(tab string, db infra.DB) *Critic {
	return &Critic{
		table: tab,
		db:    db,
	}
}

// Create creates a new rating in Critic
func (c *Critic) Create(v interface{}) (string, error) {
	rat, ok := v.(model.Rating)
	if !ok {
		return "", ErrUnsupportedType
	}
	rat.ID = uuid.NewV4().String()

	if err := rat.Validate(); err != nil {
		return "", err
	}

	stmt := fmt.Sprintf(`INSERT INTO %s ("id", "product_id", "value") VALUES('%s', '%s', %d)`,
		c.table, rat.ID, rat.ProductID, rat.Value,
	)
	err := c.db.Exec(stmt)
	if err != nil {
		return "", err
	}
	return rat.ID, nil
}

// Avg returns the aggregated average rating value selected by query
func (c *Critic) Avg(q Query, field string) (float64, error) {
	stmt := fmt.Sprintf(`SELECT AVG("%s") FROM %s`, field, c.table)
	vals := []interface{}{}
	if pdtID := q["product_id"]; len(pdtID) != 0 {
		vals = append(vals, pdtID[0])
		stmt = stmt + ` WHERE "product_id" = $1`
	}

	rows, err := c.db.Query(stmt, vals...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, nil
	}
	var f sql.NullFloat64
	err = rows.Scan(&f)
	if err != nil {
		return 0, err
	}
	if !f.Valid {
		return 0, nil
	}
	return f.Float64, nil
}
