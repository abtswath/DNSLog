package dns

import (
	"dnslog/internal/model"
	"dnslog/internal/store"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/dns/dnsmessage"
)

type DNSServer struct {
	Domain  string
	Addr    string
	udpAddr *net.UDPAddr
	Store   store.Store
	conn    *net.UDPConn
	A       [4]byte
}

func NewServer(o Option, store store.Store) *DNSServer {
	return &DNSServer{
		Domain: o.Domain,
		Addr:   o.Addr,
		Store:  store,
		A:      o.A,
	}
}

func (d *DNSServer) Serve() error {
	var err error
	d.udpAddr, err = net.ResolveUDPAddr("udp", d.Addr)
	if err != nil {
		return err
	}
	d.conn, err = net.ListenUDP("udp", d.udpAddr)
	if err != nil {
		return err
	}
	defer d.conn.Close()
	logrus.Infof("Listen and serving DNS on address %s", d.udpAddr)
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
		logrus.Debugf("Could not unpack the message: %v", err)
		return
	}
	if len(dnsMessage.Questions) < 1 {
		return
	}
	defer d.response(addr, dnsMessage)
	question := dnsMessage.Questions[0]
	queryName := question.Name.String()
	queryName = queryName[:len(queryName)-1]
	if d.Domain == "" || strings.HasSuffix(queryName, d.Domain) {
		d.Store.Put(model.Record{
			IP:        addr.IP.String(),
			Domain:    queryName,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
}

func (d *DNSServer) response(addr *net.UDPAddr, message dnsmessage.Message) {
	message.Response = true
	message.Answers = append(message.Answers, d.newResource(message))
	response, err := message.Pack()
	if err != nil {
		logrus.Debugf("make response failed: %v\n", err)
		return
	}
	if _, err := d.conn.WriteToUDP(response, addr); err != nil {
		logrus.Debugf("UDP write failed: %v\n", err)
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
