package model

type Record struct {
	ID        uint64 `json:"id"`
	IP        string `json:"ip"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
}
