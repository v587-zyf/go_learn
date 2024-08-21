package grpc_server_stream_mgr

import (
	"google.golang.org/grpc"
)

var defGrpcStreamClientMgr *GrpcServerStreamMgr

func InitGrpcClientStream(opts ...any) (err error) {
	defGrpcStreamClientMgr = NewGrpcClientStream()
	if err = defGrpcStreamClientMgr.Init(opts...); err != nil {
		return err
	}

	return nil
}

func Add(st int32, id uint64, stream grpc.ServerStream) {
	defGrpcStreamClientMgr.Add(st, id, stream)
}

func Del(st int32, id uint64) {
	defGrpcStreamClientMgr.Del(st, id)
}

func GetStreamByType(st int32) map[uint64]grpc.ServerStream {
	return defGrpcStreamClientMgr.GetStreamByType(st)
}

func RandStreamByType(st int32) grpc.ServerStream {
	return defGrpcStreamClientMgr.RandStreamByType(st)
}
