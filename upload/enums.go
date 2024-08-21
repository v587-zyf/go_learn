package main

const (
	UP_TYPE_CLIENT = "c"
	UP_TYPE_SERVER = "s"
)

var (
	SKIP_FILE_NAME = map[string]struct{}{
		"image": {},
		"web":   {},
	}
)

const (
	SegmentSize = 3 * 1024 * 1024 // 传输文件512kb
)

const (
	CONN_TYPE_PASSWORD   = "pw"
	CONN_TYPE_PUBLIC_KEY = "pk"
)
