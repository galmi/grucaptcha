package example

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"grucaptcha"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const RuCaptchaKey = "<insert your api key>"
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:73.0) Gecko/20100101 Firefox/73.0"

var jar, _ = cookiejar.New(nil)

func getHttpClient() *http.Client {
	client := http.DefaultClient
	client.Jar = jar
	return client
}

func ResolveCaptcha(uri string) {
	//IMPORTANT!!! Use same httpClient, cookiejar variable and UserAgent header for all future requests
	httpClient := getHttpClient()
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Add("User-Agent", UserAgent)
	resp, _ := httpClient.Do(req)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Resolve captcha if sitekey exists in HTML
	nodes := doc.Find("script[data-sitekey]")
	if len(nodes.Nodes) == 0 {
		return
	}
	siteKey, _ := nodes.Attr("data-sitekey")

	//Start resolve hcaptcha
	fmt.Println(time.Now().Format(time.StampMicro), "[INFO] Start resolve captcha")
	hcaptcha := grucaptcha.NewHcaptcha(RuCaptchaKey)
	params := grucaptcha.HCaptchaParams{
		SiteKey: siteKey,
		PageUrl: uri,
	}
	ch, err := hcaptcha.Resolve(params)
	if err != nil {
		log.Fatalln(time.Now().Format(time.StampMicro), "[ERR] Rucaptcha response", err)
	}
	var res grucaptcha.RuCaptchaResult
	//Waiting hcaptcha resolving with 1 minute timeout
	select {
	case res = <-ch:
	case <-time.After(time.Minute):
		log.Fatalln(time.Now().Format(time.StampMicro), "[ERR] Captcha resolve timeout")
	}
	if res.Result == "" {
		log.Fatalln(time.Now().Format(time.StampMicro), "[ERR] Rucaptcha empty result", res)
	}

	//Hcaptcha was resolved, prepare and submit challenge form
	form := doc.Find("#challenge-form")
	formAction, _ := form.Attr("action")
	formUrl := uri + formAction
	r, _ := doc.Find("input[name=r]").Attr("value")
	captchaKind, _ := doc.Find("input[name=cf_captcha_kind]").Attr("value")
	vc, _ := doc.Find("input[name=vc]").Attr("value")

	formData := url.Values{
		"r":               {r},
		"cf_captcha_kind": {captchaKind},
		"vc":              {vc},
	}
	formData.Set("h-captcha-response", res.Result)
	formData.Set("g-recaptcha-response", res.Result)

	req, err = http.NewRequest("POST", formUrl, strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", UserAgent)
	resp, _ = httpClient.Do(req)
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//Checking response, if captcha was resolved successfully, sitekey will not exist
	nodes = doc.Find("script[data-sitekey]")
	if len(nodes.Nodes) > 0 {
		log.Fatal(time.Now().Format(time.StampMicro), "[INFO] Captcha not resolved")
	}
	fmt.Println(time.Now().Format(time.StampMicro), "[INFO] Captcha resolved")
}
