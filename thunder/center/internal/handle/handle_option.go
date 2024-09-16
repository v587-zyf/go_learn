package handle

type HandleOption struct {
	Dev bool
	SID int64
}

type Option func(opts *HandleOption)

func NewHandleOption() *HandleOption {
	o := &HandleOption{}

	return o
}

func WithDev(Dev bool) Option {
	return func(opts *HandleOption) {
		opts.Dev = Dev
	}
}
func WithSID(SID int64) Option {
	return func(opts *HandleOption) {
		opts.SID = SID
	}
}
