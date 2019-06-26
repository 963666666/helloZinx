package ziface

type IRouter interface {
	BeforeHandler()
	Handler()
	AfterHandler()
}
