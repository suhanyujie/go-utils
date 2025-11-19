package logx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/suhanyujie/go_utils/core/consts"
	"github.com/suhanyujie/go_utils/core/model/vo"
	"github.com/suhanyujie/go_utils/jsonx"
)

var (
	WebHookUrl = ""
)

func SetWebHookUrl(url string) {
	WebHookUrl = url
}

func Send2FsWithWarningMsg(env string, content string) (*vo.BaseSucResp, error) {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(fmt.Sprintf("**env:**%v\n", env))
	strBuilder.WriteString(fmt.Sprintf("**content:**%v", content))
	return Send2FsGroup(strBuilder.String())
}

// Send2FsGroup 发送告警消息到飞书群
// ref: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN?lang=zh-CN
func Send2FsGroup(msg string) (*vo.BaseSucResp, error) {
	bodyObj := vo.Send2FsGroupReq{
		MsgType: "interactive",
		Card: vo.Send2FsGroupReqCard{
			Header: vo.Send2FsGroupReqCardHeader{
				Title: vo.Send2FsGroupReqCardHeaderTitle{
					Tag:     "plain_text",
					Content: "[服务告警]",
				},
			},
			Elements: []vo.Send2FsGroupReqElement{
				{
					Tag: "div",
					Text: vo.Send2FsGroupReqElementText{
						Tag:     "lark_md",
						Content: msg,
					},
				},
			},
		},
	}
	if WebHookUrl == "" {
		Logger.Infof("[Send2FsGroup] please set the webhook url")
		return nil, errors.New("no web hook url")
	}
	resp, code, err := SimplePost(WebHookUrl, jsonx.ToJsonIgnoreErr(bodyObj))
	if err != nil {
		// 这里使用 Infof 是为了防止死循环。error 级别会触发发送日志。
		Logger.Infof("err: %v", err)
		return nil, err
	}
	if code != consts.PlatformHttpApiOk {
		Logger.Infof("http code err code: %v", code)
		return nil, fmt.Errorf("http code err code: %v", code)
	}
	respVo := new(vo.BaseSucResp)
	jsonx.FromJson(resp, &respVo)
	return respVo, nil
}
