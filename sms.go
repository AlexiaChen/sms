package sms

import (
	"encoding/json"
	"fmt"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func createSmsClient(smsEndpoint, protocol, accessKeyId, accessKeySecret string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(accessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(accessKeySecret),
		// Http
		Protocol: tea.String(protocol),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String(smsEndpoint)
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func sendSms(client *dysmsapi20170525.Client, phoneNumber string, signName string, templateCode string, content string) (result string, _err error) {

	params := map[string]string{
		"content": content,
	}

	templateParam, _ := json.Marshal(params)

	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(string(templateParam)),
	}

	res, tryErr := func() (_res *dysmsapi20170525.SendSmsResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_res, _err := client.SendSmsWithOptions(request, &util.RuntimeOptions{})
		if _err != nil {
			return nil, fmt.Errorf("SendSmsWithOptions error: %s", _err.Error())
		}

		return _res, nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		//logs.GetLogger().Error(tea.StringValue(error.Message), tryErr)
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend := m["Recommend"]
			fmt.Println(recommend)
			//logs.GetLogger().Error(fmt.Sprintf("Recommend: %+v", recommend), tryErr)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return "", _err
		}

		return "", tryErr
	}

	return res.String(), nil
}
