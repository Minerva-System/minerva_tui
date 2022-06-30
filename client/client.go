package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Tenant string `json:"tenant"`
}

type UserInfo struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DefaultResponse struct {
	Message string `json:"message"`
}

type MinervaClient struct {
	url    string
	tenant string
	client *http.Client
}

func Create(url string, tenant string) (MinervaClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return MinervaClient{}, err
	}

	return MinervaClient{
		url:    url,
		tenant: tenant,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
	}, nil
}

func (c MinervaClient) Url(endpoint string) string {
	return c.url + endpoint
}

func (c MinervaClient) Tenant() string {
	return c.tenant
}

func (c *MinervaClient) Login(req LoginRequest) (int, LoginResponse, string) {
	var res LoginResponse

	url := c.Url("/" + c.Tenant() + "/login")
	payload, err := json.Marshal(req)
	if err != nil {
		return 0, LoginResponse{}, fmt.Sprintf("Erro: %v", err)
	}

	buffer := bytes.NewBuffer(payload)
	response, err := c.client.Post(url, "application/json; charset=utf-8", buffer)
	if err != nil {
		return 0, LoginResponse{}, fmt.Sprintf("Erro: %v", err)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	return response.StatusCode, res, ""
}

func (c *MinervaClient) Logout() (int, DefaultResponse, string) {
	var res DefaultResponse

	url := c.Url("/logout")
	response, err := c.client.Post(url, "", nil)
	if err != nil {
		return 0, DefaultResponse{}, fmt.Sprintf("Erro: %v", err)
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return 0, DefaultResponse{}, fmt.Sprintf("Erro: %v", err)
	}

	return response.StatusCode, res, ""
}

func (c *MinervaClient) UserList(page int) (int, []UserInfo, string) {
	res := make([]UserInfo, 0)

	url := c.Url(fmt.Sprintf("/user?page=%d", page))
	response, err := c.client.Get(url)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}
	
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	return response.StatusCode, res, ""
}
