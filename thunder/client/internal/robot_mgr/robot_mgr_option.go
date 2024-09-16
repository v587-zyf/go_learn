package robot_mgr

type RobotOption struct {
	addr string

	pem string
	key string

	https bool
}

type Option func(opts *RobotOption)

func NewRobotOption() *RobotOption {
	o := &RobotOption{}

	return o
}

func WithAddr(addr string) Option {
	return func(opts *RobotOption) {
		opts.addr = addr
	}
}

func WithHttps(https bool) Option {
	return func(opts *RobotOption) {
		opts.https = https
	}
}
func WithPem(pem string) Option {
	return func(opts *RobotOption) {
		opts.pem = pem
	}
}
func WithKey(key string) Option {
	return func(opts *RobotOption) {
		opts.key = key
	}
}
