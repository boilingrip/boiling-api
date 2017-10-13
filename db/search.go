package db

import (
	"fmt"
)

type tok int

const (
	leftParen tok = iota
	rightParen
	eq
	neq
	lt
	lte
	gt
	gte
	comma
	dot

	and
	or
	where
	orderBy
	ascending
	descending

	table
	column
	placeholder
)

func (t tok) isSpecial() bool {
	switch t {
	case table, column, placeholder:
		return true
	}
	return false
}

func (t tok) string() string {
	switch t {
	case leftParen:
		return "("
	case rightParen:
		return ")"
	case eq:
		return "="
	case neq:
		return "!="
	case lt:
		return "<"
	case lte:
		return "<="
	case gt:
		return ">"
	case gte:
		return ">="
	case comma:
		return ","
	case dot:
		return "."
	case and:
		return "AND"
	case or:
		return "OR"
	case where:
		return "WHERE"
	case orderBy:
		return "ORDER BY"
	case ascending:
		return "ASC"
	case descending:
		return "DESC"
	case placeholder, column, table:
		return ""
	default:
		panic("unknown token")
	}
}

type token interface {
	Tok() tok
}

type simpleToken struct {
	tok tok
}

func (s simpleToken) Tok() tok { return s.tok }

type placeholderToken struct {
	value interface{}
}

func (t placeholderToken) Tok() tok { return placeholder }

func (t placeholderToken) Value() interface{} {
	return t.value
}

type tableToken struct {
	name string
}

func (t tableToken) Tok() tok { return table }

func (t tableToken) Table() string {
	return t.name
}

type columnToken struct {
	name string
}

func (t columnToken) Tok() tok { return column }

func (t columnToken) Column() string {
	return t.name
}

type queryBuilder struct {
	placeholderIndex int
	tokens           []token
}

func (b *queryBuilder) Build() (string, []interface{}) {
	b.placeholderIndex = 0
	var s string
	var placeholders []interface{}

	for i := 0; i < len(b.tokens); i++ {
		t := b.tokens[i]

		if !t.Tok().isSpecial() {
			s += t.Tok().string() + " "
		}

		switch t.Tok() {
		case column:
			c, ok := t.(columnToken)
			if !ok {
				panic("invalid column token")
			}

			s += fmt.Sprintf("%q ", c.Column())
		case placeholder:
			p, ok := t.(placeholderToken)
			if !ok {
				panic("invalid placeholder token")
			}

			b.placeholderIndex++
			placeholders = append(placeholders, p.Value())
			s += fmt.Sprintf("$%d ", b.placeholderIndex)
		case table:
			tab, ok := t.(tableToken)
			if !ok {
				panic("invalid table  token")
			}

			if len(b.tokens) > i+2 &&
				b.tokens[i+1].Tok() == dot &&
				b.tokens[i+2].Tok() == column {
				i++ // skip the dot
				s += fmt.Sprintf("%q.", tab.Table())
				break
			}
			panic("not a table-dot-column combo")
		}
	}

	if len(s) > 0 {
		s = s[:len(s)-1]
	}

	return s, placeholders
}

type Query struct {
	b      Boolean
	sorter Sorter
}

func NewQuery(b Boolean) *Query {
	return &Query{b: b}
}

func (q *Query) SetSorter(s Sorter) {
	q.sorter = s
}

func (q *Query) Build() (string, []interface{}) {
	b := queryBuilder{
		tokens: q.b.tokens(),
	}
	if q.sorter != nil {
		b.tokens = append(b.tokens, q.sorter.tokens()...)
	}

	return b.Build()
}

type tokenizer interface {
	tokens() []token
}

type ColumnSelector interface {
	column() string
}

type Sorter interface {
	tokenizer
	s()
}

type Boolean interface {
	tokenizer
	b()
}

type simpleSorter struct {
	column ColumnSelector
	desc   bool
}

func (s simpleSorter) tokens() []token {
	if s.desc {
		return []token{
			simpleToken{
				tok: orderBy,
			},
			columnToken{
				name: s.column.column(),
			},
			simpleToken{
				tok: descending,
			},
		}
	}

	return []token{
		simpleToken{
			tok: orderBy,
		},
		columnToken{
			name: s.column.column(),
		},
		simpleToken{
			tok: ascending,
		},
	}
}
func (s simpleSorter) s() {}

func SortAscending(column ColumnSelector) Sorter {
	return simpleSorter{
		column: column,
	}
}

func SortDescending(column ColumnSelector) Sorter {
	return simpleSorter{
		column: column,
		desc:   true,
	}
}

type booleanComparator struct {
	column ColumnSelector
	val    interface{}
	op     tok
}

func (b booleanComparator) tokens() []token {
	return []token{
		columnToken{
			name: b.column.column(),
		},
		simpleToken{
			tok: b.op,
		},
		placeholderToken{
			value: b.val,
		},
	}
}
func (b booleanComparator) b() {}

func Eq(column ColumnSelector, v interface{}) Boolean {
	return booleanComparator{
		column: column,
		val:    v,
		op:     eq,
	}
}

func Neq(column ColumnSelector, v interface{}) Boolean {
	return booleanComparator{
		column: column,
		val:    v,
		op:     neq,
	}
}

type booleanBinary struct {
	b1, b2      Boolean
	conjunction tok
}

func (b booleanBinary) tokens() []token {
	t := []token{simpleToken{tok: leftParen}}
	t = append(t, b.b1.tokens()...)
	t = append(t, simpleToken{tok: b.conjunction})
	t = append(t, b.b2.tokens()...)
	t = append(t, simpleToken{tok: rightParen})

	return t
}
func (b booleanBinary) b() {}

func And(b1, b2 Boolean) Boolean {
	return booleanBinary{
		b1:          b1,
		b2:          b2,
		conjunction: and,
	}
}

func Or(b1, b2 Boolean) Boolean {
	return booleanBinary{
		b1:          b1,
		b2:          b2,
		conjunction: or,
	}
}
