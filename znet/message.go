package znet

type Message struct {
	DataLen uint32
	Id uint32
	Data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		Id: id,
		Data: data,
	}
}

func (message *Message) GetDataLen() uint32 {
	return message.DataLen
}


func (message *Message) GetMsgId() uint32 {
	return message.Id
}


func (message *Message) GetData() []byte {
	return message.Data
}


func (message *Message) SetMsgId(msgId uint32) {
	message.Id = msgId
}


func (message *Message) SetDataLen(msgLen uint32) {
	message.DataLen = msgLen
}


func (message *Message) SetData(msgData []byte) {
	message.Data = msgData
}