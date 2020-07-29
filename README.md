# Go library for Rucaptcha

Golang module for auto resolve google recaptcha using [`rucaptcha`][1] API.

### Resolve hcaptcha

It is a few steps for resolve captcha:

1. Find `data-sitekey` attribute in HTML page
2. Wait resolving the captcha
3. Put value to `h-captcha-response` and `g-recaptcha-response`
4. Submit hcaptcha form
5. Use same Cookie and User-Agent headers for all future requests

Check example directory for full example

```go
hcaptcha := grucaptcha.NewHcaptcha(api.RuCaptchaKey)
params := grucaptcha.HCaptchaParams{
    SiteKey: siteKey,
    PageUrl: uri,
}
ch, err := hcaptcha.Resolve(params)
if err != nil {
    log.Fatalln(time.Now().Format(time.StampMicro), "[ERR] Rucaptcha response", err)
}
var res grucaptcha.RuCaptchaResult
select {
case res = <-ch:
case <-time.After(time.Minute):
    log.Fatalln(time.Now().Format(time.StampMicro), "[ERR] Captcha resolve timeout")
}

if res.Result == "" {
    log.Fatalln(time.Now().Format(time.StampMicro), "[ERR] Rucaptcha empty result", res)
}
```

### Resolve google recaptcha

It is very similar to hcaptcha. 

1. Find `data-sitekey` attribute in HTML page
2. Wait resolving the captcha
3. Put value to `g-recaptcha-response`
4. Submit recaptcha form
5. Use same Cookie and User-Agent headers for all future requests

Check example folder for more detailed example.

```go
rucaptcha := grucaptcha.NewRecaptcha("<your_api_key>")
captchaParams := grucaptcha.ReCaptchaV2Params{
    GoogleKey: "<data-sitekey attribute>",
    PageUrl:   "https://full-url.com/auth",
}
resultChan, err := rucaptcha.Resolve(captchaParams)
if err != nil {
    panic(err)
}
results := <-resultChan

//Example for selenium
//Insert resolved captcha string into the hidden textarea. And submit the form
script := fmt.Sprintf(
    `document.getElementById("g-recaptcha-response").value="%s"`, 
    results.Result)
wd.ExecuteScript(script, nil)
```

[1]: https://rucaptcha.com?from=7563013