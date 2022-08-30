package dns

type Option struct {
	Domain  string
	Address string
	A       [4]byte
}

type DefaultOptionSetFunc = func(*Option)

func WithDomain(domain string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.Domain == "" {
			o.Domain = domain
		}
	}
}

func WithAddress(address string) DefaultOptionSetFunc {
	return func(o *Option) {
		if o.Address == "" {
			o.Address = address
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
