package ziface

import "hello_zinx/ziface"

type IMsgHandler interface {
	DoMsgHandler(request ziface.IRequest)
}


