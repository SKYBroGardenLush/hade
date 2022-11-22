package gin

import (
	"github.com/spf13/cast"
	"mime/multipart"
)

// IRequest 代表请求的方法
type IRequest interface {
	//	请求地url带的参数
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)
	DefaultQuery(key string) interface{}

	//	路由匹配中带的参数
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	paramstring(key string, def string) (string, bool)
	DefaultParam(key string) interface{}

	//	DefaultForm 表单中带的参数
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	DefaultFormFile(key string) (*multipart.FileHeader, error)
	DefaultForm(key string) interface{}

	//	基础信息
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	//	header
	Headers() map[string][]string
	Header(key string) (string, bool)

	//cookie
	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

//#region DefaultQuery url

func (ctx *Context) QueryAll() map[string][]string {
	ctx.initQueryCache()
	return map[string][]string(ctx.queryCache)
}

func (ctx *Context) DefaultQueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat32(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToBool(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1], true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryStringSlice(key string, def []string) []string {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		return value
	}
	return def
}

//#end DefaultQuery url

// #region DefaultParam

func (ctx *Context) DefaultParamInt(key string, def int) (int, bool) {
	if val := ctx.HadeParam(key); val != nil {
		// 通过cast进行类型转换
		return cast.ToInt(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	if val := ctx.HadeParam(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	if val := ctx.HadeParam(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	if val := ctx.HadeParam(key); val != nil {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	if val := ctx.HadeParam(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (ctx *Context) DefaultParamString(key string, def string) (string, bool) {
	if val := ctx.HadeParam(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

// 获取路由参数
func (ctx *Context) HadeParam(key string) interface{} {
	if val, ok := ctx.Params.Get(key); ok {
		return val
	}
	return nil
}

//#end DefaultParam

// #region DefaultForm post

func (ctx *Context) FormAll() map[string][]string {
	ctx.initFormCache()
	return map[string][]string(ctx.formCache)
}

func (ctx *Context) DefaultFormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat32(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToBool(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) DefaultFormString(key string, def string) string {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1]
		}
	}
	return def
}

func (ctx *Context) DefaultFormStringSlice(key string, def []string) []string {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		return value
	}
	return def
}

func (ctx *Context) DefaultFormFile(key string) (*multipart.FileHeader, error) {
	if ctx.Request.MultipartForm == nil {
		if err := ctx.Request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := ctx.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

func (ctx *Context) DefaultForm(key string) interface{} {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1]
		}
	}
	return nil
}

//#end DefaultForm post

//	基础信息
func (ctx *Context) Uri() string {
	return ctx.Request.RequestURI
}
func (ctx *Context) Method() string {
	return ctx.Request.Method
}
func (ctx *Context) Host() string {
	return ctx.Request.URL.Host
}
func (ctx *Context) ClientIp() string {
	r := ctx.Request
	ipAddress := r.Header.Get("X-Real-IP")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

//	header
func (ctx *Context) Headers() map[string][]string {
	return ctx.Request.Header
}

//cookie
func (ctx *Context) Cookies() map[string]string {
	cookies := ctx.Request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}
