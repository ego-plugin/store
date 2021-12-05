package edb

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"github.com/ego-plugin/store/edb/dialect"
)

// InsertStmt builds `INSERT INTO ...`.
type InsertStmt struct {
	runner
	EventReceiver
	Dialect

	raw

	Table        string
	Column       []string
	Value        [][]interface{}
	Ignored      bool
	ReturnColumn []string
	RecordID     *int64
	comments     Comments
}

type InsertBuilder = InsertStmt

func (b *InsertStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Column) == 0 {
		return ErrColumnNotSpecified
	}

	err := b.comments.Build(d, buf)
	if err != nil {
		return err
	}

	if b.Ignored {
		buf.WriteString("INSERT IGNORE INTO ")
	} else {
		buf.WriteString("INSERT INTO ")
	}

	buf.WriteString(d.QuoteIdent(b.Table))

	var placeholderBuf strings.Builder
	placeholderBuf.WriteString("(")
	buf.WriteString(" (")
	for i, col := range b.Column {
		if i > 0 {
			buf.WriteString(",")
			placeholderBuf.WriteString(",")
		}
		buf.WriteString(d.QuoteIdent(col))
		placeholderBuf.WriteString(placeholder)
	}
	buf.WriteString(")")

	if d == dialect.MSSQL && len(b.ReturnColumn) > 0 {
		buf.WriteString(" OUTPUT ")
		for i, col := range b.ReturnColumn {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString("INSERTED." + d.QuoteIdent(col))
		}
	}

	buf.WriteString(" VALUES ")
	placeholderBuf.WriteString(")")
	placeholderStr := placeholderBuf.String()

	for i, tuple := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(placeholderStr)

		buf.WriteValue(tuple...)
	}

	if d != dialect.MSSQL && len(b.ReturnColumn) > 0 {
		buf.WriteString(" RETURNING ")
		for i, col := range b.ReturnColumn {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(d.QuoteIdent(col))
		}
	}

	return nil
}

// InsertInto creates an InsertStmt.
func InsertInto(table string) *InsertStmt {
	return &InsertStmt{
		Table: table,
	}
}

// InsertInto creates an InsertStmt.
func (sess *Session) InsertInto(table string) *InsertStmt {
	b := InsertInto(table)
	b.runner = sess
	b.EventReceiver = sess.EventReceiver
	b.Dialect = sess.Dialect
	return b
}

// InsertInto creates an InsertStmt.
func (tx *Tx) InsertInto(table string) *InsertStmt {
	b := InsertInto(table)
	b.runner = tx
	b.EventReceiver = tx.EventReceiver
	b.Dialect = tx.Dialect
	return b
}

// InsertBySql creates an InsertStmt from raw query.
func InsertBySql(query string, value ...interface{}) *InsertStmt {
	return &InsertStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// InsertBySql creates an InsertStmt from raw query.
func (sess *Session) InsertBySql(query string, value ...interface{}) *InsertStmt {
	b := InsertBySql(query, value...)
	b.runner = sess
	b.EventReceiver = sess.EventReceiver
	b.Dialect = sess.Dialect
	return b
}

// InsertBySql creates an InsertStmt from raw query.
func (tx *Tx) InsertBySql(query string, value ...interface{}) *InsertStmt {
	b := InsertBySql(query, value...)
	b.runner = tx
	b.EventReceiver = tx.EventReceiver
	b.Dialect = tx.Dialect
	return b
}

func (b *InsertStmt) Columns(column ...string) *InsertStmt {
	b.Column = column
	return b
}

// Comment adds a comment to prepended. All multi-line sql comment characters are stripped
func (b *InsertStmt) Comment(comment string) *InsertStmt {
	b.comments = b.comments.Append(comment)
	return b
}

// Ignore any insertion errors
func (b *InsertStmt) Ignore() *InsertStmt {
	b.Ignored = true
	return b
}

// Values adds a tuple to be inserted.
// The order of the tuple should match Columns.
func (b *InsertStmt) Values(value ...interface{}) *InsertStmt {
	b.Value = append(b.Value, value)
	return b
}

// StructValues  Struct 里 tag 的 db 不等于空的写入 values
func (b *InsertStmt) ScanStruct(value interface{}, column ...string) *InsertStmt {
	valueType := reflect.TypeOf(value)
	valueValue := reflect.ValueOf(value)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
		valueValue = valueValue.Elem()
	}
	if valueType.Kind() != reflect.Struct {
		return b
	}
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	columnLen := len(column)
	for index := 0; index < valueType.NumField(); index++ {
		typeField := valueType.Field(index)
		valueField := valueValue.Field(index)
		columnName := typeField.Tag.Get("db")
		columnSlice := strings.SplitN(columnName, ",", 2)
		columnName = columnSlice[0]
		if columnName == "" || columnName == "id" || columnName == "-" {
			continue
		}
		// 检查列是否存在
		if columnLen > 0 && !IsSliceContainsString(columnName, column...) {
			continue
		}
		columns = append(columns, columnName)
		values = append(values, valueField.Interface())
	}
	return b.Columns(columns...).Values(values...)
}

// Record adds a tuple for columns from a struct.
//
// If there is a field called "Id" or "ID" in the struct,
// it will be set to LastInsertId.
func (b *InsertStmt) Record(structValue interface{}) *InsertStmt {
	v := reflect.Indirect(reflect.ValueOf(structValue))

	if v.Kind() == reflect.Struct {
		found := make([]interface{}, len(b.Column)+1)
		// ID is recommended by golint here
		s := newTagStore()
		s.findValueByName(v, append(b.Column, "id"), found, false)

		value := found[:len(found)-1]
		for i, v := range value {
			if v != nil {
				value[i] = v.(reflect.Value).Interface()
			}
		}

		if v.CanSet() {
			switch idField := found[len(found)-1].(type) {
			case reflect.Value:
				if idField.Kind() == reflect.Int64 {
					b.RecordID = idField.Addr().Interface().(*int64)
				}
			}
		}
		b.Values(value...)
	}
	return b
}

// Returning specifies the returning columns for postgres/mssql.
func (b *InsertStmt) Returning(column ...string) *InsertStmt {
	b.ReturnColumn = column
	return b
}

// Pair adds (column, value) to be inserted.
// It is an error to mix Pair with Values and Record.
func (b *InsertStmt) Pair(column string, value interface{}) *InsertStmt {
	b.Column = append(b.Column, column)
	switch len(b.Value) {
	case 0:
		b.Values(value)
	case 1:
		b.Value[0] = append(b.Value[0], value)
	default:
		panic("pair only allows one record to insert")
	}
	return b
}

func (b *InsertStmt) Exec() (sql.Result, error) {
	return b.ExecContext(context.Background())
}

func (b *InsertStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	result, err := exec(ctx, b.runner, b.EventReceiver, b, b.Dialect)
	if err != nil {
		return nil, err
	}

	if b.RecordID != nil {
		if id, err := result.LastInsertId(); err == nil {
			*b.RecordID = id
		}
		b.RecordID = nil
	}

	return result, nil
}

func (b *InsertStmt) LoadContext(ctx context.Context, value interface{}) error {
	_, err := query(ctx, b.runner, b.EventReceiver, b, b.Dialect, value)
	return err
}

func (b *InsertStmt) Load(value interface{}) error {
	return b.LoadContext(context.Background(), value)
}
