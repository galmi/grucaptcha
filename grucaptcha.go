package grucaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const SOFT_ID = "7563013"
const SEND_JOB_URL = "https://2captcha.com/in.php"
const CHECK_JOB_URL = "https://2captcha.com/res.php"

type RuCaptcha struct {
	key string
}

type ProxyType string

const HTTP ProxyType = "HTTP"
const HTTPS ProxyType = "HTTPS"
const SOCKS4 ProxyType = "SOCKS4"
const SOCKS5 ProxyType = "SOCKS5"

type ReCaptchaV2Params struct {
	Method    string    //userrecaptcha - defines that you're sending a ReCaptcha V2 with new method
	GoogleKey string    //Value of k or data-sitekey parameter you found on page
	PageUrl   string    //Full URL of the page where you see the ReCaptcha
	Invisible int       //1 - means that ReCaptcha is invisible. 0 - normal ReCaptcha. Default - 0
	Proxy     string    //Format: login:password@123.123.123.123:3128
	ProxyType ProxyType //Type of your proxy: HTTP, HTTPS, SOCKS4, SOCKS5.
}

type RuCaptchaResp struct {
	Status  int
	Request string
}

type RuCaptchaResult struct {
	JobId  string //Job ID
	Result string //Result string
	Error  error  //Error message
}

func (r *RuCaptcha) requestJob(params map[string]string) (string, error) {
	req, err := http.NewRequest("GET", SEND_JOB_URL, nil)
	if err != nil {
		return "", err
	}
	params["key"] = r.key
	params["soft_id"] = SOFT_ID
	params["json"] = "1"
	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	respData := RuCaptchaResp{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}

	if respData.Status == 0 {
		return "", errors.New(respData.Request)
	}

	return respData.Request, nil
}

func (r *RuCaptcha) checkJob(jobId string) (string, error) {
	req, err := http.NewRequest("GET", CHECK_JOB_URL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", r.key)
	q.Add("soft_id", SOFT_ID)
	q.Add("json", "1")
	q.Add("action", "get")
	q.Add("id", jobId)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	respData := RuCaptchaResp{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}

	if respData.Status == 0 {
		return "", errors.New(respData.Request)
	}

	return respData.Request, nil
}

func (r *RuCaptcha) ResolveReCaptchaV2(params ReCaptchaV2Params) (chan RuCaptchaResult, error) {
	respChan := make(chan RuCaptchaResult, 1)
	requestParams := map[string]string{}
	if params.GoogleKey == "" {
		return nil, errors.New("GoogleKey is empty")
	}
	if params.PageUrl == "" {
		return nil, errors.New("PageUrl is empty")
	}

	if params.Method == "" {
		params.Method = "userrecaptcha"
	}
	requestParams["method"] = params.Method
	requestParams["googlekey"] = params.GoogleKey
	requestParams["pageurl"] = params.PageUrl
	requestParams["invisible"] = strconv.Itoa(params.Invisible)
	if params.Proxy != "" && params.ProxyType != "" {
		requestParams["proxy"] = params.Proxy
		requestParams["proxytype"] = string(params.ProxyType)
	}

	jobId, err := r.requestJob(requestParams)
	if err != nil {
		return nil, err
	}

	go func(jobId string) {
		defer close(respChan)
		for {
			time.Sleep(time.Second * 5)
			jobResult, err := r.checkJob(jobId)
			if err != nil && err.Error() == "CAPCHA_NOT_READY" {
				continue
			}
			if err != nil {
				respChan <- RuCaptchaResult{
					JobId: jobId,
					Error: err,
				}
				break
			}
			respChan <- RuCaptchaResult{
				JobId:  jobId,
				Result: jobResult,
			}
			break
		}
	}(jobId)

	return respChan, nil
}

func NewRucaptcha(key string) RuCaptcha {
	return RuCaptcha{
		key: key,
	}
}
