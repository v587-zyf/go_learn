Head	总共 18 字节

- 0-4		len			消息的总长度	uint32
- 4-6		msgID		消息id  *_msgId.proto文件定义  uint16
- 6-10		tag			预留位 方便后续扩展  uint32
- 10-18	userID		用户id    uint64
- 18-...	msg			具体消息

Msg

- maxSize	65535 * 5
