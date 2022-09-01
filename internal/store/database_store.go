package store

import (
	"database/sql"
	"dnslog/internal/model"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseStore struct {
	db *sql.DB
}

func NewDatabaseStore(dbFile string) (Store, error) {
	var err error
	d := DatabaseStore{}
	d.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (d *DatabaseStore) Put(record model.Record) error {
	stmt, err := d.db.Prepare("insert into records(domain, ip, created_at) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(record.Domain, record.IP, record.CreatedAt)
	return err
}

func (m *DatabaseStore) Size(domain, ip string) (int64, error) {
	sql := "select count(*) from records"
	var conditions []string
	var args []any
	if domain != "" {
		conditions = append(conditions, "domain=?")
		args = append(args, domain)
	}
	if ip != "" {
		conditions = append(conditions, "ip=?")
		args = append(args, ip)
	}
	if len(conditions) > 0 {
		sql += " where " + strings.Join(conditions, " and ")
	}
	stmt, err := m.db.Prepare(sql)
	var size int64
	if err != nil {
		return size, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(args...).Scan(&size)
	return size, err
}

func (m *DatabaseStore) Paginate(offset, limit int, domain, ip string) ([]model.Record, error) {
	sqlStmt := "select * from records"
	var conditions []string
	var args []any
	if domain != "" {
		conditions = append(conditions, "domain=?")
		args = append(args, domain)
	}
	if ip != "" {
		conditions = append(conditions, "ip=?")
		args = append(args, ip)
	}
	if len(conditions) > 0 {
		sqlStmt += " where " + strings.Join(conditions, " and ")
	}
	sqlStmt += fmt.Sprintf(" order by id desc limit %d offset %d", limit, offset)
	stmt, err := m.db.Prepare(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var records []model.Record
	for rows.Next() {
		var id uint64
		var domain string
		var ip string
		var created_at string
		err = rows.Scan(&id, &domain, &ip, &created_at)
		if err != nil {
			return nil, err
		}
		records = append(records, model.Record{
			ID:        id,
			Domain:    domain,
			IP:        ip,
			CreatedAt: created_at,
		})
	}
	return records, nil
}

func (d *DatabaseStore) GetByDomain(domain string) (*model.Record, error) {
	stmt, err := d.db.Prepare("select * from records where domain=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(domain)
	var id uint64
	var domainStr string
	var ip string
	var created_at string
	err = row.Scan(&id, &domainStr, &ip, &created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model.Record{
		ID:        id,
		Domain:    domainStr,
		IP:        ip,
		CreatedAt: created_at,
	}, nil
}
