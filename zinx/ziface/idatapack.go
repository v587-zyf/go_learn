package ziface

type IDataPack interface {
	// 获取包头
	GetHeadLen() uint32
	// 封包
	Pack(msg IMessage) ([]byte, error)
	// 拆包
	Unpack([]byte) (IMessage, error)
}
