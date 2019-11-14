package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Gvinaxu/cli/util"
	"github.com/valyala/fasthttp"
)

var (
	token = &AccessToken{}
)

const (
	loginUrl        = "http://127.0.0.1:18080/api/1/user/login"
	refreshTokenUrl = "http://127.0.0.1:18080/api/1/user/refreshToken"
)

type Account struct {
	name     string
	password string
}

func NewAccount(name, password string) *Account {
	return &Account{name: name, password: password}
}

func (a *Account) Login() (*AccessToken, error) {
	args := &fasthttp.Args{}
	args.Add("name", a.name)
	args.Add("pass", a.password)

	status, resp, err := util.DoTimeout(args, "POST", loginUrl, nil)
	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response code is %d", status))
	}
	type R struct {
		Code int          `json:"code"`
		Data *AccessToken `json:"data"`
	}
	r := &R{}
	err = json.Unmarshal(resp, r)
	if err != nil {
		return nil, err
	}
	token = r.Data
	return r.Data, nil
}

func (a *Account) RefreshToken() (*AccessToken, error) {

	head := map[string]interface{}{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	status, resp, err := util.DoTimeout(nil, "POST", refreshTokenUrl, head)
	if status != http.StatusOK {
		fmt.Println(status)
		return nil, errors.New(fmt.Sprintf("response code is %d", status))
	}
	type R struct {
		Code int          `json:"code"`
		Data *AccessToken `json:"data"`
	}
	r := &R{}
	err = json.Unmarshal(resp, r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	token = r.Data
	return r.Data, nil
}

func (a *Account) GetToken() *AccessToken {
	return token
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int64  `json:"expires"`
}
