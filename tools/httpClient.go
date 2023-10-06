package tools

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpGet(url, getInfo string) (string, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	//resp, _ := http.Get("http://localhost:5000/login?username=admin&password=111111")
	resp, err1 := client.Get(url + "?" + getInfo)
	if err1 != nil {
		Log.Error(fmt.Sprintf("httpGetError:%v,%v", url+"?"+getInfo, err1.Error()))
		return "", err1
	}
	defer resp.Body.Close()
	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		Log.Error(fmt.Sprintf("httpGetError:%v,%v", url+"?"+getInfo, err2.Error()))
		return "", err2
	}
	return string(body), nil

}

func HttpGetByUrl(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err1 := client.Get(url)
	if err1 != nil {
		Log.Error(fmt.Sprintf("httpGetError:%v,%v", url, err1.Error()))
		return "", err1
	}
	defer resp.Body.Close()
	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		Log.Error(fmt.Sprintf("httpGetError:%v,%v", url, err2.Error()))
		return "", err2
	}
	return string(body), nil
}

func HttpPost(url, contentType, getInfo string) (string, error) {
	client := &http.Client{}
	resp, err1 := client.Post(url, contentType, strings.NewReader(getInfo))
	if err1 != nil {
		Log.Error(fmt.Sprintf("HttpPostError:%v,%v", url+"?"+getInfo, err1.Error()))
		return "", err1
	}
	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		Log.Error(fmt.Sprintf("HttpPostError:%v,%v", url+"?"+getInfo, err2.Error()))
		return "", err2
	}
	resp.Body.Close()
	return string(body), nil
}

func HttpPostRaw(url, data string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("err", err.Error()))
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("err", err.Error()))
		return nil
	}
	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil
	}
	err = resp.Body.Close()
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("err", err.Error()))
	}
	return body
}

func HttpPostJson(url, data string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("err", err.Error()))
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("err", err.Error()), zap.Any("req", data), zap.Any("res", response))
		return nil
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			Log.Error(RunFuncName(), zap.Any("err", err.Error()))
		}
	}()
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("err", err.Error()), zap.Any("req", data), zap.Any("res", response))
		return nil
	}
	return bodyData
}
