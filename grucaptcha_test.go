package grucaptcha

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestRuCaptcha_ResolveReCaptchaV2(t *testing.T) {
	body, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Fatalln(err, "not found token.txt example file")
	}
	response, err := func(ruCaptchaApiKey string, timeout time.Duration) (string, error) {
		rucaptcha := NewRucaptcha(ruCaptchaApiKey)
		captchaParams := ReCaptchaV2Params{
			GoogleKey: "6LdGiQ0TAAAAAJ8Gm_YHJXAbv1lTeJV9Ez-aT_Jp",
			PageUrl:   "https://passport.ngs.ru/register/",
		}
		start := time.Now()
		resultChan, err := rucaptcha.ResolveReCaptchaV2(captchaParams)
		if err != nil {
			return "", nil
		}
		select {
		case results := <-resultChan:
			if results.Result != "" {
				return results.Result, nil
			} else {
				err = fmt.Errorf("captcha result is empty, %+v", results)
				return "", nil
			}
		case <-time.After(timeout):
			err = fmt.Errorf("timeout elapsed request is too long %s, max %s", time.Now().Sub(start), timeout)
			return "", nil
		}
	}(string(body), time.Second*60)
	fmt.Println(response, err)
}
