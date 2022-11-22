package contract

import "net/http"

// KernelKey 提供kernel 服务凭证
const KernelKey = "hade:kernel"

// 服务接口

type Kernel interface {
	HttpEngine() http.Handler
}
