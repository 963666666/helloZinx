package znet

import (
	"bytes"
	"encoding/binary"
	"helloZinx/ziface"
)

type DataPack struct {
	
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack)GetDataLen() uint32 {
	return 8
}

func (dp *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	err := binary.Write(dataBuff, binary.LittleEndian, message.GetDataLen())
	if err != nil {
		return nil, err
	}

	//写msgID
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetMsgId())
	if err != nil {
		return nil, err
	}

	//写data数据
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetData())
	if err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewBuffer(binaryData)

	//直接压message的信息，得到dataLen和msgId
	msg := &Message{}

	//读dataLen的信息
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	return msg, nil
}