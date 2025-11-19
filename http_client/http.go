package http_client

import (
	"errors"
	"fmt"
	"log"

	"github.com/suhanyujie/go_utils/encoding"
	"github.com/suhanyujie/go_utils/unsafex"
)

type DemoQueryParam1 struct {
	Token string `json:"token"`
}

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// DemoQueryResp1 定义接口返回值的结构
type DemoQueryResp1 struct {
	ErrorInfo error
	Data      *DemoQueryResp1Data `json:"data"`
}

type DemoQueryResp1Data struct {
	IsOk bool `json:"isOk"`
}

func DemoApiQuery(req DemoQueryParam1) *DemoQueryResp1 {
	respVo := &DemoQueryResp1{Data: &DemoQueryResp1Data{}}

	// test url
	reqUrl := fmt.Sprintf("%s%s/column/delete", "http://127.0.0.1", "/api")
	respVo.ErrorInfo = Request(reqUrl, req.Token, nil, respVo.Data)
	return respVo
}

func Request(reqUrl string, token string, req, respVo interface{}) error {
	queryParams := map[string]interface{}{}
	requestBody := marshalToString(req)
	fullUrl := reqUrl + ConvertToQueryParams(queryParams)
	fullUrl += "|" + requestBody
	extraHeaders := []HeaderOption{
		{Name: AppPlatformApiHeadTokenKey, Value: token},
	}

	respBody, respStatusCode, err := Post(reqUrl, queryParams, requestBody, extraHeaders...)
	//Process the response
	if err != nil {
		log.Printf("request [%s] failed, response status code [%d], err [%v]", fullUrl, respStatusCode, err)
		return err
	}
	//接口响应错误
	if respStatusCode < 200 || respStatusCode > 299 {
		respObj := Response{
			Code: -1,
			Msg:  "请求服务接口异常",
		}
		if len(respBody) != 0 {
			_ = unmarshalFromString(respBody, &respObj)
		}
		log.Printf("request [%s] failed , response status code [%d], err [%v]", fullUrl, respStatusCode, respObj.Msg)
		message := fmt.Sprintf("http code: %d", respStatusCode)
		if respObj.Msg != "" {
			message = respObj.Msg
		}
		return errors.New(message)
	}
	if len(respBody) > 0 {
		jsonConvertErr := unmarshalFromString(respBody, respVo)
		if jsonConvertErr != nil {
			return jsonConvertErr
		}
	}

	return nil
}

func RequestGet(reqUrl string, token string, req, respVo interface{}) error {
	extraHeaders := []HeaderOption{
		{Name: AppPlatformApiHeadTokenKey, Value: token},
	}
	respBody, _, err := Get(reqUrl, nil, extraHeaders...)
	if err != nil {
		return err
	}
	jsonConvertErr := unmarshalFromString(respBody, respVo)
	if jsonConvertErr != nil {
		return jsonConvertErr
	}
	return nil
}

func marshalToString(v interface{}) string {
	bts, _ := encoding.GetJsonCodec().Marshal(v)
	return unsafex.BytesString(bts)
}

func unmarshalFromString(str string, v interface{}) error {
	return encoding.GetJsonCodec().Unmarshal(unsafex.StringBytes(str), v)
}
