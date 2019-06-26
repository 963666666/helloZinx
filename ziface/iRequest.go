package ziface

type IRequest interface {
	GetConnection() //获取请求的连接信息
	GetData()       //获取请求的数据
	GetMsgId()      //获取消息的Id
}
