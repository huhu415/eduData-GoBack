package identimage

import (
	"eduData/bootstrap"
	"fmt"
	"testing"
)

var base64image string = "R0lGODlhNAAUAPcAAAAAAAAAMwAAZgAAmQAAzAAA/wArAAArMwArZgArmQArzAAr/wBVAABVMwBVZgBVmQBVzABV/wCAAACAMwCAZgCAmQCAzACA/wCqAACqMwCqZgCqmQCqzACq/wDVAADVMwDVZgDVmQDVzADV/wD/AAD/MwD/ZgD/mQD/zAD//zMAADMAMzMAZjMAmTMAzDMA/zMrADMrMzMrZjMrmTMrzDMr/zNVADNVMzNVZjNVmTNVzDNV/zOAADOAMzOAZjOAmTOAzDOA/zOqADOqMzOqZjOqmTOqzDOq/zPVADPVMzPVZjPVmTPVzDPV/zP/ADP/MzP/ZjP/mTP/zDP//2YAAGYAM2YAZmYAmWYAzGYA/2YrAGYrM2YrZmYrmWYrzGYr/2ZVAGZVM2ZVZmZVmWZVzGZV/2aAAGaAM2aAZmaAmWaAzGaA/2aqAGaqM2aqZmaqmWaqzGaq/2bVAGbVM2bVZmbVmWbVzGbV/2b/AGb/M2b/Zmb/mWb/zGb//5kAAJkAM5kAZpkAmZkAzJkA/5krAJkrM5krZpkrmZkrzJkr/5lVAJlVM5lVZplVmZlVzJlV/5mAAJmAM5mAZpmAmZmAzJmA/5mqAJmqM5mqZpmqmZmqzJmq/5nVAJnVM5nVZpnVmZnVzJnV/5n/AJn/M5n/Zpn/mZn/zJn//8wAAMwAM8wAZswAmcwAzMwA/8wrAMwrM8wrZswrmcwrzMwr/8xVAMxVM8xVZsxVmcxVzMxV/8yAAMyAM8yAZsyAmcyAzMyA/8yqAMyqM8yqZsyqmcyqzMyq/8zVAMzVM8zVZszVmczVzMzV/8z/AMz/M8z/Zsz/mcz/zMz///8AAP8AM/8AZv8Amf8AzP8A//8rAP8rM/8rZv8rmf8rzP8r//9VAP9VM/9VZv9Vmf9VzP9V//+AAP+AM/+AZv+Amf+AzP+A//+qAP+qM/+qZv+qmf+qzP+q///VAP/VM//VZv/Vmf/VzP/V////AP//M///Zv//mf//zP///wAAAAAAAAAAAAAAACH5BAEAAPwALAAAAAA0ABQAAAj/APcJHEiwoMGDCBMqXMiwocOHECNKnEjxoL5JkkIhvJjRoLBJaCbBOTgsZEiNC3mdRKhyEkqB+iSlkSlzUrSB89DQnJlGk0JhaUAuOwhUKMFhklzumzdzpMCPaYbuA9rRoD5GM1cSvJpVKUysTvfxCnqTnkxKN5eC3GQwJ8hJM18uNQm3qlo0PgUiRcP2IL2QYQeajepW7uBlhXE2xQkY4Uc0cgXOczrMKE7KlvfFnKRm4OO8Bf9mPpiYtFaBYzsylQmaMcjWB9vFRSjbruakdBsTHCxTasLSBpP1Lsix52q5bkOmVVgZje+CtSMXbP5y3ts0yxUCLwico5q0UNPGOOT5XGF02rMHlpTE1jpezzyzLxQeFSH95+N3akpLD+PMMTLphp5tBJ23VS8m9SXZTjS9VtGDCgUEADs="
var test string = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAAZAFADASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD1bxlri+GfCWqawzIHtYGaLzFZlaU/LGpC84LlR269R1rynw78V9X1m0aW/utFtR5UszNa2F1dPbLEVDedErZCsJFYSBto2MDyfl9d8UaPd6zpJt9P1WfS71JUlhu4stsZWBIZAwDqRkFWyvOcHArz7RfhE3h3Qr3S7LxA4XVFEWpytaAtJGC3yw/NiMlHZSW8zqCNuKANvxH4ptfDPg/+3LqaG9/dJ5X2dgi3TsBjZkn5Ty3BYhQTzivO/AOuWdn4m8U3Ot69pwuLz7JMZGu0Ee5kZmjRi5DLGW2cE8KK9Ym8P6TJpNrpc+nW1xY2qokEFxGJVQKu1fvZ5A4z1rkU8Dppeu6vcaUkNrFq8aRGS3hWE2EartcRFTku5IIwAAV3HJUKwBvQ3Vtf2yXNncQ3Fu+dssLh1bBwcEcHkEVFIKuCCK3gSGCNIoo1CJGihVVQMAADoAKryCgChIKpyCr8oqpIKAM+QVUkFX5BVOQUAe2SCqkq1elqnJQBRkFU5BV6WqktAFGQVTkFXpKpyUAUpBVOVavS1TkoApSCqcgq9JVOWgD/2Q=="

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
