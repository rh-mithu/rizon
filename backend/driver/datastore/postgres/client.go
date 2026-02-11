package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rh-mithu/rizon/backend/config"
	"github.com/rh-mithu/rizon/backend/driver/datastore"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Store implements the datastore.DataStore interface
type Store struct {
	db *bun.DB
}

// NewStore initializes a new PostgresSQL Bun connection
func NewStore(cfg *config.Config, l *slog.Logger) *Store {
	dsn := cfg.SQLDatabaseURL
	if dsn == "" {
		log.Fatal("DATABASE_URL not found in environment variables")
	}
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	db := bun.NewDB(sqlDB, pgdialect.New())
	if cfg.Env != "production" {
		db = db.WithQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	if err = db.Ping(); err != nil {
		l.Error("failed to ping Postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	db.WithQueryHook(&QueryHook{})
	err = db.Ping()
	if err != nil {
		l.Error("failed to ping Postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	l.Info("Connected to Postgres successfully")
	return &Store{db: db}
}

// DB exposes the underlying *bun.DB (for advanced usage)
func (s *Store) DB() *bun.DB {
	return s.db
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Store) GetDatabase() interface{} {
	return s.db
}

func (s *Store) Disconnect(ctx context.Context) error {
	return s.db.Close()
}

func (s *Store) Insert(ctx context.Context, table string, data interface{}) error {
	_, err := s.db.NewInsert().
		Model(data).
		ModelTableExpr(table).
		Exec(ctx)
	return err
}

func (s *Store) InsertMany(ctx context.Context, table string, data []interface{}) error {
	_, err := s.db.NewInsert().Model(&data).ModelTableExpr(table).Exec(ctx)
	return err
}

func (s *Store) Update(ctx context.Context, filter map[string]interface{}, data interface{}) error {
	q := s.db.NewUpdate().Model(data).OmitZero()

	for k, v := range filter {
		q = q.Where(fmt.Sprintf("%s = ?", k), v)
	}

	_, err := q.Exec(ctx)
	return err
}

func (s *Store) UpdateMany(ctx context.Context, filter map[string]interface{}, data interface{}) error {
	return s.Update(ctx, filter, data)
}

func (s *Store) Upsert(ctx context.Context, table string, filter map[string]interface{}, data interface{}) error {
	q := s.db.NewInsert().Model(data).ModelTableExpr(table)
	for k := range filter {
		q = q.On("CONFLICT (" + k + ") DO UPDATE").Set(fmt.Sprintf("%s = EXCLUDED.%s", k, k))
		break
	}
	_, err := q.Exec(ctx)
	return err
}

func (s *Store) Delete(ctx context.Context, table string, filter map[string]interface{}) error {
	q := s.db.NewDelete().ModelTableExpr(table)
	for k, v := range filter {
		q = q.Where(fmt.Sprintf("%s = ?", k), v)
	}
	_, err := q.Exec(ctx)
	return err
}

func (s *Store) DeleteMany(ctx context.Context, table string, filter map[string]interface{}) error {
	return s.Delete(ctx, table, filter)
}

// FindOne fetches a single record from the given table with optional alias and query options.
func (s *Store) FindOne(ctx context.Context, table, alias string, dest interface{}, opts *datastore.QueryOption) error {
	if dest == nil {
		return fmt.Errorf("dest cannot be nil")
	}

	q := s.db.NewSelect().
		Model(dest)

	applyQueryOptions(q, opts, alias)

	return q.Scan(ctx)
}

// FindMany fetches multiple records from the table with optional alias and query options.
func (s *Store) FindMany(
	ctx context.Context,
	table, alias string,
	dest interface{},
	opts *datastore.QueryOption,
) error {
	if dest == nil {
		return fmt.Errorf("dest cannot be nil")
	}

	// FIX: Reverting to the canonical TableExpr to resolve the 'WrapWith undefined' compilation error.
	q := s.db.NewSelect().
		Model(dest)

	applyQueryOptions(q, opts, alias)

	return q.Scan(ctx)
}

// applyQueryOptions applies QueryOption struct to a Bun SelectQuery
func applyQueryOptions(q *bun.SelectQuery, opts *datastore.QueryOption, alias string) {
	if opts == nil {
		return
	}

	// --- SELECT main table columns ---
	q = q.Column(fmt.Sprintf(`%s.*`, alias))

	// --- Select extra columns if needed ---
	for _, c := range opts.Select {
		q = q.Column(fmt.Sprintf(`"%s"."%s"`, alias, c))
	}

	// --- Relations ---
	for _, rel := range opts.Relations {
		q = q.Relation(rel) // Bun will handle JOINs and nested structs automatically
	}

	// --- Filters ---
	for k, v := range opts.Filter {
		switch {
		// ✅ Support IS NULL
		case strings.HasSuffix(k, "__is_null"):
			col := strings.TrimSuffix(k, "__is_null")
			q = q.Where(fmt.Sprintf(`"%s"."%s" IS NULL`, alias, col))

		// ✅ Support IS NOT NULL
		case strings.HasSuffix(k, "__is_not_null"):
			col := strings.TrimSuffix(k, "__is_not_null")
			q = q.Where(fmt.Sprintf(`"%s"."%s" IS NOT NULL`, alias, col))

		// ✅ Support IN (slice values)
		case strings.HasSuffix(k, "__in"):
			col := strings.TrimSuffix(k, "__in")
			q = q.Where(fmt.Sprintf(`"%s"."%s" IN (?)`, alias, col), bun.In(v))

		// ✅ Support not equal
		case strings.HasSuffix(k, "__ne"):
			col := strings.TrimSuffix(k, "__ne")
			q = q.Where(fmt.Sprintf(`"%s"."%s" != ?`, alias, col), v)

		// ✅ Default equal
		default:
			q = q.Where(fmt.Sprintf(`"%s"."%s" = ?`, alias, k), v)
		}
	}

	// --- Sorting ---
	for field, order := range opts.Sort {
		q = q.Order(fmt.Sprintf(`%s.%s %s`, alias, field, order))
	}

	// --- Distinct ---
	if opts.Distinct {
		q = q.Distinct()
	}

	// --- Pagination ---
	if opts.Limit > 0 {
		q = q.Limit(int(opts.Limit))
	}
	if opts.Skip > 0 {
		q = q.Offset(int(opts.Skip))
	}
}

func (s *Store) Count(ctx context.Context, table, alias string, model interface{}, opts *datastore.QueryOption) (int, error) {
	q := s.db.NewSelect().Model(model)
	applyQueryOptions(q, opts, alias)
	return q.Count(ctx)
}

func (s *Store) Distinct(ctx context.Context, table, field string, filter map[string]interface{}, dest interface{}) error {
	q := s.db.NewSelect().ModelTableExpr(table).ColumnExpr("DISTINCT ?", bun.Ident(field))
	for k, v := range filter {
		q = q.Where(fmt.Sprintf("%s = ?", k), v)
	}
	return q.Scan(ctx, dest)
}

func (s *Store) Aggregate(ctx context.Context, table string, pipeline interface{}, dest interface{}) error {
	return fmt.Errorf("aggregate not implemented; use RawQuery instead")
}

func (s *Store) RawQuery(ctx context.Context, query string, args []interface{}, dest interface{}) error {
	return s.db.NewRaw(query, args...).Scan(ctx, dest)
}

func (s *Store) EnsureIndices(ctx context.Context, table string, indices []datastore.Index) error {
	// Implementation depends on schema management needs
	return nil
}

func (s *Store) DropIndices(ctx context.Context, table string, indices []datastore.Index) error {
	return nil
}

// --- Transaction Management ---

type pgTx struct {
	tx bun.Tx
}

func (t *pgTx) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *pgTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

func (s *Store) BeginTx(ctx context.Context) (datastore.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &pgTx{tx: tx}, nil
}

func (s *Store) RunInTransaction(ctx context.Context, fn func(ctx context.Context, tx datastore.Transaction) error) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, bunTx bun.Tx) error {
		tx := &pgTx{tx: bunTx}
		return fn(ctx, tx)
	})
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if contains := len(sub) > 0 && len(s) >= len(sub) &&
			(len(s) >= len(sub)) &&
			(func() bool {
				return len(s) > 0 && (len(sub) == 0 ||
					(len(s) >= len(sub) && (s[len(s)-len(sub):] == sub || s[:len(sub)] == sub ||
						len(s) > len(sub) && (containsAny(s[1:], sub)))))
			})(); contains {
			return true
		}
	}
	return false
}

func buildUpdateMap(data interface{}) (map[string]interface{}, error) {
	updateMap := make(map[string]interface{})

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, fmt.Errorf("data must be a non-nil pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data must point to a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		if !value.CanInterface() {
			continue
		}

		// Skip nil pointers
		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		// Skip nil or empty slices
		if value.Kind() == reflect.Slice && (value.IsNil() || value.Len() == 0) {
			continue
		}

		// Skip BaseModel or audit fields
		if field.Type.Name() == "BaseModel" || field.Type.Name() == "UpdateAuditParam" {
			continue
		}

		tag := field.Tag.Get("bun")
		columnName := parseBunColumn(tag)
		if columnName == "" {
			columnName = field.Name
		}

		updateMap[columnName] = value.Interface()
	}

	return updateMap, nil
}

// parseBunColumn extracts the column name from the bun tag
func parseBunColumn(tag string) string {
	parts := strings.Split(tag, ",")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	fmt.Println(time.Since(event.StartTime), event.Query)
}
