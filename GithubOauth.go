package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Conf struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

var conf = Conf{
	ClientId:     "",
	ClientSecret: "",
	RedirectUrl:  "",
}

type GToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"` // 这个字段没用到
	Scope       string `json:"scope"`      // 这个字段也没用到
}

// 通过code获取token认证url
func GetTokenAuthUrl(code string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		conf.ClientId, conf.ClientSecret, code, conf.RedirectUrl,
	)
}
func GetUserInfo(token *GToken) (map[string]interface{}, error) {

	// 形成请求
	var userInfoUrl = "https://api.github.com/user" // github用户信息获取接口
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func GetToken(url string) (*GToken, error) {

	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}

	// 将响应体解析为 token，并返回
	var token GToken
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func Oauth(ctx *gin.Context) {

	var err error
	// 获取 code
	var code = ctx.Query("code")

	// 通过 code, 获取 token
	var tokenAuthUrl = GetTokenAuthUrl(code)
	var token *GToken
	if token, err = GetToken(tokenAuthUrl); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tokenAuthUrl)
	// 通过token，获取用户信息
	var userInfo map[string]interface{}
	if userInfo, err = GetUserInfo(token); err != nil {
		fmt.Println("获取用户信息失败，错误信息为:", err)
		return
	}
	fmt.Println(userInfo)
	if err != nil {
		return
	}
}
