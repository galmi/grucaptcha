package grucaptcha

import (
	"errors"
)

type HCaptcha struct {
	GoRuCaptcha
}

type HCaptchaParams struct {
	SiteKey string
	PageUrl string
}

type HcaptchaResult struct {
	JobId  string //Job ID
	Result string //Result string
	Error  error  //Error message
}

func (h *HCaptcha) Resolve(params HCaptchaParams) (chan HcaptchaResult, error) {
	requestParams := map[string]string{}
	if params.SiteKey == "" {
		return nil, errors.New("Sitekey is empty")
	}
	if params.PageUrl == "" {
		return nil, errors.New("PageUrl is empty")
	}
	requestParams["sitekey"] = params.SiteKey
	requestParams["pageurl"] = params.PageUrl

	captchaParams := CaptchaParams{
		Method: "hcaptcha",
		Params: requestParams,
	}

	resChan := make(chan HcaptchaResult, 1)

	ch, err := h.resolveCaptcha(captchaParams)
	if err == nil {
		go func() {
			for msg := range ch {
				result := HcaptchaResult{
					JobId:  msg.JobId,
					Result: msg.Result.(string),
					Error:  msg.Error,
				}
				resChan <- result
			}
		}()
	}
	return resChan, err
}

func NewHcaptcha(key string) HCaptcha {
	hcaptcha := HCaptcha{}
	hcaptcha.key = key
	return hcaptcha
}
