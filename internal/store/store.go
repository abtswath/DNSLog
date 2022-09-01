package store

import "dnslog/internal/model"

type Store interface {
	GetByDomain(domain string) (*model.Record, error)
	Size(domain, ip string) (int64, error)
	Paginate(offset, limit int, domain, ip string) ([]model.Record, error)
	Put(model.Record) error
}
