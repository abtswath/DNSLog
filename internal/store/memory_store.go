package store

import (
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Record struct {
	IP        string `json:"ip"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
}

type MemoryStore struct {
	collection []Record
	mutex      sync.Mutex
}

func (m *MemoryStore) Put(ip string, question dnsmessage.Question) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.collection = append(m.collection, Record{
		IP:        ip,
		Domain:    question.Name.String(),
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (m *MemoryStore) Get() []Record {
	return m.collection
}

func (m *MemoryStore) GetByIP(ip string) []Record {
	var result []Record
	for _, record := range m.collection {
		if record.IP == ip {
			result = append(result, record)
		}
	}
	return result
}

func (m *MemoryStore) GetByDomain(domain string) []Record {
	var result []Record
	for _, record := range m.collection {
		if record.Domain == domain {
			result = append(result, record)
		}
	}
	return result
}

func (m *MemoryStore) Paginate(offset, limit uint, domain, ip string) []Record {
	if domain != "" || ip != "" {
		var result []Record
		for _, record := range m.collection {
			if domain == "" {
				if record.Domain == domain {
					result = append(result, record)
				}
			}
			if ip == "" {
				if record.IP == ip {
					result = append(result, record)
				}
			}
		}
		return result[offset : offset+limit]
	}
	return m.collection[offset : offset+limit]
}
