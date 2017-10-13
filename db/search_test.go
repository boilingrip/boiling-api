package db

import (
	"testing"

	"time"

	"github.com/stretchr/testify/require"
)

type columnA struct{}

func (columnA) column() string {
	return "a"
}

type columnB struct{}

func (columnB) column() string {
	return "b"
}

type columnC struct{}

func (columnC) column() string {
	return "c"
}

func TestQuerySimple(t *testing.T) {
	q := NewQuery(Eq(columnA{}, "test"))
	require.NotNil(t, q)

	query, params := q.Build()

	require.Equal(t, 1, len(params))
	require.Equal(t, []interface{}{"test"}, params)
	require.Equal(t, "\"a\" = $1", query)
}

func TestQueryComplex(t *testing.T) {
	var (
		p1 = "test"
		p2 = 3
		p3 = time.Now()
	)

	q := NewQuery(
		And(
			Eq(columnA{}, p1),
			Or(
				Eq(columnB{}, p2),
				Neq(columnC{}, p3),
			),
		),
	)
	require.NotNil(t, q)

	query, params := q.Build()
	require.Equal(t, 3, len(params))
	require.Equal(t, []interface{}{p1, p2, p3}, params)
	require.Equal(t, "( \"a\" = $1 AND ( \"b\" = $2 OR \"c\" != $3 ) )", query)
}

func TestQueryWithSorting(t *testing.T) {
	var (
		p1 = 4
		p2 = 3
	)

	q := NewQuery(Or(
		Eq(columnA{}, p1),
		Eq(columnA{}, p2),
	))
	require.NotNil(t, q)

	q.SetSorter(SortAscending(columnB{}))

	query, params := q.Build()

	require.Equal(t, 2, len(params))
	require.Equal(t, []interface{}{p1, p2}, params)
	require.Equal(t, "( \"a\" = $1 OR \"a\" = $2 ) ORDER BY \"b\" ASC", query)

	q.SetSorter(SortDescending(columnC{}))

	query, params = q.Build()

	require.Equal(t, 2, len(params))
	require.Equal(t, []interface{}{p1, p2}, params)
	require.Equal(t, "( \"a\" = $1 OR \"a\" = $2 ) ORDER BY \"c\" DESC", query)
}
