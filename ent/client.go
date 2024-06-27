// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/karust/openserp/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/karust/openserp/ent/searchquery"
	"github.com/karust/openserp/ent/serp"

	stdsql "database/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// SERP is the client for interacting with the SERP builders.
	SERP *SERPClient
	// SearchQuery is the client for interacting with the SearchQuery builders.
	SearchQuery *SearchQueryClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.SERP = NewSERPClient(c.config)
	c.SearchQuery = NewSearchQueryClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		SERP:        NewSERPClient(cfg),
		SearchQuery: NewSearchQueryClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		SERP:        NewSERPClient(cfg),
		SearchQuery: NewSearchQueryClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		SERP.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.SERP.Use(hooks...)
	c.SearchQuery.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.SERP.Intercept(interceptors...)
	c.SearchQuery.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *SERPMutation:
		return c.SERP.mutate(ctx, m)
	case *SearchQueryMutation:
		return c.SearchQuery.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// SERPClient is a client for the SERP schema.
type SERPClient struct {
	config
}

// NewSERPClient returns a client for the SERP from the given config.
func NewSERPClient(c config) *SERPClient {
	return &SERPClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `serp.Hooks(f(g(h())))`.
func (c *SERPClient) Use(hooks ...Hook) {
	c.hooks.SERP = append(c.hooks.SERP, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `serp.Intercept(f(g(h())))`.
func (c *SERPClient) Intercept(interceptors ...Interceptor) {
	c.inters.SERP = append(c.inters.SERP, interceptors...)
}

// Create returns a builder for creating a SERP entity.
func (c *SERPClient) Create() *SERPCreate {
	mutation := newSERPMutation(c.config, OpCreate)
	return &SERPCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of SERP entities.
func (c *SERPClient) CreateBulk(builders ...*SERPCreate) *SERPCreateBulk {
	return &SERPCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *SERPClient) MapCreateBulk(slice any, setFunc func(*SERPCreate, int)) *SERPCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &SERPCreateBulk{err: fmt.Errorf("calling to SERPClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*SERPCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &SERPCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for SERP.
func (c *SERPClient) Update() *SERPUpdate {
	mutation := newSERPMutation(c.config, OpUpdate)
	return &SERPUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SERPClient) UpdateOne(s *SERP) *SERPUpdateOne {
	mutation := newSERPMutation(c.config, OpUpdateOne, withSERP(s))
	return &SERPUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SERPClient) UpdateOneID(id int) *SERPUpdateOne {
	mutation := newSERPMutation(c.config, OpUpdateOne, withSERPID(id))
	return &SERPUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SERP.
func (c *SERPClient) Delete() *SERPDelete {
	mutation := newSERPMutation(c.config, OpDelete)
	return &SERPDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SERPClient) DeleteOne(s *SERP) *SERPDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SERPClient) DeleteOneID(id int) *SERPDeleteOne {
	builder := c.Delete().Where(serp.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SERPDeleteOne{builder}
}

// Query returns a query builder for SERP.
func (c *SERPClient) Query() *SERPQuery {
	return &SERPQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSERP},
		inters: c.Interceptors(),
	}
}

// Get returns a SERP entity by its id.
func (c *SERPClient) Get(ctx context.Context, id int) (*SERP, error) {
	return c.Query().Where(serp.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SERPClient) GetX(ctx context.Context, id int) *SERP {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySearchQuery queries the search_query edge of a SERP.
func (c *SERPClient) QuerySearchQuery(s *SERP) *SearchQueryQuery {
	query := (&SearchQueryClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(serp.Table, serp.FieldID, id),
			sqlgraph.To(searchquery.Table, searchquery.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serp.SearchQueryTable, serp.SearchQueryColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SERPClient) Hooks() []Hook {
	return c.hooks.SERP
}

// Interceptors returns the client interceptors.
func (c *SERPClient) Interceptors() []Interceptor {
	return c.inters.SERP
}

func (c *SERPClient) mutate(ctx context.Context, m *SERPMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SERPCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SERPUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SERPUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SERPDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown SERP mutation op: %q", m.Op())
	}
}

// SearchQueryClient is a client for the SearchQuery schema.
type SearchQueryClient struct {
	config
}

// NewSearchQueryClient returns a client for the SearchQuery from the given config.
func NewSearchQueryClient(c config) *SearchQueryClient {
	return &SearchQueryClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `searchquery.Hooks(f(g(h())))`.
func (c *SearchQueryClient) Use(hooks ...Hook) {
	c.hooks.SearchQuery = append(c.hooks.SearchQuery, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `searchquery.Intercept(f(g(h())))`.
func (c *SearchQueryClient) Intercept(interceptors ...Interceptor) {
	c.inters.SearchQuery = append(c.inters.SearchQuery, interceptors...)
}

// Create returns a builder for creating a SearchQuery entity.
func (c *SearchQueryClient) Create() *SearchQueryCreate {
	mutation := newSearchQueryMutation(c.config, OpCreate)
	return &SearchQueryCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of SearchQuery entities.
func (c *SearchQueryClient) CreateBulk(builders ...*SearchQueryCreate) *SearchQueryCreateBulk {
	return &SearchQueryCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *SearchQueryClient) MapCreateBulk(slice any, setFunc func(*SearchQueryCreate, int)) *SearchQueryCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &SearchQueryCreateBulk{err: fmt.Errorf("calling to SearchQueryClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*SearchQueryCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &SearchQueryCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for SearchQuery.
func (c *SearchQueryClient) Update() *SearchQueryUpdate {
	mutation := newSearchQueryMutation(c.config, OpUpdate)
	return &SearchQueryUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SearchQueryClient) UpdateOne(sq *SearchQuery) *SearchQueryUpdateOne {
	mutation := newSearchQueryMutation(c.config, OpUpdateOne, withSearchQuery(sq))
	return &SearchQueryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SearchQueryClient) UpdateOneID(id int) *SearchQueryUpdateOne {
	mutation := newSearchQueryMutation(c.config, OpUpdateOne, withSearchQueryID(id))
	return &SearchQueryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SearchQuery.
func (c *SearchQueryClient) Delete() *SearchQueryDelete {
	mutation := newSearchQueryMutation(c.config, OpDelete)
	return &SearchQueryDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SearchQueryClient) DeleteOne(sq *SearchQuery) *SearchQueryDeleteOne {
	return c.DeleteOneID(sq.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SearchQueryClient) DeleteOneID(id int) *SearchQueryDeleteOne {
	builder := c.Delete().Where(searchquery.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SearchQueryDeleteOne{builder}
}

// Query returns a query builder for SearchQuery.
func (c *SearchQueryClient) Query() *SearchQueryQuery {
	return &SearchQueryQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSearchQuery},
		inters: c.Interceptors(),
	}
}

// Get returns a SearchQuery entity by its id.
func (c *SearchQueryClient) Get(ctx context.Context, id int) (*SearchQuery, error) {
	return c.Query().Where(searchquery.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SearchQueryClient) GetX(ctx context.Context, id int) *SearchQuery {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySerps queries the serps edge of a SearchQuery.
func (c *SearchQueryClient) QuerySerps(sq *SearchQuery) *SERPQuery {
	query := (&SERPClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := sq.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(searchquery.Table, searchquery.FieldID, id),
			sqlgraph.To(serp.Table, serp.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, searchquery.SerpsTable, searchquery.SerpsColumn),
		)
		fromV = sqlgraph.Neighbors(sq.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SearchQueryClient) Hooks() []Hook {
	return c.hooks.SearchQuery
}

// Interceptors returns the client interceptors.
func (c *SearchQueryClient) Interceptors() []Interceptor {
	return c.inters.SearchQuery
}

func (c *SearchQueryClient) mutate(ctx context.Context, m *SearchQueryMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SearchQueryCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SearchQueryUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SearchQueryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SearchQueryDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown SearchQuery mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		SERP, SearchQuery []ent.Hook
	}
	inters struct {
		SERP, SearchQuery []ent.Interceptor
	}
)

// ExecContext allows calling the underlying ExecContext method of the driver if it is supported by it.
// See, database/sql#DB.ExecContext for more information.
func (c *config) ExecContext(ctx context.Context, query string, args ...any) (stdsql.Result, error) {
	ex, ok := c.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	return ex.ExecContext(ctx, query, args...)
}

// QueryContext allows calling the underlying QueryContext method of the driver if it is supported by it.
// See, database/sql#DB.QueryContext for more information.
func (c *config) QueryContext(ctx context.Context, query string, args ...any) (*stdsql.Rows, error) {
	q, ok := c.driver.(interface {
		QueryContext(context.Context, string, ...any) (*stdsql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	return q.QueryContext(ctx, query, args...)
}
