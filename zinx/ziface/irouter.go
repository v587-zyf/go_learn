package ziface

type IRouter interface {
	// 处理业务之前的hook
	PreHandle(request IRequest)
	// 处理业务的hook
	Handle(request IRequest)
	// 处理业务之后的hook
	PostHandle(request IRequest)
}
