#Google recaptcha resolve

Golang module for auto resolve google recaptcha, using [`rucaptcha`][1] API.

###Usage
```go
rucaptcha := grucaptcha.NewRucaptcha("<your_api_key>")
captchaParams := grucaptcha.ReCaptchaV2Params{
    GoogleKey: "<data-sitekey attribute>",
    PageUrl:   "https://full-url.com/auth",
}
resultChan, err := rucaptcha.ResolveReCaptchaV2(captchaParams)
if err != nil {
    panic(err)
}
results := <-resultChan

//Example for selenium
//Insert resolved captcha string into the hidden textarea. And submit form
script := fmt.Sprintf(
    `document.getElementById("g-recaptcha-response").value="%s"`, 
    results.Result)
wd.ExecuteScript(script, nil)

```

[1]: https://rucaptcha.com?from=7563013