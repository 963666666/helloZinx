package ziface

type IRouter interface {
	BeforeHandler(request IRequest)
	Handler(request IRequest)
	AfterHandler(request IRequest)
}
