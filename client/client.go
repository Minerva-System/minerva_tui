package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

type NewUserRequest struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DefaultResponse struct {
	Message string `json:"message"`
}

type MinervaClient struct {
	url    string
	tenant string
	token  string
	client *http.Client
}

func Create(url string, tenant string) (MinervaClient, error) {
	return MinervaClient{
		url:    url,
		tenant: tenant,
		token:  "",
		client: &http.Client{
			Timeout: 10 * time.Second,
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

	// Save token
	c.token = res.Token

	return response.StatusCode, res, ""
}

func (c *MinervaClient) Logout() (int, DefaultResponse, string) {
	var res DefaultResponse

	url := c.Url("/" + c.Tenant() + "/logout")

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return 0, DefaultResponse{}, fmt.Sprintf("Erro: %v", err)
	}

	req.Header.Add("Authorization", "Bearer " + c.token)
	response, err := c.client.Do(req)
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

	url := c.Url(fmt.Sprintf("/" + c.Tenant() + "/user?page=%d", page))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}
	
	req.Header.Add("Authorization", "Bearer " + c.token)
	response, err := c.client.Do(req)
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

func (c *MinervaClient) UserCreate(req_data NewUserRequest) (int, UserInfo, string) {
	res := UserInfo{}

	url := c.Url("/" + c.Tenant() + "/user")
	payload, err := json.Marshal(req_data)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}
	req.Header.Add("Authorization", "Bearer " + c.token)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	response, err := c.client.Do(req)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	defer response.Body.Close()
	if response.StatusCode > 299 {
		errRes := DefaultResponse{}
		err = json.NewDecoder(response.Body).Decode(&errRes)
		if err != nil {
			return 0, res, fmt.Sprintf("Erro: %v", err)
		}
		return response.StatusCode, res, fmt.Sprintf("Erro: %s", errRes.Message)
	}

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	return response.StatusCode, res, ""
}


func (c *MinervaClient) UserRemove(index int64) (int, string) {
	url := c.Url(fmt.Sprintf("/" + c.Tenant() + "/user/%d", index))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, fmt.Sprintf("Erro: %v", err)
	}

	req.Header.Add("Authorization", "Bearer " + c.token)
	response, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Sprintf("Erro: %v", err)
	}
	defer response.Body.Close()
	return response.StatusCode, ""
}

func (c *MinervaClient) UserUpdate(index int64, data NewUserRequest) (int, UserInfo, string) {
	res := UserInfo{}
	url := c.Url(fmt.Sprintf("/" + c.Tenant() + "/user/%d", index))
	payload, err := json.Marshal(data)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	req.Header.Add("Authorization", "Bearer " + c.token)
	response, err := c.client.Do(req)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		errRes := DefaultResponse{}
		err = json.NewDecoder(response.Body).Decode(&errRes)
		if err != nil {
			return 0, res, fmt.Sprintf("Erro: %v", err)
		}
		return response.StatusCode, res, fmt.Sprintf("Erro: %s", errRes.Message)
	}

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return 0, res, fmt.Sprintf("Erro: %v", err)
	}

	return response.StatusCode, res, ""
}
