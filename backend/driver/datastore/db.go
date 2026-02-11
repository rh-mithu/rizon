package datastore

import (
	"context"
)

// SortOrder defines sorting direction
type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

// QueryOption defines filtering, sorting, pagination, joins, and aggregation
type QueryOption struct {
	Filter map[string]interface{} // e.g. {"age >": 30, "status": "active"}
	Sort   map[string]SortOrder   // e.g. {"created_at": Desc}
	Limit  int64
	Skip   int64
	Select []string // projection fields (e.g. ["id","name"])

	// Advanced relational options
	Join      []JoinOption      // e.g. {Type: "INNER", Table: "orders o", On: "u.id = o.user_id"}
	Group     []string          // e.g. {"u.id"}
	Having    map[string]string // e.g. {"SUM(o.amount) >": "100"}
	Distinct  bool
	Relations []string
}

type JoinOption struct {
	Type   string   // LEFT, INNER, etc.
	Table  string   // e.g. "associations a"
	Alias  string   // e.g. "a"
	On     string   // join condition
	Select []string // columns to fetch from join
}

// Transaction interface for databases that support transactions
type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Index definition for NoSQL or relational indexing
type Index struct {
	Keys   map[string]int // e.g. {"email": 1, "created_at": -1}
	Unique bool
	Sparse bool
	Name   string
}

// DataStore — universal database interface
type DataStore interface {
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
	GetDatabase() interface{}
	Insert(ctx context.Context, table string, data interface{}) error
	InsertMany(ctx context.Context, table string, data []interface{}) error
	Update(ctx context.Context, filter map[string]interface{}, data interface{}) error
	UpdateMany(ctx context.Context, filter map[string]interface{}, data interface{}) error
	Upsert(ctx context.Context, table string, filter map[string]interface{}, data interface{}) error
	Delete(ctx context.Context, table string, filter map[string]interface{}) error
	DeleteMany(ctx context.Context, table string, filter map[string]interface{}) error

	FindOne(ctx context.Context, table, alias string, dest interface{}, opts *QueryOption) error
	FindMany(ctx context.Context, table, alias string, dest interface{}, opts *QueryOption) error
	Count(ctx context.Context, table, alias string, dest interface{}, opts *QueryOption) (int, error)
	Distinct(ctx context.Context, table, field string, filter map[string]interface{}, dest interface{}) error

	Aggregate(ctx context.Context, table string, pipeline interface{}, dest interface{}) error
	RawQuery(ctx context.Context, query string, args []interface{}, dest interface{}) error

	BeginTx(ctx context.Context) (Transaction, error)
	RunInTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error

	EnsureIndices(ctx context.Context, table string, indices []Index) error
	DropIndices(ctx context.Context, table string, indices []Index) error
}
