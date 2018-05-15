package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/LyricTian/retry"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

// M 定义字典
type M map[string]interface{}

// 自定义数据库语句打印日志
type dbLogger struct {
	logger *log.Logger
}

func (l *dbLogger) Init() gorp.GorpLogger {
	l.logger = log.New(os.Stdout, "", log.Lmicroseconds)
	return l
}

func (l *dbLogger) Printf(format string, v ...interface{}) {
	query := fmt.Sprint(v[1])
	query = strings.Replace(query, "\n", " ", -1)
	query = strings.Replace(query, "\t", "", -1)
	v[1] = query

	l.logger.Printf(format, v...)
}

type options struct {
	dsn          string        // 连接串
	trace        bool          // 追踪调试
	maxLifetime  time.Duration // 设置连接可以被重新使用的最大时间量
	maxOpenConns int           // 设置打开连接到数据库的最大数量
	maxIdleConns int           // 设置空闲连接池中的最大连接数
}

// Option 配置项
type Option func(*options)

// SetDSN 设置连接串
func SetDSN(dsn string) Option {
	return func(o *options) {
		o.dsn = dsn
	}
}

// SetTrace 设置追踪调试
func SetTrace(t bool) Option {
	return func(o *options) {
		o.trace = t
	}
}

// SetMaxLifetime 设置连接可以被重新使用的最大时间量
func SetMaxLifetime(maxLifetime time.Duration) Option {
	return func(o *options) {
		o.maxLifetime = maxLifetime
	}
}

// SetMaxOpenConns 设置打开连接到数据库的最大数量
func SetMaxOpenConns(maxOpenConns int) Option {
	return func(o *options) {
		o.maxOpenConns = maxOpenConns
	}
}

// SetMaxIdleConns 设置空闲连接池中的最大连接数
func SetMaxIdleConns(maxIdleConns int) Option {
	return func(o *options) {
		o.maxIdleConns = maxIdleConns
	}
}

// DB 数据库
type DB struct {
	*gorp.DbMap
}

// NewMySQL 创建MySQL数据库实例
func NewMySQL(opts ...Option) (*sql.DB, bool, error) {
	o := &options{
		maxLifetime:  time.Hour * 2,
		maxOpenConns: 150,
		maxIdleConns: 50,
	}

	for _, opt := range opts {
		opt(o)
	}

	db, err := sql.Open("mysql", o.dsn)
	if err != nil {
		return nil, o.trace, err
	}

	// 尝试发送Ping包
	err = retry.DoFunc(3, func() error {
		perr := db.Ping()
		if perr != nil {
			fmt.Println("发送ping值错误：", perr.Error())
		}
		return perr
	}, func(i int) time.Duration {
		return time.Second
	})
	if err != nil {
		return nil, o.trace, err
	}

	db.SetMaxOpenConns(o.maxOpenConns)
	db.SetMaxIdleConns(o.maxIdleConns)
	db.SetConnMaxLifetime(o.maxLifetime)

	return db, o.trace, nil
}

// NewMySQLWithDB 创建DB
func NewMySQLWithDB(db *sql.DB, trace bool) *DB {
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Encoding: "UTF8", Engine: "InnoDB"}}
	if trace {
		dbMap.TraceOn("[db]", new(dbLogger).Init())
	}

	return &DB{dbMap}
}

// Close 关闭数据库连接
func (m *DB) Close() error {
	if m.DbMap == nil {
		return nil
	}
	return m.Db.Close()
}

// InsertSQL 获取插入SQL
func (m *DB) InsertSQL(table string, info M) (string, []interface{}) {
	q := fmt.Sprintf("INSERT INTO %s", table)

	var (
		cols []string
		vals []interface{}
	)

	for k, v := range info {
		cols = append(cols, k)
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s(%s) VALUES(%s)", q, strings.Join(cols, ","), strings.Repeat(",?", len(cols))[1:])
	return q, vals
}

// InsertM 插入数据
func (m *DB) InsertM(table string, info M) (int64, error) {
	q, vals := m.InsertSQL(table, info)
	result, err := m.Exec(q, vals...)
	if err != nil {
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()

	return lastInsertID, nil
}

// InsertMWithTran 使用事物插入数据
func (m *DB) InsertMWithTran(tran *gorp.Transaction, table string, info M) (int64, error) {
	q, vals := m.InsertSQL(table, info)
	result, err := tran.Exec(q, vals...)
	if err != nil {
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()

	return lastInsertID, nil
}

// UpdateSQL 获取更新SQL
func (m *DB) UpdateSQL(table string, pk, info M) (string, []interface{}) {
	q := fmt.Sprintf("UPDATE %s SET", table)

	var (
		cols []string
		vals []interface{}
	)

	for k, v := range info {
		cols = append(cols, fmt.Sprintf("%s=?", k))
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s %s", q, strings.Join(cols, ","))
	cols = nil

	for k, v := range pk {
		cols = append(cols, fmt.Sprintf("%s=?", k))
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s WHERE %s", q, strings.Join(cols, " and "))
	return q, vals
}

// UpdateByPK 更新表数据
func (m *DB) UpdateByPK(table string, pk, info M) (int64, error) {
	q, vals := m.UpdateSQL(table, pk, info)
	result, err := m.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// UpdateByPKWithTran 使用事物更新表数据
func (m *DB) UpdateByPKWithTran(tran *gorp.Transaction, table string, pk, info M) (int64, error) {
	q, vals := m.UpdateSQL(table, pk, info)
	result, err := tran.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// DeleteSQL 获取删除SQL
func (m *DB) DeleteSQL(table string, pk M) (string, []interface{}) {
	q := fmt.Sprintf("DELETE FROM %s", table)

	var (
		cols []string
		vals []interface{}
	)

	for k, v := range pk {
		cols = append(cols, fmt.Sprintf("%s=?", k))
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s WHERE %s", q, strings.Join(cols, " and "))
	return q, vals
}

// DeleteByPK 删除表数据
func (m *DB) DeleteByPK(table string, pk M) (int64, error) {
	q, vals := m.DeleteSQL(table, pk)
	result, err := m.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// DeleteByPKWithTran 使用事物删除表数据
func (m *DB) DeleteByPKWithTran(tran *gorp.Transaction, table string, pk M) (int64, error) {
	q, vals := m.DeleteSQL(table, pk)
	result, err := tran.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// In 组织带有IN查询的SQL
func (m *DB) In(query string, args ...interface{}) (string, []interface{}, error) {
	type argMeta struct {
		v      reflect.Value
		i      interface{}
		length int
	}

	var flatArgsCount int
	var anySlices bool

	meta := make([]argMeta, len(args))

	for i, arg := range args {
		v := reflect.ValueOf(arg)

		t := v.Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() == reflect.Slice {
			meta[i].length = v.Len()
			meta[i].v = v

			anySlices = true
			flatArgsCount += meta[i].length

			if meta[i].length == 0 {
				return "", nil, errors.New("empty slice passed to 'in' query")
			}
		} else {
			meta[i].i = arg
			flatArgsCount++
		}
	}

	if !anySlices {
		return query, args, nil
	}

	newArgs := make([]interface{}, 0, flatArgsCount)
	buf := bytes.NewBuffer(make([]byte, 0, len(query)+len(", ?")*flatArgsCount))

	var arg, offset int

	for i := strings.IndexByte(query[offset:], '?'); i != -1; i = strings.IndexByte(query[offset:], '?') {
		if arg >= len(meta) {
			return "", nil, errors.New("number of bindVars exceeds arguments")
		}

		argMeta := meta[arg]
		arg++

		if argMeta.length == 0 {
			offset = offset + i + 1
			newArgs = append(newArgs, argMeta.i)
			continue
		}

		buf.WriteString(query[:offset+i+1])

		for si := 1; si < argMeta.length; si++ {
			buf.WriteString(", ?")
		}

		newArgs = m.appendReflectSlice(newArgs, argMeta.v, argMeta.length)

		query = query[offset+i+1:]
		offset = 0
	}

	buf.WriteString(query)

	if arg < len(meta) {
		return "", nil, errors.New("number of bindVars less than number arguments")
	}

	return buf.String(), newArgs, nil
}

func (m *DB) appendReflectSlice(args []interface{}, v reflect.Value, vlen int) []interface{} {
	switch val := v.Interface().(type) {
	case []interface{}:
		args = append(args, val...)
	case []int:
		for i := range val {
			args = append(args, val[i])
		}
	case []string:
		for i := range val {
			args = append(args, val[i])
		}
	default:
		for si := 0; si < vlen; si++ {
			args = append(args, v.Index(si).Interface())
		}
	}

	return args
}
