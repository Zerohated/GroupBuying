package constant

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

const (
	// AppId       string = "gxbf44d41f1aed601cc"
	// AppSecurity string = "b5ed40fceb2944349224716ca5203ea6"
	Format = "2006-01-02"
	// Format = "2006-01-02T15:04:05"
)

// PostRequest receive content and url, return body,header,err
func PostRequest(content string, url string) (*http.Response, error) {
	req := bytes.NewBuffer([]byte(content))
	bodyType := "application/json;charset=utf-8"
	resp, err := http.Post(url, bodyType, req)
	return resp, err
}

// ConvertMD5 jus convert input string to MD5 encrypted result
func ConvertMD5(input string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(input))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
