package ziface

type IMessage interface {
	 GetDataLen() uint32
	 GetMsgId() uint32
	 GetData() []byte

	 SetMsgId(msgId uint32)
	 SetDataLen(msgLen uint32)
	 SetData(msgData []byte)
}