package grucaptcha

import (
	"errors"
	"strconv"
)

type ReCaptcha struct {
	GoRuCaptcha
}

type ReCaptchaV2Params struct {
	GoogleKey string //Value of k or data-sitekey parameter you found on page
	PageUrl   string //Full URL of the page where you see the ReCaptcha
	Invisible int    //1 - means that ReCaptcha is invisible. 0 - normal ReCaptcha. Default - 0
}

func (r *ReCaptcha) Resolve(params ReCaptchaV2Params) (chan RuCaptchaResult, error) {
	requestParams := map[string]string{}
	if params.GoogleKey == "" {
		return nil, errors.New("GoogleKey is empty")
	}
	if params.PageUrl == "" {
		return nil, errors.New("PageUrl is empty")
	}
	requestParams["googlekey"] = params.GoogleKey
	requestParams["pageurl"] = params.PageUrl
	requestParams["invisible"] = strconv.Itoa(params.Invisible)

	captchaParams := CaptchaParams{
		Method: "userrecaptcha",
		Params: requestParams,
	}

	return r.resolveCaptcha(captchaParams)
}

func NewRecaptcha(key string) ReCaptcha {
	recaptcha := ReCaptcha{}
	recaptcha.key = key
	return recaptcha
}
