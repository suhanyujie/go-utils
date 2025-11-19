package jsonx

import (
	"github.com/goccy/go-json"
)

func ToJsonIgnoreErr(v interface{}) string {
	bArr, _ := json.Marshal(v)
	return string(bArr)
}

func ToJson(v interface{}) (string, error) {
	bArr, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bArr), nil
}

func FromJson(jsonStr string, o interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), &o)
	if err != nil {
		return err
	}
	return nil
}

func ToJsonFormatIgnoreErr(v interface{}) string {
	bArr, _ := json.MarshalIndent(v, "", "    ")
	return string(bArr)
}
