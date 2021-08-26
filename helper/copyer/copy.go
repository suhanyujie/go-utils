package copyer

import "github.com/suhanyujie/go-utils/libs/jsonx"

func Copy(src interface{}, dst interface{}) error {
	json1 := jsonx.ToJsonIgnoreErr(src)
	err := jsonx.FromJson(json1, &dst)
	if err != nil {
		return err
	}
	return nil
}
