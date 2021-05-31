package condition

import (
	"github.com/ego-plugin/store/edb"
	"reflect"
)

type Func func() (edb.Builder, bool)
type WhereFunc func() (edb.Builder, bool)
type GroupFunc func() (edb.Builder, bool)
type HavingFunc func() (edb.Builder, bool)

var (
	whereType   = reflect.TypeOf(WhereFunc(func() (edb.Builder, bool) { return nil, false }))
	groupByType = reflect.TypeOf(GroupFunc(func() (edb.Builder, bool) { return nil, false }))
	havingType  = reflect.TypeOf(HavingFunc(func() (edb.Builder, bool) { return nil, false }))
)

func (b Func) Build() (edb.Builder, bool) {
	return b()
}

func (b WhereFunc) Build() (edb.Builder, bool) {
	return b()
}

func (b GroupFunc) Build() (edb.Builder, bool) {
	return b()
}

func (b HavingFunc) Build() (edb.Builder, bool) {
	return b()
}

func In(cond ...Builder) edb.Builder {
	build, _ := And(cond...).Build()
	return build
}

// And creates AND from a list of conditions.
func And(cond ...Builder) Builder {
	return WhereFunc(func() (b edb.Builder, p bool) {
		builders := make([]edb.Builder, 0)
		for _, v := range cond {
			if build, ok := v.Build(); ok {
				builders = append(builders, build)
			}
		}
		if len(builders) == 1 {
			b = builders[0]
			return b, true
		}
		if len(builders) > 1 {
			b = edb.And(builders...)
			return b, true
		}
		return b, false
	})
}

// Or creates OR from a list of conditions.
func Or(cond ...Builder) Builder {
	return WhereFunc(func() (b edb.Builder, p bool) {
		builders := make([]edb.Builder, 0)
		for _, v := range cond {
			if build, ok := v.Build(); ok {
				builders = append(builders, build)
			}
		}
		if len(builders) == 1 {
			b = builders[0]
			return b, true
		}
		if len(builders) > 1 {
			b = edb.Or(builders...)
			return b, true
		}
		return b, false
	})
}

// Eq is `=`.
// When value is nil, it will be translated to `IS NULL`.
// When value is a slice, it will be translated to `IN`.
// Otherwise it will be translated to `=`.
func Eq(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Eq(column, value), ok
	})
}

// Neq is `!=`.
// When value is nil, it will be translated to `IS NOT NULL`.
// When value is a slice, it will be translated to `NOT IN`.
// Otherwise it will be translated to `!=`.
func Neq(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Neq(column, value), ok
	})
}

// Gt is `>`.
func Gt(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Gt(column, value), ok
	})
}

// Gte is '>='.
func Gte(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Gte(column, value), ok
	})
}

// Lt is '<'.
func Lt(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Lt(column, value), ok
	})
}

// Lte is `<=`.
func Lte(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Lte(column, value), ok
	})
}

// Like is `LIKE`, with an optional `ESCAPE` clause
func Like(column, value string, ok bool, escape ...string) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.Like(column, value, escape...), ok
	})
}

// NotLike is `NOT LIKE`, with an optional `ESCAPE` clause
func NotLike(column, value string, ok bool, escape ...string) Builder {
	return WhereFunc(func() (edb.Builder, bool) {
		return edb.NotLike(column, value, escape...), ok
	})
}