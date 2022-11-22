package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/spf13/cast"
	"io/ioutil"
	"mime/multipart"
)

// IRequest 代表请求的方法
type IRequest interface {
	//	请求地url带的参数
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	//	路由匹配中带的参数
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	//	form 表单中带的参数
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	//	json body
	BindJson(obj interface{}) error

	//	xml body
	BindXml(obj interface{}) error

	// GetRawData raw body
	GetRawData() ([]byte, error)

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

const defaultMultipartMemory = 32 << 20 // 32 MB

//#region query url

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query()) //注意这个的写法
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat32(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToBool(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1], true
		}
	}
	return def, false
}

func (ctx *Context) QueryStringSlice(key string, def []string) []string {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		return value
	}
	return def
}

func (ctx *Context) Query(key string) interface{} {
	params := ctx.QueryAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1]
		}

	}
	return nil
}

//#end query url

// #region Param
func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok && value != "" {
			return cast.ToInt(value), true
		}
	}
	return def, false
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok && value != "" {
			return cast.ToInt64(value), true
		}
	}
	return def, false
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok && value != "" {
			return cast.ToFloat64(value), true
		}
	}
	return def, false
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok && value != "" {
			return cast.ToFloat32(value), true
		}
	}
	return def, false
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok && value != "" {
			return cast.ToBool(value), true
		}
	}
	return def, false
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok {
			return value, true
		}
	}
	return def, false
}

func (ctx *Context) Param(key string) interface{} {
	if ctx.params != nil {
		if value, ok := ctx.params[key]; ok {
			return value
		}
	}
	return nil
}

//#end Param

// #region form post

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToInt64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat64(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToFloat32(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return cast.ToBool(value[l-1]), true
		}
	}
	return def, false
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1]
		}
	}
	return def
}

func (ctx *Context) FormStringSlice(key string, def []string) []string {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		return value
	}
	return def
}

func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil {
		if err := ctx.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := ctx.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil

}

func (ctx *Context) Form(key string) interface{} {
	params := ctx.FormAll()
	if value, ok := params[key]; ok {
		l := len(value)
		if l > 0 {
			return value[l-1]
		}
	}
	return nil
}

//#end form post

// #region application/json post

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// 重新填充 request.Body，为后续的逻辑二次读取做准备,request.Body 的读取是一次性的，读取一次之后，下个逻辑再去 request.Body 中是读取不到数据内容的
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) //?
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

//# end application/json post

//	xml body
func (ctx *Context) BindXml(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = xml.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// GetRawData raw body
func (ctx *Context) GetRawData() ([]byte, error) {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return nil, err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	} else {
		return nil, errors.New("ctx.request empty")
	}
	return nil, nil
}

//	基础信息
func (ctx *Context) Uri() string {
	return ctx.request.RequestURI
}
func (ctx *Context) Method() string {
	return ctx.request.Method
}
func (ctx *Context) Host() string {
	return ctx.request.URL.Host
}
func (ctx *Context) ClientIp() string {
	r := ctx.request
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
	return ctx.request.Header
}
func (ctx *Context) Header(key string) (string, bool) {
	values := ctx.request.Header.Values(key)
	if values == nil || len(values) < 0 {
		return "", false
	}
	return values[0], true
}

//cookie
func (ctx *Context) Cookies() map[string]string {
	cookies := ctx.request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}
func (ctx *Context) Cookie(key string) (string, bool) {
	cookieMap := ctx.Cookies()
	if cookie, ok := cookieMap[key]; ok {
		return cookie, true
	}
	return "", false
}
