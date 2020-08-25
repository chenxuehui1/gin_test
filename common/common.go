package common

import (
	"encoding/json"
	"github.com/chenxuehui1/gin_test/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)
//公共方法
var (
	// router
	router http.Handler

	// 为令牌授权等定制的请求头
	myHeaders = make(map[string]string, 0)

	logging *log.Logger
)

// 设置路由器
func SetRouter(r http.Handler) {
	router = r
}

// 设置日志
func SetLog(l *log.Logger) {
	logging = l
}

// 添加自定义请求头
func AddHeader(key, value string) {
	myHeaders[key] = value
}

// 打印日志
func printfLog(format string, v ...interface{}) {
	if logging == nil {
		return
	}

	logging.Printf(format, v...)
}
//调用处理器
func invokeHandler(req *http.Request) (bodyByte []byte, err error) {

	//初始化response记录
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 从响应记录中提取响应
	result := w.Result()
	defer result.Body.Close()

	// 提取响应主体
	bodyByte, err = ioutil.ReadAll(result.Body)

	return
}

func TestFileHandler(method, api, fileName string, fieldName string, param interface{}) (bodyByte []byte, err error) {
	// 检查路由器是否为nil
	if router == nil {
		err = utils.ErrRouterNotSet
		return
	}

	paramStr := utils.MakeQueryStrFrom(param)
	printfLog("TestFileHandler\tRequest:\t%v:%v?%v \tFileName:%v, FieldName:%v\n", method, api, paramStr, fileName, fieldName)

	// make request
	req, err := utils.MakeFileRequest(method, api, fileName, fieldName, param)
	if err != nil {
		return
	}

	for key, value := range myHeaders {
		req.Header.Add(key, value)
	}

	// 调用处理器
	bodyByte, err = invokeHandler(req)
	//printfLog("TestFileHandler\n\tRequestMethod:%v\n\tRequestUrl:%v\n\trequestBody:%v\n\tFileName:%v\n\tFieldName:%v\n\tResponse:%v\n\n\n", method, api, param,fileName,fieldName,string(bodyByte))
	printfLog("TestFileHandler\tResponse:\t%v:%v,\tResponse:%v\n\n\n", method, api, string(bodyByte))
	return
}

func TestOrdinaryHandler(method string, api string, mime string, param interface{}) (bodyByte []byte, err error) {
	if router == nil {
		err = utils.ErrRouterNotSet
		return
	}

	//printfLog("TestOrdinaryHandler\n\tRequestMethod:%v\n\tRequestUrl:%v\n\trequestBody:%v\n", method, api, param)

	// make request
	req, err := utils.MakeRequest(method, mime, api, param)
	if err != nil {
		return
	}

	// add the customed headers
	for key, value := range myHeaders {
		req.Header.Add(key, value)
	}

	// invoke handler
	bodyByte, err = invokeHandler(req)
	printfLog("TestOrdinaryHandler\n\tRequestMethod:%v\n\tRequestUrl:%v\n\trequestBody:%v\n\tResponse:%v\n\n\n", method, api, param,string(bodyByte))

	//printfLog("TestOrdinaryHandler\n\tResponseMethod:%v\n\tResponseUrl:%v\n\tResponseBody:%v\n\n\n", method, api, string(bodyByte))
	return
}

func TestHandlerUnMarshalResp(method string, uri string, way string, param interface{}, resp interface{}) error {
	bodyByte, err := TestOrdinaryHandler(method, uri, way, param)
	if err != nil {
		return err
	}
	//把json格式转换成resp对象,
	return json.Unmarshal(bodyByte, resp)
	//printfLog("TestHandlerUnMarshalResp\n\tRequestMethod:%v\n\tRequestUrl:%v\n\trequestBody:%v\n\tResponse:%v\n\n\n", method, uri, param,string(bodyByte))
}

func TestFileHandlerUnMarshalResp(method, uri, fileName string, filedName string, param interface{}, resp interface{}) error {
	bodyByte, err := TestFileHandler(method, uri, fileName, filedName, param)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyByte, resp)
}
