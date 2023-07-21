package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gwaycc/bchain-storage/lib/utils"
	"github.com/gwaylib/errors"
)

var (
	ErrAuthFailed = errors.New("auth failed")
)

type AuthClient struct {
	Host   string
	User   string
	Passwd string
}

func NewAuthClient(host, user, passwd string) *AuthClient {
	// TODO: check the host format
	return &AuthClient{Host: host, User: user, Passwd: passwd}
}

func (auth *AuthClient) Check(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+auth.Host+"/check", nil)
	if err != nil {
		return nil, errors.As(err)
	}
	req.SetBasicAuth(auth.User, auth.Passwd)
	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
	}
	return respBody, nil
}

func (auth *AuthClient) AddUser(ctx context.Context, user, space string) ([]byte, error) {
	params := make(url.Values)
	params.Add("user", user)
	params.Add("space", space)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+auth.Host+"/sys/auth/add", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.User, auth.Passwd)

	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	switch resp.StatusCode {
	case 200:
		return respBody, nil
	case 401:
		// auth failed
		return nil, ErrAuthFailed.As(resp.StatusCode, string(respBody))
	}
	return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
}

func (auth *AuthClient) ResetUserPasswd(ctx context.Context, user string) ([]byte, error) {
	params := make(url.Values)
	params.Add("user", user)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+auth.Host+"/sys/auth/reset", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.User, auth.Passwd)

	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	switch resp.StatusCode {
	case 200:
		return respBody, nil
	case 401:
		// auth failed
		return nil, ErrAuthFailed.As(resp.StatusCode, string(respBody))
	}
	return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
}

func (auth *AuthClient) ChangeAuth(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+auth.Host+"/sys/auth/change", nil)
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.User, auth.Passwd)

	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	switch resp.StatusCode {
	case 200:
		return respBody, nil
	case 401:
		// auth failed
		return nil, ErrAuthFailed.As(resp.StatusCode, string(respBody))
	}
	return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
}

func (auth *AuthClient) NewFileToken(ctx context.Context, authFile string) ([]byte, error) {
	params := url.Values{}
	params.Add("file", authFile)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+auth.Host+"/sys/file/token?"+params.Encode(), nil)
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.User, auth.Passwd)

	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
	}
	return respBody, nil
}

func (auth *AuthClient) DelayFileToken(ctx context.Context, authFile string) ([]byte, error) {
	params := url.Values{}
	params.Add("file", authFile)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+auth.Host+"/sys/file/token", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.User, auth.Passwd)

	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
	}
	return respBody, nil
}

func (auth *AuthClient) DeleteFileToken(ctx context.Context, authFile string) ([]byte, error) {
	params := url.Values{}
	params.Add("file", authFile)
	req, err := http.NewRequestWithContext(ctx, "DELETE", "https://"+auth.Host+"/sys/file/token?"+params.Encode(), nil)
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(auth.User, auth.Passwd)

	resp, err := utils.HttpsClient.Do(req)
	if err != nil {
		return nil, errors.As(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.As(err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Parse(string(respBody)).As(resp.StatusCode)
	}
	return respBody, nil
}
