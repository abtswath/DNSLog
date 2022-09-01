package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Option struct {
	HTTPAddr string       `yaml:"http-addr"`
	DNSAddr  string       `yaml:"dns-addr"`
	Domain   string       `yaml:"domain"`
	A        [4]byte      `yaml:"a"`
	Username string       `yaml:"username"`
	Password string       `yaml:"password"`
	Tokens   []string     `yaml:"tokens"`
	LogLevel logrus.Level `yaml:"log-level"`
}

type DefaultOptionSetFunc = func(*Option)

func WithHTTPAddr(addr string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.HTTPAddr == "" {
			o.HTTPAddr = addr
		}
	}
}

func WithUDPAddr(addr string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.DNSAddr == "" {
			o.DNSAddr = addr
		}
	}
}

func WithDomain(domain string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.Domain == "" {
			o.Domain = domain
		}
	}
}

func WithUsername(username string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.Username == "" {
			o.Username = username
		}
	}
}

func WithPassword(password string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.Password == "" {
			o.Password = password
		}
	}
}

func WithA(a [4]byte) DefaultOptionSetFunc {
	var initialValue [4]byte
	return func(o *Option) {
		if o.A == initialValue {
			o.A = a
		}
	}
}

func WithLogLevel(level logrus.Level) DefaultOptionSetFunc {
	return func(o *Option) {
		valid := false
		for _, v := range []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		} {
			if o.LogLevel == v {
				valid = true
				break
			}
		}
		if o.LogLevel == 0 || !valid {
			o.LogLevel = level
		}
	}
}

var (
	DefaultHTTPAddr = ":8080"
	DefaultUDPAddr  = ":53"
	DefaultA        = [4]byte{127, 0, 0, 1}
	DefaultUsername = "admin"
	DefaultPassword = "dnslog2022"
)

func NewDaemonOptions(dir string) Option {
	opts := Option{}
	_ = viper.Unmarshal(opts)
	for _, f := range []DefaultOptionSetFunc{
		WithHTTPAddr(DefaultHTTPAddr),
		WithUDPAddr(DefaultUDPAddr),
		WithA(DefaultA),
		WithUsername(DefaultUsername),
		WithPassword(DefaultPassword),
		WithLogLevel(logrus.InfoLevel),
	} {
		f(&opts)
	}
	if opts.Domain == "" {
		logrus.Warning("The domain name is not set, the system will record all DNS resolution records.")
	}
	return opts
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.dnslog")
}
