package condition

import (
	"github.com/ego-plugin/store/edb"
	"reflect"
)

type Builder interface {
	Build() (edb.Builder, bool)
}

type Builders struct {
	Value []Builder
}

func NewBuilders() *Builders {
	return &Builders{
		Value: make([]Builder, 0),
	}
}

// ScanSelectBuilder 扫描Builder写入查询SQL语句
func (b *Builders) ScanSelectBuilder(stmt edb.SelectBuilder) {
	for _, v := range b.Value {
		switch reflect.TypeOf(v) {
		case whereType:
			if build, ok := v.Build(); ok {
				stmt.Where(build)
			}
		case groupByType:
			if build, ok := v.Build(); ok {
				stmt.Group = append(stmt.Group, build)
			}
		case havingType:
			if build, ok := v.Build(); ok {
				stmt.Having(build)
			}
		}
	}
}

func (b *Builders) Append(v Builder) {
	b.Value = append(b.Value, v)
}