package store

import "golang.org/x/net/dns/dnsmessage"

type Store interface {
	Get() []Record
	GetByIP() []Record
	GetByDomain() []Record
	Paginate(offset, limit uint, domain, ip string) []Record
	Put(string, dnsmessage.Question)
}
