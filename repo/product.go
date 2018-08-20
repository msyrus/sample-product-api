package repo

import (
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/msyrus/simple-product-inv/infra"
	"github.com/msyrus/simple-product-inv/model"
)

// Product interface is the repo wrapper of product
type Product interface {
	Creator
	Fetcher
	Updater
	Deleter
	Lister
	Counter
	Searcher
}

// Chef is an implementation of Product interface
type Chef struct {
	table string
	db    infra.DB
}

// NewChef returns new Chef with table name tab
func NewChef(tab string, db infra.DB) *Chef {
	return &Chef{
		table: tab,
		db:    db,
	}
}

// Create a new product
func (c *Chef) Create(v interface{}) (string, error) {
	pdt, ok := v.(model.Product)
	if !ok {
		return "", ErrUnsupportedType
	}
	pdt.ID = uuid.NewV4().String()

	if err := pdt.Validate(); err != nil {
		return "", err
	}

	err := c.db.Exec(fmt.Sprintf(`INSERT INTO %s ("id", "name", "price", "weight", "available") VALUES('%s', '%s', %d, %d, %t)`,
		c.table, pdt.ID, pdt.Name, pdt.Price, pdt.Weight, pdt.Available,
	))
	if err != nil {
		return "", err
	}
	return pdt.ID, nil
}

// Fetch returns a model.Product finding by its id
func (c *Chef) Fetch(id string) (interface{}, error) {
	pdt := model.Product{}

	row, err := c.db.Query(fmt.Sprintf(`SELECT * FROM %s WHERE "id"='%s' AND "deleted"=FALSE`, c.table, id))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, nil
	}
	err = row.Scan(&pdt.ID, &pdt.Name, &pdt.Price, &pdt.Weight, &pdt.Available,
		&pdt.Deleted, &pdt.CreatedAt, &pdt.UpdatedAt, &pdt.DeletedAt)
	if err != nil {
		return nil, err
	}
	return pdt, nil
}

// Update updates a product
func (c *Chef) Update(id string, v interface{}) error {
	pdt, ok := v.(model.Product)
	if !ok {
		return ErrUnsupportedType
	}
	if err := pdt.Validate(); err != nil {
		return err
	}

	stmt := fmt.Sprintf(`UPDATE %s SET ("name", "price", "weight", "available", "updated_at") = ('%s', %d, %d, %t, CURRENT_TIMESTAMP)
		WHERE "id"='%s' AND "deleted"=FALSE`, c.table, pdt.Name, pdt.Price, pdt.Weight, pdt.Available, id)

	return c.db.Exec(stmt)
}

// Delete deletes a product
func (c *Chef) Delete(id string) error {
	return c.db.Exec(fmt.Sprintf(`UPDATE %s SET ("deleted", "deleted_at") = (TRUE, CURRENT_TIMESTAMP) WHERE "id"='%s' AND "deleted"=FALSE`, c.table, id))
}

// List lists products
func (c *Chef) List(skip, limit int) ([]interface{}, error) {
	pdts := []interface{}{}

	rows, err := c.db.Query(fmt.Sprintf(`SELECT * FROM %s WHERE "deleted"=FALSE ORDER BY "created_at" OFFSET %d LIMIT %d`, c.table, skip, limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pdt := model.Product{}
		err = rows.Scan(&pdt.ID, &pdt.Name, &pdt.Price, &pdt.Weight, &pdt.Available,
			&pdt.Deleted, &pdt.CreatedAt, &pdt.UpdatedAt, &pdt.DeletedAt)
		if err != nil {
			return nil, err
		}
		pdts = append(pdts, pdt)
	}

	return pdts, nil
}

// Count counts the number of products
func (c *Chef) Count() (int, error) {
	rows, err := c.db.Query(fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE "deleted"=FALSE`, c.table))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}
	var n int
	err = rows.Scan(&n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Search search products with query
func (c *Chef) Search(q Query, skip, limit int) ([]interface{}, error) {
	qstmt, vals := buildProductQuery(q)
	str := fmt.Sprintf(`SELECT * FROM %s WHERE "deleted"=FALSE `, c.table)
	if len(vals) != 0 {
		str = str + " AND " + qstmt
	}
	str = str + fmt.Sprintf(` ORDER BY "created_at" OFFSET %d LIMIT %d`, skip, limit)

	rows, err := c.db.Query(str, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pdts := []interface{}{}
	for rows.Next() {
		pdt := model.Product{}
		err = rows.Scan(&pdt.ID, &pdt.Name, &pdt.Price, &pdt.Weight, &pdt.Available,
			&pdt.Deleted, &pdt.CreatedAt, &pdt.UpdatedAt, &pdt.DeletedAt)
		if err != nil {
			return nil, err
		}
		pdts = append(pdts, pdt)
	}
	return pdts, nil
}

// SearchCount returns number of products that matches query
func (c *Chef) SearchCount(q Query) (int, error) {
	qstmt, vals := buildProductQuery(q)
	str := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE "deleted"=FALSE `, c.table)
	if len(vals) != 0 {
		str = str + " AND " + qstmt
	}

	rows, err := c.db.Query(str, vals...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}

	var n int
	err = rows.Scan(&n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func buildProductQuery(q Query) (string, []interface{}) {
	str := ""
	vals := []interface{}{}
	cnt := 0
	if name := q["name"]; len(name) != 0 {
		cnt++
		str = fmt.Sprintf(`"name" LIKE $%d`, cnt)
		vals = append(vals, name[0])
	}
	if prc := q["price"]; len(prc) != 0 {
		if cnt != 0 {
			str = str + " AND "
		}
		cnt++
		str = str + fmt.Sprintf(`"price" <= $%d`, cnt)
		vals = append(vals, prc[0])
	}
	if wgt := q["weight"]; len(wgt) != 0 {
		if cnt != 0 {
			str = str + " AND "
		}
		cnt++
		str = str + fmt.Sprintf(`"weight" <= $%d`, cnt)
		vals = append(vals, wgt[0])
	}
	if avl := q["available"]; len(avl) != 0 {
		if cnt != 0 {
			str = str + " AND "
		}
		cnt++
		str = str + fmt.Sprintf(`"available" = $%d`, cnt)
		vals = append(vals, avl[0])
	}
	return str, vals
}
