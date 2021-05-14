package grucaptcha

import (
	"errors"
	"fmt"
)

type Geetest struct {
	GoRuCaptcha
}

type GeetestParams struct {
	PageUrl   string
	Gt        string
	Challenge string
	ApiServer string
}

type GeetestResultResult struct {
	GeetestChallenge string `json:"geetest_challenge"`
	GeetestValidate  string `json:"geetest_validate"`
	GeetestSeccode   string `json:"geetest_seccode"`
}

type GeetestResult struct {
	JobId  string              //Job ID
	Result GeetestResultResult //Result string [INFO] {"geetest_challenge":"1b833024591a94fe750a7946464a54a5he","geetest_validate":"646539128f52b2bb0556d7d30bce52a2","geetest_seccode":"646539128f52b2bb0556d7d30bce52a2|jordan"}}
	Error  error               //Error message
}

func (h *Geetest) Resolve(params GeetestParams) (chan GeetestResult, error) {
	requestParams := map[string]string{}
	if params.PageUrl == "" {
		return nil, errors.New("PageUrl is empty")
	}
	if params.Gt == "" {
		return nil, errors.New("gt is empty")
	}
	if params.Challenge == "" {
		return nil, errors.New("Challenge is empty")
	}
	requestParams["pageurl"] = params.PageUrl
	requestParams["gt"] = params.Gt
	requestParams["challenge"] = params.Challenge
	if params.ApiServer != "" {
		requestParams["api_server"] = params.ApiServer
	}

	captchaParams := CaptchaParams{
		Method: "geetest",
		Params: requestParams,
	}

	resChan := make(chan GeetestResult, 1)
	ch, err := h.resolveCaptcha(captchaParams)
	if err == nil {
		go func() {
			for msg := range ch {
				fmt.Printf("[INFO] Incoming msg %+v\n", msg)
				var result GeetestResult
				if msg.Result == nil {
					result = GeetestResult{
						JobId:  msg.JobId,
						Result: GeetestResultResult{},
						Error:  msg.Error,
					}
				} else {
					resultResult := msg.Result.(map[string]interface{})
					result = GeetestResult{
						JobId: msg.JobId,
						Result: GeetestResultResult{
							GeetestChallenge: resultResult["geetest_challenge"].(string),
							GeetestSeccode:   resultResult["geetest_seccode"].(string),
							GeetestValidate:  resultResult["geetest_validate"].(string),
						},
						Error: msg.Error,
					}
				}
				resChan <- result
			}
		}()
	}
	return resChan, err
}

func NewGeetest(key string) Geetest {
	geetest := Geetest{}
	geetest.key = key
	return geetest
}
