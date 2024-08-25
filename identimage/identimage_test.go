package identimage

import (
	"eduData/bootstrap"
	"fmt"
	"testing"
)

var base64image string = "R0lGODlhNAAUAPcAAAAAAAAAMwAAZgAAmQAAzAAA/wArAAArMwArZgArmQArzAAr/wBVAABVMwBVZgBVmQBVzABV/wCAAACAMwCAZgCAmQCAzACA/wCqAACqMwCqZgCqmQCqzACq/wDVAADVMwDVZgDVmQDVzADV/wD/AAD/MwD/ZgD/mQD/zAD//zMAADMAMzMAZjMAmTMAzDMA/zMrADMrMzMrZjMrmTMrzDMr/zNVADNVMzNVZjNVmTNVzDNV/zOAADOAMzOAZjOAmTOAzDOA/zOqADOqMzOqZjOqmTOqzDOq/zPVADPVMzPVZjPVmTPVzDPV/zP/ADP/MzP/ZjP/mTP/zDP//2YAAGYAM2YAZmYAmWYAzGYA/2YrAGYrM2YrZmYrmWYrzGYr/2ZVAGZVM2ZVZmZVmWZVzGZV/2aAAGaAM2aAZmaAmWaAzGaA/2aqAGaqM2aqZmaqmWaqzGaq/2bVAGbVM2bVZmbVmWbVzGbV/2b/AGb/M2b/Zmb/mWb/zGb//5kAAJkAM5kAZpkAmZkAzJkA/5krAJkrM5krZpkrmZkrzJkr/5lVAJlVM5lVZplVmZlVzJlV/5mAAJmAM5mAZpmAmZmAzJmA/5mqAJmqM5mqZpmqmZmqzJmq/5nVAJnVM5nVZpnVmZnVzJnV/5n/AJn/M5n/Zpn/mZn/zJn//8wAAMwAM8wAZswAmcwAzMwA/8wrAMwrM8wrZswrmcwrzMwr/8xVAMxVM8xVZsxVmcxVzMxV/8yAAMyAM8yAZsyAmcyAzMyA/8yqAMyqM8yqZsyqmcyqzMyq/8zVAMzVM8zVZszVmczVzMzV/8z/AMz/M8z/Zsz/mcz/zMz///8AAP8AM/8AZv8Amf8AzP8A//8rAP8rM/8rZv8rmf8rzP8r//9VAP9VM/9VZv9Vmf9VzP9V//+AAP+AM/+AZv+Amf+AzP+A//+qAP+qM/+qZv+qmf+qzP+q///VAP/VM//VZv/Vmf/VzP/V////AP//M///Zv//mf//zP///wAAAAAAAAAAAAAAACH5BAEAAPwALAAAAAA0ABQAAAj/APcJHEiwoMGDCBMqXMiwocOHECNKnEjxoL5JkkIhvJjRoLBJaCbBOTgsZEiNC3mdRKhyEkqB+iSlkSlzUrSB89DQnJlGk0JhaUAuOwhUKMFhklzumzdzpMCPaYbuA9rRoD5GM1cSvJpVKUysTvfxCnqTnkxKN5eC3GQwJ8hJM18uNQm3qlo0PgUiRcP2IL2QYQeajepW7uBlhXE2xQkY4Uc0cgXOczrMKE7KlvfFnKRm4OO8Bf9mPpiYtFaBYzsylQmaMcjWB9vFRSjbruakdBsTHCxTasLSBpP1Lsix52q5bkOmVVgZje+CtSMXbP5y3ts0yxUCLwico5q0UNPGOOT5XGF02rMHlpTE1jpezzyzLxQeFSH95+N3akpLD+PMMTLphp5tBJ23VS8m9SXZTjS9VtGDCgUEADs="

// TestNumberIdentify 百度的
func TestNumberIdentify(t *testing.T) {
	bootstrap.Loadconfig()
	ocr := NewBaiduOcr(bootstrap.C.BaiduRequestUrl, bootstrap.C.BaiduAccesstoken)
	verify, err := ocr.Identify(&base64image)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	if verify != "1192" {
		t.Errorf("expected 1192, got %s", verify)
	}
}

// TestCommonVerify 云码的
func TestCommonVerify(t *testing.T) {
	bootstrap.Loadconfig()
	ocr := NewJfbymOcr(bootstrap.C.JfymRequestUrl, bootstrap.C.JfymToken)
	verify, err := ocr.Identify(&base64image)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	if verify != "1192" {
		t.Errorf("expected 1192, got %s", verify)
	}
}

func TestYesCaptcha(t *testing.T) {
	bootstrap.Loadconfig()
	ocr := NewYescaptchaOcr(bootstrap.C.YescaptchaRequestUrl, bootstrap.C.YesCaptchaToken)
	verify, err := ocr.Identify(&base64image)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	if verify != "1192" {
		t.Errorf("expected 1192, got %s", verify)
	}
}
