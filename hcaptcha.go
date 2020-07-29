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

func (h *HCaptcha) Resolve(params HCaptchaParams) (chan RuCaptchaResult, error) {
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

	return h.resolveCaptcha(captchaParams)
}

func NewHcaptcha(key string) HCaptcha {
	hcaptcha := HCaptcha{}
	hcaptcha.key = key
	return hcaptcha
}
