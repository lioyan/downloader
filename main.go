package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)


func main() {
	url := "https://test-xjy-file.obs.cn-east-2.myhuaweicloud.com/202006%2Fb917e15c-e9ef-4e8c-89ed-b2e4f9047ec9.pdf?response-content-disposition=attachment%3Bfilename%3DPDF%25E6%25B5%258B%25E8%25AF%2595%25E6%2596%2587%25E4%25BB%25B6.pdf%3Bfilename%2A%3Dutf-8%27%27PDF%25E6%25B5%258B%25E8%25AF%2595%25E6%2596%2587%25E4%25BB%25B6.pdf&response-content-type=binary%2Foctet-stream&AWSAccessKeyId=R7WWGCCN8NPIWHPLFOLT&Expires=1653795386&Signature=rU94FizzTCzB9JneBpLy%2F8%2BAF6I%3D"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body);err == nil {
		ioutil.WriteFile("测试文件.pdf", data, 0644)
	}
}
