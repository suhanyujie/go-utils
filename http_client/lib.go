package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	nacosModel "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/suhanyujie/go_utils/jsonx"
	"github.com/suhanyujie/go_utils/logx"
	"github.com/suhanyujie/go_utils/times"
)

const (
	defaultContentType = "application/json"
	BlankString        = ""
)

var (
	Logger *logx.Logx
)

var httpClient = &http.Client{}

// 可以手动设定超时时间的 http 客户端
var httpClientWithTimeout = &http.Client{}

type HeaderOption struct {
	Name  string
	Value string
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(3) * time.Second,
	}
}

func Init(ctx context.Context) {
	v, ok := ctx.Value("logger").(*logx.Logx)
	if ok {
		SetLogger(v)
	}
}

func SetLogger(logger *logx.Logx) {
	if Logger == nil {
		Logger = logger
	}
}

func Post(url string, params map[string]interface{}, body string, headerOptions ...HeaderOption) (string, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(body))
	req.Header.Set("Content-Type", defaultContentType)

	if err != nil {
		return BlankString, 0, errors.Wrap(err, "http Post request error")
	}

	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	headers := jsonx.ToJsonIgnoreErr(req.Header)
	Logger.Printf("http type: POST|request [%s] starting|request body [%s]|request headers [%s]", fullUrl, body, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)

	//截取下日志长度
	cutRespBody := respBody
	if len(cutRespBody) > 1000 {
		cutRespBody = cutRespBody[:1000]
	}
	Logger.Printf("[Post] http type: POST| request [%s] successful| request body [%s]|request headers [%s]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, cutRespBody, timeConsuming)
	return respBody, httpCode, err
}

func Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) (string, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("GET", fullUrl, nil)

	if err != nil {
		return BlankString, 0, errors.Wrap(err, "http Get request error")
	}

	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}
	headers := jsonx.ToJsonIgnoreErr(req.Header)
	Logger.Printf("http type: GET|request [%s] starting|request headers [%s]", fullUrl, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)
	cutRespBody := respBody
	if len(cutRespBody) > 1000 {
		cutRespBody = cutRespBody[:1000]
	}
	Logger.Printf("http type: GET| request [%s] successful|request headers [%s]|response status code [%d]|response body [%s]| time-consuming [%s]", fullUrl, headers, httpCode, cutRespBody, timeConsuming)
	return respBody, httpCode, err
}

func responseHandle(resp *http.Response, err error) (string, int, error) {
	if err != nil {
		Logger.Println(err)
		return "", 0, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Println(err)
		return "", resp.StatusCode, err
	}
	return string(b), resp.StatusCode, nil
}

func Float64Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

func Float64Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

// 不可存在结构体值，否则会报错。
func ConvertToQueryParams(params map[string]interface{}) string {
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		vStr := SerialToString(v)
		buffer.WriteString(fmt.Sprintf("%s=%v&", k, vStr))
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}

func ConvertToQueryParamsOld(params map[string]interface{}) string {
	// 反序列化是为了更加统一地处理，但可能会带来一些问题。
	paramsJson := jsonx.ToJsonIgnoreErr(params)
	params = map[string]interface{}{}
	_ = jsonx.FromJson(paramsJson, &params)

	if &params == nil || len(params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		if v == nil {
			continue
		}
		if fv, ok := v.(float64); ok {
			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatFloat(fv, 'f', -1, 64)))
		} else if fv, ok := v.(float32); ok {
			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatFloat(float64(fv), 'f', -1, 32)))
		} else {
			buffer.WriteString(fmt.Sprintf("%s=%v&", k, v))
		}
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}

func PostWithTimeout(url string, params map[string]interface{}, body string, timeout uint32, headerOptions ...HeaderOption) (string, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(body))
	req.Header.Set("Content-Type", defaultContentType)

	if err != nil {
		return BlankString, 0, err
	}

	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	headers := jsonx.ToJsonIgnoreErr(req.Header)
	Logger.Printf("http type: POST|request [%s] starting|request body [%s]|request headers [%s]", fullUrl, body, headers)

	start := times.GetNowMillisecond()
	resp, err := getHttpClientWithTimeout(timeout).Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)

	Logger.Printf("http type: POST| request [%s] successful| request body [%s]|request headers [%s]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, respBody, timeConsuming)
	return respBody, httpCode, err
}

// 设定超时时间并返回 http 客户端
// timeoutOpt 表示设定的超时时间的秒数
func getHttpClientWithTimeout(timeoutOpt uint32) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClientWithTimeout = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeoutOpt) * time.Second,
	}
	return httpClientWithTimeout
}

func isNil(dest interface{}) bool {
	if dest == nil {
		return true
	}
	v := reflect.ValueOf(dest)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// 将数据转化为字符串
func SerialToString(dest interface{}) string {
	var key string
	if isNil(dest) {
		return key
	}
	switch dest.(type) {
	case float64:
		key = decimal.NewFromFloat(dest.(float64)).String()
	case *float64:
		if dest.(*float64) != nil {
			key = decimal.NewFromFloat(*dest.(*float64)).String()
		}
	case float32:
		key = decimal.NewFromFloat32(dest.(float32)).String()
	case *float32:
		if dest.(*float32) != nil {
			key = decimal.NewFromFloat32(*dest.(*float32)).String()
		}
	case int:
		key = strconv.Itoa(dest.(int))
	case *int:
		if dest.(*int) != nil {
			key = strconv.Itoa(*dest.(*int))
		}
	case uint:
		key = strconv.Itoa(int(dest.(uint)))
	case *uint:
		key = strconv.Itoa(int(*dest.(*uint)))
	case int8:
		key = strconv.Itoa(int(dest.(int8)))
	case *int8:
		if dest.(*int8) != nil {
			key = strconv.Itoa(int(*dest.(*int8)))
		}
	case uint8:
		key = strconv.Itoa(int(dest.(uint8)))
	case *uint8:
		if dest.(*uint8) != nil {
			key = strconv.Itoa(int(*dest.(*uint8)))
		}
	case int16:
		key = strconv.Itoa(int(dest.(int16)))
	case *int16:
		if dest.(*int16) != nil {
			key = strconv.Itoa(int(*dest.(*int16)))
		}
	case uint16:
		key = strconv.Itoa(int(dest.(uint16)))
	case *uint16:
		if dest.(*uint16) != nil {
			key = strconv.Itoa(int(*dest.(*uint16)))
		}
	case int32:
		key = strconv.Itoa(int(dest.(int32)))
	case *int32:
		if dest.(*int32) != nil {
			key = strconv.Itoa(int(*dest.(*int32)))
		}
	case uint32:
		key = strconv.Itoa(int(dest.(uint32)))
	case *uint32:
		if dest.(*uint32) != nil {
			key = strconv.Itoa(int(*dest.(*uint32)))
		}
	case int64:
		key = strconv.FormatInt(dest.(int64), 10)
	case *int64:
		if dest.(*int64) != nil {
			key = strconv.FormatInt(*dest.(*int64), 10)
		}
	case uint64:
		key = strconv.FormatUint(dest.(uint64), 10)
	case *uint64:
		if dest.(*uint64) != nil {
			key = strconv.FormatUint(*dest.(*uint64), 10)
		}
	case string:
		key = dest.(string)
	case *string:
		if dest.(*string) != nil {
			key = *dest.(*string)
		}
	case []byte:
		key = string(dest.([]byte))
	case *[]byte:
		if dest.(*[]byte) != nil {
			key = string(*dest.(*[]byte))
		}
	case bool:
		if dest.(bool) {
			key = "true"
		} else {
			key = "false"
		}
	case *bool:
		if dest.(*bool) != nil {
			if *dest.(*bool) {
				key = "true"
			} else {
				key = "false"
			}
		}
	default:
	}
	return key
}

func Put(url string, params map[string]interface{}, body string, headerOptions ...HeaderOption) (string, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(body))
	req.Header.Set("Content-Type", defaultContentType)

	if err != nil {
		return BlankString, 0, errors.Wrap(err, "http Put request error")
	}

	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	headers := jsonx.ToJsonIgnoreErr(req.Header)
	Logger.Printf("http type: Put|request [%s] starting|request body [%s]|request headers [%s]", fullUrl, body, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)

	//截取下日志长度
	cutRespBody := respBody
	if len(cutRespBody) > 500 {
		cutRespBody = cutRespBody[:500]
	}
	Logger.Printf("http type: Put| request [%s] successful| request body [%s]|request headers [%s]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, cutRespBody, timeConsuming)
	return respBody, httpCode, err
}

func Delete(url string, params map[string]interface{}, body string, headerOptions ...HeaderOption) (string, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("DELETE", fullUrl, strings.NewReader(body))
	req.Header.Set("Content-Type", defaultContentType)

	if err != nil {
		return BlankString, 0, errors.Wrap(err, "http Delete request error")
	}
	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	headers := jsonx.ToJsonIgnoreErr(req.Header)
	Logger.Printf("http type: Delete|request [%s] starting|request body [%s]|request headers [%s]", fullUrl, body, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)

	//截取下日志长度
	cutRespBody := respBody
	if len(cutRespBody) > 500 {
		cutRespBody = cutRespBody[:500]
	}
	Logger.Printf("http type: DELETE| request [%s] successful| request body [%s]|request headers [%s]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, cutRespBody, timeConsuming)
	return respBody, httpCode, err
}

func PostForNacos(nacosModel *nacosModel.Service, url string, params map[string]interface{}, body string, headerOptions ...HeaderOption) (string, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(body))
	req.Header.Set("Content-Type", defaultContentType)

	if err != nil {
		return BlankString, 0, errors.Wrap(err, "http Post request error")
	}

	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	headers := jsonx.ToJsonIgnoreErr(req.Header)
	Logger.Printf("http type: POST|request [%s] starting|request body [%s]|request headers [%s]", fullUrl, body, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)

	//截取下日志长度
	cutRespBody := respBody
	if len(cutRespBody) > 1000 {
		cutRespBody = cutRespBody[:1000]
	}
	Logger.Printf("[Post] http type: POST| request [%s] successful| request body [%s]|request headers [%s]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, cutRespBody, timeConsuming)
	return respBody, httpCode, err
}
