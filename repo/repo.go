package repo

// Query represents the query object
type Query map[string][]interface{}

// Add adds new field and value to query
func (q Query) Add(key string, val interface{}) {
	q[key] = append(q[key], val)
}

// Creator interface holds the necessery dependencies to create a entry in repo
// Create takes a model and returns its id on success
type Creator interface {
	Create(v interface{}) (string, error)
}

// Fetcher interface holds the necessery dependencies to Fetch a entry from repo
// Fetch takes an id and returns the model on success
type Fetcher interface {
	Fetch(id string) (interface{}, error)
}

// Updater interface holds the necessery dependencies to Update a entry in repo
// Update takes an id and model, it select the entry by id and replace its fields
// with the new model except the id
type Updater interface {
	Update(id string, v interface{}) error
}

// Deleter interface holds the necessery dependencies to delete a entry in repo
// Delete takes an id and delete the entry
type Deleter interface {
	Delete(id string) error
}

// Searcher interface holds the necessery dependencies to search and count entries
// with the matching query parameters
type Searcher interface {
	Search(q Query, skip, limit int) ([]interface{}, error)
	SearchCount(q Query) (int, error)
}

// Lister interface holds the necessery dependencies to list paginated entries
type Lister interface {
	List(skip, limit int) ([]interface{}, error)
}

// Counter interface holds the necessery dependencies to count the entries in repo
type Counter interface {
	Count() (int, error)
}

// AvgAggrigator interface holds the necessery dependencies to aggregate
// average value of the field with matching query q
type AvgAggrigator interface {
	Avg(q Query, field string) (float64, error)
}
