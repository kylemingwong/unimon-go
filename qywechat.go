package main

import (
	//	"bufio"
	//	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//	"os"
)

const (
	//send message url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//get token url
	get_token = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

var requestError = errors.New("request error,check url or network")

//access_token struct
type access_token struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

//config file struct
type AppConfig struct {
	Corpid  string
	Secret  string
	Agentid string
}

func main() {
	JsonParse := NewJsonStruct()

	v := AppConfig{}

	JsonParse.Load("./config.json", &v)

	fmt.Println(v.Corpid)
	fmt.Println(v.Secret)
	fmt.Println(v.Agentid)
	my_token, err := GetToken(v.Corpid, v.Secret)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(my_token.Access_token)
	fmt.Println(my_token.Expires_in)

}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

//load json file!
func (jst *JsonStruct) Load(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

//get wechat token
func GetToken(corpid, corpsecret string) (at access_token, err error) {
	resp, err := http.Get(get_token + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if at.Access_token == "" {
		err = errors.New("corpid or corpsecret error.")
	}
	return
}
