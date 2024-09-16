package handle

type Handle struct {
	options *HandleOption
}

var h *Handle

func init() {
	h = &Handle{
		options: NewHandleOption(),
	}
}

func GetHandle() *Handle {
	return h
}

func GetHandleOps() *HandleOption {
	return h.options
}

func HandleInit(opts ...Option) error {
	for _, opt := range opts {
		opt(h.options)
	}

	return nil
}
