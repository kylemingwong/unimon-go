package main

import (
	"bytes"
	"encoding/json"
	"errors"
	//"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/larspensjo/config"
)

const (
	//send message url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//get token url
	get_token = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

var requestError = errors.New("request error,check url or network")

// message struct
type send_msg struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

//access_token struct
type access_token struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

//send messae error type
type send_msg_error struct {
	Errcode int    `json:"errcode`
	Errmsg  string `json:"errmsg"`
}

//config
func main() {

	corpid, secret, agentid := loadConfig("config.ini")
	fmt.Println(corpid)
	fmt.Println(secret)
	fmt.Println(agentid)

	msg := "this is a test message send by golang app!"

	my_token, err := getToken(corpid, secret)
	if err != nil {
		panic(err)
	}
	fmt.Println(my_token.Access_token)
	var m send_msg = send_msg{Touser: "@all", Msgtype: "text", Agentid: agentid, Text: map[string]string{"content": msg}}

	bufmsg, senderr := json.Marshal(m)
	if senderr != nil {
		fmt.Println(senderr)
		return
	}
	err = sendMsg(my_token.Access_token, bufmsg)
	if err != nil {
		fmt.Println(err)
	}

}

//load configfile return corpid secret agentid
func loadConfig(configFile string) (corpid string, secret string, agentid int) {
	c, err := config.ReadDefault("config.ini")
	if err != nil {
		panic(err)
	}
	if c.HasSection("qywechat") {
		corpid, _ = c.String("qywechat", "corpid")
		secret, _ = c.String("qywechat", "secret")
		agentid, _ = c.Int("qywechat", "agentid")

		return
	}
	return
}

//get wechat token
func getToken(corpid string, secret string) (at access_token, err error) {

	resp, err := http.Get(get_token + corpid + "&corpsecret=" + secret)
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

//send message to wechat
func sendMsg(Access_token string, msgbody []byte) error {
	body := bytes.NewBuffer(msgbody)
	resp, err := http.Post(sendurl+Access_token, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e send_msg_error
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}

////db option
//1 get alarm config value
//2 get record match alarm config values from db and return a alarm message!
//
