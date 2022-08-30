package dns

import (
	"dnslog/internal/store"
	"log"
	"net"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

type DNSServer struct {
	Domain string
	Addr   *net.UDPAddr
	Logger *log.Logger
	Store  store.Store
	conn   *net.UDPConn
	A      [4]byte
}

func NewServer(o Option, store store.Store) (*DNSServer, error) {
	addr, err := net.ResolveUDPAddr("udp", o.Address)
	if err != nil {
		return nil, err
	}
	return &DNSServer{
		Domain: o.Domain,
		Addr:   addr,
		Logger: log.Default(),
		Store:  store,
		A:      o.A,
	}, nil
}

func (d *DNSServer) ListenAndServe() error {
	var err error
	d.conn, err = net.ListenUDP("udp", d.Addr)
	if err != nil {
		return err
	}
	for {
		message := make([]byte, 512)
		_, addr, err := d.conn.ReadFromUDP(message)
		if err != nil {
			return err
		}
		go d.serve(addr, message)
	}
}

func (d *DNSServer) serve(addr *net.UDPAddr, message []byte) {
	var dnsMessage dnsmessage.Message
	err := dnsMessage.Unpack(message)
	if err != nil {
		d.Logger.Printf("cannot pack the message: %v\n", err)
		return
	}
	if len(dnsMessage.Questions) < 1 {
		return
	}
	defer d.response(dnsMessage)
	question := dnsMessage.Questions[0]
	queryName := question.Name.String()
	if strings.HasSuffix(queryName, d.Domain) {
		d.Store.Put(addr.IP.String(), question)
	}
}

func (d *DNSServer) response(message dnsmessage.Message) {
	message.Response = true
	message.Answers = append(message.Answers, d.newResource(message))
	response, err := message.Pack()
	if err != nil {
		d.Logger.Printf("make response failed: %v\n", err)
		return
	}
	if _, err := d.conn.WriteToUDP(response, d.Addr); err != nil {
		d.Logger.Printf("UDP write failed: %v\n", err)
	}
}

func (d *DNSServer) newResource(message dnsmessage.Message) dnsmessage.Resource {
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  message.Questions[0].Name,
			Class: dnsmessage.ClassINET,
			TTL:   0,
		},
		Body: &dnsmessage.AResource{
			A: d.A,
		},
	}
}

func (d *DNSServer) Close() {
	d.conn.Close()
}
