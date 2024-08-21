package enums

import "time"

const (
	READ_BUFF_SIZE_INIT  int = 8 * 1024
	READ_BUFF_SIZE_MAX   int = 64 * 1024 * 10
	WRITE_BUFF_SIZE_INIT int = 16 * 1024
)

const (
	MSG_HEADER_SIZE = 18

	MSG_MAX_PACKET_SIZE = 65535 * 5

	CONN_WRITE_WAIT_TIME = 10 * time.Second

	HEARTBEAT_INTERVAL = 10 * time.Second
)
