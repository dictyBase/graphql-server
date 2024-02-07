package authentication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dictyBase/graphql-server/internal/repository"
	"golang.org/x/exp/slices"
)

type LogtoClient interface {
	AccessToken() (*AccessTokenResp, error)
	CheckUserWithUserName(string) (bool, string, error)
	UserWithEmail(string) (*UserResp, error)
	CheckUser(string) (bool, string, error)
	AddCustomUserInformation(string, string, *APIUsersPatchCustomData) error
	User(string) (*UserResp, error)
	Roles(string) ([]*RoleResp, error)
	Permissions(string) ([]*PermissionResp, error)
	CreateUser(string, *APIUsersPostReq) (string, error)
}

type authClient struct {
	baseURL     string
	httpClient  *http.Client
	appID       string
	appSecret   string
	apiResource string
	cache       repository.Repository
	cacheKey    string
}

type LogtoClientParams struct {
	URL, AppID, AppSecret, APIResource, Key string
	TokenCache                              repository.Repository
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type UserResp struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PrimaryEmail string `json:"primaryEmail"`
	PrimaryPhone string `json:"primaryPhone"`
	Name         string `json:"name"`
	Avatar       any    `json:"avatar,omitempty"`
	CustomData   struct {
		City             string `json:"city,omitempty"`
		Phone            string `json:"phone,omitempty"`
		State            string `json:"state,omitempty"`
		Region           string `json:"region,omitempty"`
		Address          string `json:"address,omitempty"`
		Country          string `json:"country,omitempty"`
		Zipcode          string `json:"zipcode,omitempty"`
		JobTitle         string `json:"job_title,omitempty"`
		Profession       string `json:"profession,omitempty"`
		Subscribed       bool   `json:"subscribed,omitempty"`
		Institution      string `json:"institution,omitempty"`
		ResearchInterest string `json:"research_interest,omitempty"`
		SecondaryAddress string `json:"secondary_address,omitempty"`
	} `json:"customData"`
	Identities struct {
	} `json:"identities,omitempty"`
	CreatedAt     CustomTime `json:"createdAt"`
	ApplicationID any        `json:"applicationId,omitempty"`
	IsSuspended   bool       `json:"isSuspended,omitempty"`
}

type APIUsersPostReq struct {
	PrimaryPhone string `json:"primaryPhone"`
	PrimaryEmail string `json:"primaryEmail"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
}

type APIUsersPostRes struct {
	ID string `json:"id"`
}

type APIUsersSearchRes struct {
	Email    string `json:"primaryEmail"`
	ID       string `json:"id"`
	UserName string `json:"username"`
}

type APIUsersPatchCustomData struct {
	CustomData AdditionalUserInformation `json:"customData"`
}

type AdditionalUserInformation struct {
	Profession       string `json:"profession"`
	JobTitle         string `json:"job_title"`
	Institution      string `json:"institution"`
	Address          string `json:"address"`
	SecondaryAddress string `json:"secondary_address"`
	City             string `json:"city"`
	State            string `json:"state"`
	Region           string `json:"region"`
	Country          string `json:"country"`
	Zipcode          string `json:"zipcode"`
	Subscribed       bool   `json:"subscribed"`
	ResearchInterest string `json:"research_interest"`
	Phone            string `json:"phone"`
}

type RoleResp struct {
	TenantID    string `json:"tenantId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rtype       string `json:"type"`
}

type PermissionResp struct {
	TenantID    string    `json:"tenantId"`
	ID          string    `json:"id"`
	ResourceID  string    `json:"resourceId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Resource    struct {
		TenantID       string `json:"tenantId"`
		ID             string `json:"id"`
		Name           string `json:"name"`
		Indicator      string `json:"indicator"`
		IsDefault      bool   `json:"isDefault"`
		AccessTokenTTL int    `json:"accessTokenTtl"`
	} `json:"resource"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(input []byte) error {
	sinput := strings.Trim(string(input), `"`)
	num, err := strconv.ParseInt(sinput, 10, 64)
	if err != nil {
		return fmt.Errorf("error converting string to int64:%s", err)
	}
	tstamp := time.UnixMilli(num)
	ct.Time = tstamp
	return nil
}

// NewClient creates a new instance of the Client struct.
// It takes an endpoint string as a parameter and returns a pointer to the Client struct.
func NewClient(params *LogtoClientParams) LogtoClient {
	return &authClient{
		baseURL:     params.URL,
		httpClient:  &http.Client{},
		appID:       params.AppID,
		appSecret:   params.AppSecret,
		apiResource: params.APIResource,
		cache:       params.TokenCache,
		cacheKey:    params.Key,
	}
}

func (clnt *authClient) AccessToken() (*AccessTokenResp, error) {
	acresp := &AccessTokenResp{}
	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("resource", clnt.apiResource)
	params.Set("scope", "all")
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oidc/token", clnt.baseURL),
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return acresp, fmt.Errorf("error in creating request %s ", err)
	}
	req.SetBasicAuth(clnt.appID, clnt.appSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := clnt.reqToResponse(req)
	if err != nil {
		return acresp, err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(acresp); err != nil {
		return acresp, fmt.Errorf("error in decoding json response %s", err)
	}

	return acresp, nil
}

func (clnt *authClient) CheckUserWithUserName(
	username string,
) (bool, string, error) {
	var userStruct string
	usrs, err := clnt.checkUser("username", username)
	if err != nil {
		return false, userStruct, fmt.Errorf(
			"error in finding user name %s",
			err,
		)
	}
	index := slices.IndexFunc(usrs, func(usr *APIUsersSearchRes) bool {
		return usr.UserName == username
	})
	if index == -1 {
		return false, userStruct, nil
	}
	return true, usrs[index].ID, nil
}

func (clnt *authClient) CheckUser(email string) (bool, string, error) {
	var userStruct string
	usrs, err := clnt.checkUser("primaryEmail", email)
	if err != nil {
		return false, userStruct, fmt.Errorf(
			"error in finding user email %s",
			err,
		)
	}
	index := slices.IndexFunc(usrs, func(usr *APIUsersSearchRes) bool {
		return usr.Email == email
	})
	if index == -1 {
		return false, userStruct, nil
	}
	return true, usrs[index].ID, nil
}

func (clnt *authClient) UserWithEmail(email string) (*UserResp, error) {
	userStruct := &UserResp{}
	token, err := clnt.retrieveToken()
	if err != nil {
		return userStruct, fmt.Errorf("error in getting token %s", err)
	}
	params := url.Values{}
	params.Set("search.primaryEmail", email)
	params.Set("mode.name", "exact")
	parsedURL, err := url.Parse(fmt.Sprintf("%s/api/users", clnt.baseURL))
	if err != nil {
		return userStruct, fmt.Errorf(
			"error in parsing url for query %s",
			err,
		)
	}
	parsedURL.RawQuery = params.Encode()
	ureq, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return userStruct, fmt.Errorf(
			"error in making new request %s",
			err,
		)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return userStruct, err
	}
	defer uresp.Body.Close()
	usrs := make([]*UserResp, 0)
	if err := json.NewDecoder(uresp.Body).Decode(&usrs); err != nil {
		return userStruct, fmt.Errorf(
			"error in decoding json response %s",
			err,
		)
	}
	index := slices.IndexFunc(usrs, func(usr *UserResp) bool {
		return usr.PrimaryEmail == email
	})
	if index == -1 {
		return userStruct, fmt.Errorf(
			"user with email %s not found %s",
			email,
			err,
		)
	}

	return usrs[index], nil
}

func (clnt *authClient) AddCustomUserInformation(
	token,
	userStruct string,
	user *APIUsersPatchCustomData,
) error {
	content, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error in converting to json %s", err)
	}
	ureq, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/api/users/%s/custom-data", clnt.baseURL, userStruct),
		bytes.NewBuffer(content),
	)
	if err != nil {
		return fmt.Errorf(
			"error in making new request for user custom data%s ",
			err,
		)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return fmt.Errorf("error in adding user custom information %s", err)
	}
	defer uresp.Body.Close()
	return nil
}

func (clnt *authClient) User(userID string) (*UserResp, error) {
	userStruct := &UserResp{}
	token, err := clnt.retrieveToken()
	if err != nil {
		return userStruct, fmt.Errorf("error in getting token %s", err)
	}
	ureq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/users/%s", clnt.baseURL, userID),
		nil,
	)
	if err != nil {
		return userStruct, fmt.Errorf(
			"error in making new request for fetching user %s",
			err,
		)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return userStruct, fmt.Errorf(
			"error in getting response during fetching of user %s",
			err,
		)
	}
	defer uresp.Body.Close()
	if err := json.NewDecoder(uresp.Body).Decode(&userStruct); err != nil {
		return userStruct, fmt.Errorf(
			"error in decoding json response %s",
			err,
		)
	}

	return userStruct, nil
}

func (clnt *authClient) Roles(userStruct string) ([]*RoleResp, error) {
	rolesStruct := make([]*RoleResp, 0)
	token, err := clnt.retrieveToken()
	if err != nil {
		return rolesStruct, fmt.Errorf("error in getting token %s", err)
	}
	ureq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/users/%s/roles", clnt.baseURL, userStruct),
		nil,
	)
	if err != nil {
		return rolesStruct, fmt.Errorf(
			"error in making new request for fetching user roles %s",
			err,
		)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return rolesStruct, fmt.Errorf(
			"error in getting response during fetching of roles %s",
			err,
		)
	}
	defer uresp.Body.Close()
	if err := json.NewDecoder(uresp.Body).Decode(&rolesStruct); err != nil {
		return rolesStruct, fmt.Errorf(
			"error in decoding json response %s",
			err,
		)
	}

	return rolesStruct, nil
}

func (clnt *authClient) Permissions(roleID string) ([]*PermissionResp, error) {
	permissionStruct := make([]*PermissionResp, 0)
	token, err := clnt.retrieveToken()
	if err != nil {
		return permissionStruct, fmt.Errorf("error in getting token %s", err)
	}
	ureq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/roles/%s/scopes", clnt.baseURL, roleID),
		nil,
	)
	if err != nil {
		return permissionStruct, fmt.Errorf(
			"error in making new request for fetching permissions %s",
			err,
		)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return permissionStruct, fmt.Errorf(
			"error in getting response during fetching of permissions %s",
			err,
		)
	}
	defer uresp.Body.Close()
	if err := json.NewDecoder(uresp.Body).Decode(&permissionStruct); err != nil {
		return permissionStruct, fmt.Errorf(
			"error in decoding json response %s",
			err,
		)
	}

	return permissionStruct, nil
}

func (clnt *authClient) CreateUser(
	token string,
	user *APIUsersPostReq,
) (string, error) {
	var userStruct string
	content, err := json.Marshal(user)
	if err != nil {
		return userStruct, fmt.Errorf("error in converting to json %s", err)
	}
	ureq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/users", clnt.baseURL),
		bytes.NewBuffer(content),
	)
	if err != nil {
		return userStruct, fmt.Errorf("error in making new request %s", err)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return userStruct, err
	}
	defer uresp.Body.Close()
	usr := &APIUsersPostRes{}
	if err := json.NewDecoder(uresp.Body).Decode(usr); err != nil {
		return userStruct, fmt.Errorf("error in decoding json response %s", err)
	}
	return usr.ID, nil
}

func (clnt *authClient) retrieveToken() (string, error) {
	var token string
	ok, err := clnt.cache.Exists(clnt.cacheKey)
	if err != nil {
		return token, fmt.Errorf("error in finding token key %s", err)
	}
	if ok {
		token, err := clnt.cache.Get(clnt.cacheKey)
		if err != nil {
			return token, fmt.Errorf(
				"error in fetching token from cache %s",
				err,
			)
		}
	}
	aresp, err := clnt.AccessToken()
	if err != nil {
		return token, fmt.Errorf("error in retrieving access token %s", err)
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", aresp.ExpiresIn-1000))
	if err != nil {
		return token, fmt.Errorf(
			"error in parsing duration %d",
			aresp.ExpiresIn,
		)
	}
	err = clnt.cache.SetWithTTL(token, aresp.AccessToken, dur)
	if err != nil {
		return token, fmt.Errorf("error in setting token %s", err)
	}

	return aresp.AccessToken, nil
}

func commonHeader(lreq *http.Request, token string) {
	lreq.Header.Set("Content-Type", "application/json")
	lreq.Header.Set("Accept", "application/json")
	lreq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (clnt *authClient) reqToResponse(
	creq *http.Request,
) (*http.Response, error) {
	uresp, err := clnt.httpClient.Do(creq)
	if err != nil {
		return uresp, fmt.Errorf("error in making request %s", err)
	}
	if uresp.StatusCode != 200 {
		cnt, err := io.ReadAll(uresp.Body)
		if err != nil {
			return uresp, fmt.Errorf(
				"error in response and the reading the body %d %s",
				uresp.StatusCode,
				err,
			)
		}
		return uresp, fmt.Errorf(
			"unexpected error response %d %s",
			uresp.StatusCode,
			string(cnt),
		)
	}
	return uresp, nil
}

func (clnt *authClient) checkUser(
	queryField string,
	value string,
) ([]*APIUsersSearchRes, error) {
	usrs := make([]*APIUsersSearchRes, 0)
	token, err := clnt.retrieveToken()
	if err != nil {
		return usrs, fmt.Errorf("error in getting token %s", err)
	}
	params := url.Values{}
	params.Set(fmt.Sprintf("search.%s", queryField), value)
	params.Set("mode.name", "exact")
	parsedURL, err := url.Parse(fmt.Sprintf("%s/api/users", clnt.baseURL))
	if err != nil {
		return usrs, fmt.Errorf(
			"error in parsing url for query %s",
			err,
		)
	}
	parsedURL.RawQuery = params.Encode()
	ureq, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return usrs, fmt.Errorf(
			"error in making new request %s",
			err,
		)
	}
	commonHeader(ureq, token)
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return usrs, err
	}
	defer uresp.Body.Close()
	if err := json.NewDecoder(uresp.Body).Decode(&usrs); err != nil {
		return usrs, fmt.Errorf(
			"error in decoding json response %s",
			err,
		)
	}

	return usrs, nil
}
