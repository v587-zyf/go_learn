package module

type ModuleMgrOption struct {
	Dev bool
	SID int64
}

type Option func(opts *ModuleMgrOption)

func NewModuleMgrOption() *ModuleMgrOption {
	o := &ModuleMgrOption{}

	return o
}

func WithDev(Dev bool) Option {
	return func(opts *ModuleMgrOption) {
		opts.Dev = Dev
	}
}
func WithSID(SID int64) Option {
	return func(opts *ModuleMgrOption) {
		opts.SID = SID
	}
}
