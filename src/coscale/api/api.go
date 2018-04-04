package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Object represents an CoScale API Object with an Id.
type Object interface {
	// GetId returns the Id of the CoScale API Object.
	GetId() int64
}

const (
	// DEFAULT_STRING_VALUE defines the default value for a string flag.
	DEFAULT_STRING_VALUE string = `!>dUmmy<!`
	// DEFAULT_INT64_VALUE defines the default value for an integer flag.
	DEFAULT_INT64_VALUE int64 = math.MinInt64
	// MaxUploadSize for data call.
	MaxUploadSize int = 10 * 1024 * 1024 // 10MB
)

const (
	connectionTimeout time.Duration = 30 * time.Second
	readWriteTimeout  time.Duration = 30 * time.Second
	downloadTimeout   time.Duration = 300 * time.Second
)

// AuthenticationError is returned by the API when a login fails.
type AuthenticationError string

// UnauthorizedError is returned by the API when the session has insufficient acces rights.
type UnauthorizedError string

// NotFoundError is returned by the API when the requested object was not found.
type NotFoundError string

// RequestError is returned by the API for invalid requests.
type RequestError string

// Duplicate is returned by the API when a duplicate is detected during insert.
type Duplicate int64

// Disabled is returned by the API when the feature is disabled.
type Disabled string

// InvalidConfig is returned when the API connection configuration is not valid.
type InvalidConfig string

func (ue AuthenticationError) Error() string {
	return string(ue)
}

func (ue UnauthorizedError) Error() string {
	return string(ue)
}

func (nfe NotFoundError) Error() string {
	return string(nfe)
}

func (re RequestError) Error() string {
	return string(re)
}

func (d Duplicate) Error() string {
	return fmt.Sprintf("Duplicate with id %d", d)
}

func (d Disabled) Error() string {
	return string(d)
}

func (i InvalidConfig) Error() string {
	return string(i)
}

// IsAuthenticationError checks if an error is a AuthenticationError.
func IsAuthenticationError(err error) bool {
	_, ok := err.(AuthenticationError)
	return ok
}

// IsNotFoundError checks if an error is an NotFoundError.
func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

// IsRequestError checks if an error is a RequestError.
func IsRequestError(err error) bool {
	_, ok := err.(RequestError)
	return ok
}

// IsDuplicate checks if the error is a duplicate error, if true, this also returns the duplicate id.
func IsDuplicate(err error) (bool, int64) {
	d, ok := err.(Duplicate)
	if ok {
		return ok, int64(d)
	}
	return ok, 0
}

// IsDisabled checks if an error indicates that the feature was disabled.
func IsDisabled(err error) bool {
	_, ok := err.(Disabled)
	return ok
}

// IsInvalidConfig checks if an error is caused by an invalid api configuration.
func IsInvalidConfig(err error) bool {
	_, ok := err.(InvalidConfig)
	return ok
}

// Api contains the Api connection configuration.
type Api struct {
	// BaseUrl is in the following format: https://<coscale-hostname>
	BaseUrl string
	// AccessToken is a UUID giving access permissions on the application.
	AccessToken string
	// AppID is a UUID defining the application.
	AppID       string
	rawOutput   bool
	token       string
	validConfig bool
}

// NewApi creates a new Api connector using an email and a password.
func NewApi(baseUrl string, accessToken string, appID string, rawOutput bool) *Api {
	api := &Api{baseUrl, accessToken, appID, rawOutput, "", true}
	return api
}

// NewFakeApi creates a new Api connector using an email and a password.
func NewFakeApi() *Api {
	api := &Api{"", "", "", true, "", false}
	return api
}

// GetSource gets the source name for the CoScale cli.
func GetSource() string {
	return "CLI"
}

// newTimeoutDialer creates a new Dailer with the given timeouts.
func newTimeoutDialer(connectionTimeout time.Duration, readWriteTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectionTimeout)
		if err != nil {
			return nil, err
		}
		if err = conn.SetDeadline(time.Now().Add(readWriteTimeout)); err != nil {
			return nil, err
		}
		return conn, nil
	}
}

// newTimeoutClient creates a http client with timeouts on the connection an reads/writes.
func newTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial:  newTimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}

// RequestErrorResponse is response in case of a RequestError.
type RequestErrorResponse struct {
	Msg  string
	Type string
	ID   int64
}

// Do an http request.
func (api *Api) doHttpRequest(method string, uri string, token string, data map[string][]string, timeout time.Duration) ([]byte, error) {
	requestBody := url.Values(data).Encode()
	req, err := http.NewRequest(method, uri, strings.NewReader(requestBody))

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "CoScale CLI")

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if token != "" {
		req.Header.Add("HTTPAuthorization", token)
	}

	client := newTimeoutClient(connectionTimeout, timeout)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return nil, UnauthorizedError(fmt.Sprintf("Unauthorized: %s", body))
	} else if resp.StatusCode == 404 {
		return nil, NotFoundError(fmt.Sprintf("Instance not found: %s", body))
	} else if resp.StatusCode == 409 {
		var reqErr RequestErrorResponse
		if jsonErr := json.Unmarshal(body, &reqErr); jsonErr == nil {
			if reqErr.Type == "DUPLICATE" {
				return nil, Duplicate(reqErr.ID)
			} else if reqErr.Type == "DISABLED" {
				return nil, Disabled(reqErr.Msg)
			}
		}
		return nil, RequestError(fmt.Sprintf("Request error: %s", body))
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received bad status code: %d -- %s", resp.StatusCode, body)
	}

	if err != nil {
		return nil, err
	}

	return body, nil
}

// LoginData contains the required fields for the login API function.
type LoginData struct {
	// Token should contain the AccessToken.
	Token string
}

// Login to the Api, returns the token and an error.
func (api *Api) Login() error {

	data := map[string][]string{
		"accessToken": {api.AccessToken},
	}

	var loginData LoginData
	bytes, err := api.doHttpRequest("POST", fmt.Sprintf("%s/api/v1/app/%s/login/", api.BaseUrl, api.AppID), "", data, readWriteTimeout)
	if err != nil {
		return AuthenticationError(fmt.Sprintf("Authentication error: %s", err.Error()))
	}

	if err := json.Unmarshal(bytes, &loginData); err != nil {
		return err
	}

	api.token = loginData.Token
	return nil
}

// Make a call to the api, returns the bytes returned.
func (api *Api) makeRawCall(method string, uri string, data map[string][]string, timeout time.Duration) ([]byte, error) {
	// Not authenticated yet, try login.
	if api.token == "" {
		if err := api.Login(); err != nil {
			return nil, err
		}
	}

	// Do the actual request.
	bytes, err := api.doHttpRequest(method, api.BaseUrl+uri, api.token, data, timeout)
	if err != nil {
		if _, ok := err.(UnauthorizedError); ok {
			// unauthorizedError: the token might have experied. Performing login again
			// and retrying the request.
			if err := api.Login(); err != nil {
				api.token = ""
				return nil, err
			}
			return api.doHttpRequest(method, api.BaseUrl+uri, api.token, data, timeout)
		}
		return bytes, err
	}

	return bytes, nil
}

// makeCall Make a call to the api and parse the json response into target.
func (api *Api) makeCall(method string, uri string, data map[string][]string, jsonOut bool, target interface{}) error {
	if !api.validConfig {
		return InvalidConfig("Could not find valid authentication configuration.")
	}

	b, err := api.makeRawCall(method, uri, data, readWriteTimeout)
	if err != nil {
		return err
	}
	if target != nil {
		if jsonOut {
			//if the output will be json, check if we need to format it or no
			var result string
			if api.rawOutput {
				result = string(b)
			} else {
				var out bytes.Buffer
				json.Indent(&out, b, "", " ")
				result = out.String()
			}
			if t, ok := target.(*string); ok {
				*t = result
			}
		} else {
			return json.Unmarshal(b, target)
		}
	}
	return nil
}

// GetObjectGroupName will add group sufix to a object name.
func GetObjectGroupName(objectName string) (objectGroupName string) {
	objectGroupName = objectName

	// Make it compatible for both adding object to group and also group to group.
	if !strings.HasSuffix(objectGroupName, "group") {
		objectGroupName += "group"
	}
	return
}
