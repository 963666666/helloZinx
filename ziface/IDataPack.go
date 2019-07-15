package ziface

type IDataPack interface {
	GetDataLen() uint32                //获取包头长度方法
	Pack(msg IMessage) ([]byte, error) //封包方法
	UnPack([]byte) (IMessage, error)           //拆包方法
}
