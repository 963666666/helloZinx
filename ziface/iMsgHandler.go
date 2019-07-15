package ziface

type IMsgHandler interface {
	//DoMsgHandler(request IRequest)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest) //将消息给TaskQueue, 由worker进行处理
}


