package custom

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aliyun/saml2alibabacloud/pkg/cfg"
	"github.com/aliyun/saml2alibabacloud/pkg/creds"
	"github.com/aliyun/saml2alibabacloud/pkg/provider"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var logger = logrus.WithField("provider", "custom")

// Client is a wrapper representing a custom SAML client
type Client struct {
	client *provider.HTTPClient
	mfa    string
}

// New creates a new custom client
func New(idpAccount *cfg.IDPAccount) (*Client, error) {

	tr := provider.NewDefaultTransport(idpAccount.SkipVerify)

	client, err := provider.NewHTTPClient(tr, provider.BuildHttpClientOpts(idpAccount))
	if err != nil {
		return nil, errors.Wrap(err, "error building http client")
	}

	// assign a response validator to ensure all responses are either success or a redirect
	// this is to avoid have explicit checks for every single response
	client.CheckResponseStatus = provider.SuccessOrRedirectResponseValidator

	return &Client{
		client: client,
		// TODO currently not supported
		mfa: idpAccount.MFA,
	}, nil
}

// Authenticate using an API endpoint with username and password then returns a SAML response
func (oc *Client) Authenticate(loginDetails *creds.LoginDetails) (string, error) {

	_, err := url.Parse(loginDetails.URL)
	if err != nil {
		return "", errors.Wrap(err, "error building login request URL")
	}

	//authenticate using x-www-form-urlencoded
	authReq := url.Values{}
	authReq.Set("username", loginDetails.Username)
	authReq.Set("password", loginDetails.Password)

	authBody := strings.NewReader(authReq.Encode())

	req, err := http.NewRequest("POST", loginDetails.URL, authBody)
	if err != nil {
		return "", errors.Wrap(err, "error building authentication request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(authReq.Encode())))

	res, err := oc.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error retrieving auth response")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "error retrieving body from response")
	}

	resp := string(body)

	successResponse := gjson.Get(resp, "success").String()
	samlResponse := gjson.Get(resp, "data").String()

	// error response
	if successResponse != "true" {
		return "", errors.Wrap(err, "error retrieving SAML response")
	}

	decodedSamlResponse, err := base64.StdEncoding.DecodeString(samlResponse)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode SAML response")
	}
	logger.WithField("type", "saml-response").WithField("saml-response", string(decodedSamlResponse)).Debug("custom auth response")
	return samlResponse, nil
}
