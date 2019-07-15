package znet

import (
	"fmt"
	"helloZinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: 10,
		TaskQueue:      make([]chan ziface.IRequest, 10),
	}
}

func (msgHandler *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := msgHandler.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), "is not found")
		return
	}

	//执行对应的方法
	handler.BeforeHandler(request)
	handler.Handler(request)
	handler.AfterHandler(request)
}

//启动一个Worker工作流程
func (msgHandler *MsgHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("workerId is ", workerId, "is started")
	for {
		select {
		case request := <-taskQueue:
			msgHandler.DoMsgHandler(request)
		}
	}
}

func (msgHandler *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(msgHandler.WorkerPoolSize); i ++ {

		msgHandler.TaskQueue[i] = make(chan ziface.IRequest, 100)

		go msgHandler.StartOneWorker(i, msgHandler.TaskQueue[i])
	}
}

func (msgHandler *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerId := request.GetConnection().GetConnId() % msgHandler.WorkerPoolSize
	msgHandler.TaskQueue[workerId] <- request
}
