// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/karust/openserp/ent/predicate"
	"github.com/karust/openserp/ent/searchquery"
	"github.com/karust/openserp/ent/serp"
)

// SearchQueryUpdate is the builder for updating SearchQuery entities.
type SearchQueryUpdate struct {
	config
	hooks    []Hook
	mutation *SearchQueryMutation
}

// Where appends a list predicates to the SearchQueryUpdate builder.
func (squ *SearchQueryUpdate) Where(ps ...predicate.SearchQuery) *SearchQueryUpdate {
	squ.mutation.Where(ps...)
	return squ
}

// SetQuery sets the "query" field.
func (squ *SearchQueryUpdate) SetQuery(s string) *SearchQueryUpdate {
	squ.mutation.SetQuery(s)
	return squ
}

// SetNillableQuery sets the "query" field if the given value is not nil.
func (squ *SearchQueryUpdate) SetNillableQuery(s *string) *SearchQueryUpdate {
	if s != nil {
		squ.SetQuery(*s)
	}
	return squ
}

// SetLocation sets the "location" field.
func (squ *SearchQueryUpdate) SetLocation(s string) *SearchQueryUpdate {
	squ.mutation.SetLocation(s)
	return squ
}

// SetNillableLocation sets the "location" field if the given value is not nil.
func (squ *SearchQueryUpdate) SetNillableLocation(s *string) *SearchQueryUpdate {
	if s != nil {
		squ.SetLocation(*s)
	}
	return squ
}

// SetLanguage sets the "language" field.
func (squ *SearchQueryUpdate) SetLanguage(s string) *SearchQueryUpdate {
	squ.mutation.SetLanguage(s)
	return squ
}

// SetNillableLanguage sets the "language" field if the given value is not nil.
func (squ *SearchQueryUpdate) SetNillableLanguage(s *string) *SearchQueryUpdate {
	if s != nil {
		squ.SetLanguage(*s)
	}
	return squ
}

// SetIsCanceled sets the "is_canceled" field.
func (squ *SearchQueryUpdate) SetIsCanceled(b bool) *SearchQueryUpdate {
	squ.mutation.SetIsCanceled(b)
	return squ
}

// SetNillableIsCanceled sets the "is_canceled" field if the given value is not nil.
func (squ *SearchQueryUpdate) SetNillableIsCanceled(b *bool) *SearchQueryUpdate {
	if b != nil {
		squ.SetIsCanceled(*b)
	}
	return squ
}

// SetCreatedAt sets the "created_at" field.
func (squ *SearchQueryUpdate) SetCreatedAt(t time.Time) *SearchQueryUpdate {
	squ.mutation.SetCreatedAt(t)
	return squ
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (squ *SearchQueryUpdate) SetNillableCreatedAt(t *time.Time) *SearchQueryUpdate {
	if t != nil {
		squ.SetCreatedAt(*t)
	}
	return squ
}

// AddSerpIDs adds the "serps" edge to the SERP entity by IDs.
func (squ *SearchQueryUpdate) AddSerpIDs(ids ...int) *SearchQueryUpdate {
	squ.mutation.AddSerpIDs(ids...)
	return squ
}

// AddSerps adds the "serps" edges to the SERP entity.
func (squ *SearchQueryUpdate) AddSerps(s ...*SERP) *SearchQueryUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.AddSerpIDs(ids...)
}

// Mutation returns the SearchQueryMutation object of the builder.
func (squ *SearchQueryUpdate) Mutation() *SearchQueryMutation {
	return squ.mutation
}

// ClearSerps clears all "serps" edges to the SERP entity.
func (squ *SearchQueryUpdate) ClearSerps() *SearchQueryUpdate {
	squ.mutation.ClearSerps()
	return squ
}

// RemoveSerpIDs removes the "serps" edge to SERP entities by IDs.
func (squ *SearchQueryUpdate) RemoveSerpIDs(ids ...int) *SearchQueryUpdate {
	squ.mutation.RemoveSerpIDs(ids...)
	return squ
}

// RemoveSerps removes "serps" edges to SERP entities.
func (squ *SearchQueryUpdate) RemoveSerps(s ...*SERP) *SearchQueryUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squ.RemoveSerpIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (squ *SearchQueryUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, squ.sqlSave, squ.mutation, squ.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (squ *SearchQueryUpdate) SaveX(ctx context.Context) int {
	affected, err := squ.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (squ *SearchQueryUpdate) Exec(ctx context.Context) error {
	_, err := squ.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (squ *SearchQueryUpdate) ExecX(ctx context.Context) {
	if err := squ.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (squ *SearchQueryUpdate) check() error {
	if v, ok := squ.mutation.Query(); ok {
		if err := searchquery.QueryValidator(v); err != nil {
			return &ValidationError{Name: "query", err: fmt.Errorf(`ent: validator failed for field "SearchQuery.query": %w`, err)}
		}
	}
	return nil
}

func (squ *SearchQueryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := squ.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(searchquery.Table, searchquery.Columns, sqlgraph.NewFieldSpec(searchquery.FieldID, field.TypeInt))
	if ps := squ.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := squ.mutation.Query(); ok {
		_spec.SetField(searchquery.FieldQuery, field.TypeString, value)
	}
	if value, ok := squ.mutation.Location(); ok {
		_spec.SetField(searchquery.FieldLocation, field.TypeString, value)
	}
	if value, ok := squ.mutation.Language(); ok {
		_spec.SetField(searchquery.FieldLanguage, field.TypeString, value)
	}
	if value, ok := squ.mutation.IsCanceled(); ok {
		_spec.SetField(searchquery.FieldIsCanceled, field.TypeBool, value)
	}
	if value, ok := squ.mutation.CreatedAt(); ok {
		_spec.SetField(searchquery.FieldCreatedAt, field.TypeTime, value)
	}
	if squ.mutation.SerpsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   searchquery.SerpsTable,
			Columns: []string{searchquery.SerpsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.mutation.RemovedSerpsIDs(); len(nodes) > 0 && !squ.mutation.SerpsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   searchquery.SerpsTable,
			Columns: []string{searchquery.SerpsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squ.mutation.SerpsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   searchquery.SerpsTable,
			Columns: []string{searchquery.SerpsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, squ.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{searchquery.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	squ.mutation.done = true
	return n, nil
}

// SearchQueryUpdateOne is the builder for updating a single SearchQuery entity.
type SearchQueryUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SearchQueryMutation
}

// SetQuery sets the "query" field.
func (squo *SearchQueryUpdateOne) SetQuery(s string) *SearchQueryUpdateOne {
	squo.mutation.SetQuery(s)
	return squo
}

// SetNillableQuery sets the "query" field if the given value is not nil.
func (squo *SearchQueryUpdateOne) SetNillableQuery(s *string) *SearchQueryUpdateOne {
	if s != nil {
		squo.SetQuery(*s)
	}
	return squo
}

// SetLocation sets the "location" field.
func (squo *SearchQueryUpdateOne) SetLocation(s string) *SearchQueryUpdateOne {
	squo.mutation.SetLocation(s)
	return squo
}

// SetNillableLocation sets the "location" field if the given value is not nil.
func (squo *SearchQueryUpdateOne) SetNillableLocation(s *string) *SearchQueryUpdateOne {
	if s != nil {
		squo.SetLocation(*s)
	}
	return squo
}

// SetLanguage sets the "language" field.
func (squo *SearchQueryUpdateOne) SetLanguage(s string) *SearchQueryUpdateOne {
	squo.mutation.SetLanguage(s)
	return squo
}

// SetNillableLanguage sets the "language" field if the given value is not nil.
func (squo *SearchQueryUpdateOne) SetNillableLanguage(s *string) *SearchQueryUpdateOne {
	if s != nil {
		squo.SetLanguage(*s)
	}
	return squo
}

// SetIsCanceled sets the "is_canceled" field.
func (squo *SearchQueryUpdateOne) SetIsCanceled(b bool) *SearchQueryUpdateOne {
	squo.mutation.SetIsCanceled(b)
	return squo
}

// SetNillableIsCanceled sets the "is_canceled" field if the given value is not nil.
func (squo *SearchQueryUpdateOne) SetNillableIsCanceled(b *bool) *SearchQueryUpdateOne {
	if b != nil {
		squo.SetIsCanceled(*b)
	}
	return squo
}

// SetCreatedAt sets the "created_at" field.
func (squo *SearchQueryUpdateOne) SetCreatedAt(t time.Time) *SearchQueryUpdateOne {
	squo.mutation.SetCreatedAt(t)
	return squo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (squo *SearchQueryUpdateOne) SetNillableCreatedAt(t *time.Time) *SearchQueryUpdateOne {
	if t != nil {
		squo.SetCreatedAt(*t)
	}
	return squo
}

// AddSerpIDs adds the "serps" edge to the SERP entity by IDs.
func (squo *SearchQueryUpdateOne) AddSerpIDs(ids ...int) *SearchQueryUpdateOne {
	squo.mutation.AddSerpIDs(ids...)
	return squo
}

// AddSerps adds the "serps" edges to the SERP entity.
func (squo *SearchQueryUpdateOne) AddSerps(s ...*SERP) *SearchQueryUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.AddSerpIDs(ids...)
}

// Mutation returns the SearchQueryMutation object of the builder.
func (squo *SearchQueryUpdateOne) Mutation() *SearchQueryMutation {
	return squo.mutation
}

// ClearSerps clears all "serps" edges to the SERP entity.
func (squo *SearchQueryUpdateOne) ClearSerps() *SearchQueryUpdateOne {
	squo.mutation.ClearSerps()
	return squo
}

// RemoveSerpIDs removes the "serps" edge to SERP entities by IDs.
func (squo *SearchQueryUpdateOne) RemoveSerpIDs(ids ...int) *SearchQueryUpdateOne {
	squo.mutation.RemoveSerpIDs(ids...)
	return squo
}

// RemoveSerps removes "serps" edges to SERP entities.
func (squo *SearchQueryUpdateOne) RemoveSerps(s ...*SERP) *SearchQueryUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return squo.RemoveSerpIDs(ids...)
}

// Where appends a list predicates to the SearchQueryUpdate builder.
func (squo *SearchQueryUpdateOne) Where(ps ...predicate.SearchQuery) *SearchQueryUpdateOne {
	squo.mutation.Where(ps...)
	return squo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (squo *SearchQueryUpdateOne) Select(field string, fields ...string) *SearchQueryUpdateOne {
	squo.fields = append([]string{field}, fields...)
	return squo
}

// Save executes the query and returns the updated SearchQuery entity.
func (squo *SearchQueryUpdateOne) Save(ctx context.Context) (*SearchQuery, error) {
	return withHooks(ctx, squo.sqlSave, squo.mutation, squo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (squo *SearchQueryUpdateOne) SaveX(ctx context.Context) *SearchQuery {
	node, err := squo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (squo *SearchQueryUpdateOne) Exec(ctx context.Context) error {
	_, err := squo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (squo *SearchQueryUpdateOne) ExecX(ctx context.Context) {
	if err := squo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (squo *SearchQueryUpdateOne) check() error {
	if v, ok := squo.mutation.Query(); ok {
		if err := searchquery.QueryValidator(v); err != nil {
			return &ValidationError{Name: "query", err: fmt.Errorf(`ent: validator failed for field "SearchQuery.query": %w`, err)}
		}
	}
	return nil
}

func (squo *SearchQueryUpdateOne) sqlSave(ctx context.Context) (_node *SearchQuery, err error) {
	if err := squo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(searchquery.Table, searchquery.Columns, sqlgraph.NewFieldSpec(searchquery.FieldID, field.TypeInt))
	id, ok := squo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SearchQuery.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := squo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, searchquery.FieldID)
		for _, f := range fields {
			if !searchquery.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != searchquery.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := squo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := squo.mutation.Query(); ok {
		_spec.SetField(searchquery.FieldQuery, field.TypeString, value)
	}
	if value, ok := squo.mutation.Location(); ok {
		_spec.SetField(searchquery.FieldLocation, field.TypeString, value)
	}
	if value, ok := squo.mutation.Language(); ok {
		_spec.SetField(searchquery.FieldLanguage, field.TypeString, value)
	}
	if value, ok := squo.mutation.IsCanceled(); ok {
		_spec.SetField(searchquery.FieldIsCanceled, field.TypeBool, value)
	}
	if value, ok := squo.mutation.CreatedAt(); ok {
		_spec.SetField(searchquery.FieldCreatedAt, field.TypeTime, value)
	}
	if squo.mutation.SerpsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   searchquery.SerpsTable,
			Columns: []string{searchquery.SerpsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.mutation.RemovedSerpsIDs(); len(nodes) > 0 && !squo.mutation.SerpsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   searchquery.SerpsTable,
			Columns: []string{searchquery.SerpsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := squo.mutation.SerpsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   searchquery.SerpsTable,
			Columns: []string{searchquery.SerpsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &SearchQuery{config: squo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, squo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{searchquery.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	squo.mutation.done = true
	return _node, nil
}
