package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

// DefaultTimout for requests.
var DefaultTimout = time.Duration(10 * time.Second)

// Client to communicate with Raven Mailer API
type Client struct {
	endpoint *url.URL
	client   *http.Client
}

// Config parameters to configure a new HTTP client.
type Config struct {
	// Endpoint e.g. https://api.ravenmailer.com/v1
	Endpoint string

	// Timeout in seconds for request. Left unset defaults to
	Timeout time.Duration
}

// NewClient creates a new Raven Mailer HTTP client.
func NewClient(c Config) (*Client, error) {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 10,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   c.Timeout,
	}

	_, err := url.Parse(c.Endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "url parse")
	}

	return &Client{
		endpoint: u,
		client:   client,
	}, nil
}

func (c *Client) buildURL(path string, query url.Values) *url.URL {
	uri := url.URL{
		Scheme:     c.endpoint.Scheme,
		Host:       fmt.Sprintf("%s:%s", c.endpoint.Hostname(), c.endpoint.Port()),
		Path:       c.endpoint.Path + "/" + path,
		ForceQuery: false,
	}
	if query != nil {
		uri.RawQuery = query.Encode()
	}
	return &uri
}

// ListProjects fetches a slice of projects for a given user.
func (c *Client) ListProjects(ctx context.Context, userID string) ([]Project, error) {
	// build the URL including query params
	query := url.Values{
		"userId": []string{userID},
	}
	uri := c.buildURL("projects", query)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// json decode
	var container struct {
		Data []Project `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode list projects")
	}
	return container.Data, nil
}

// ListTransports fetches a slice of transports for the current project.
func (c *Client) ListTransports(ctx context.Context, projectID string) ([]Transport, error) {
	// build the URL including query params
	path := fmt.Sprintf("projects/%s/transports", projectID)
	uri := c.buildURL(path, nil)
	fmt.Printf("%s\n", uri)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	var container struct {
		Data []Transport `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode list transports")
	}
	return container.Data, nil
}

// CreateGroup creates a new named group.
func (c *Client) CreateGroup(ctx context.Context, projectID, name string) (*Group, error) {
	type createGroupRequest struct {
		ProjectID string `json:"projectId"`
		Name      string `json:"name"`
	}

	// request body
	req := createGroupRequest{
		ProjectID: projectID,
		Name:      name,
	}
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(&req); err != nil {
		return nil, err
	}

	// do request
	path := fmt.Sprintf("projects/%s/groups", projectID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodPost, uri.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "http post request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	// decode response
	return decodeGroupResponse(res.Body)
}

// ListGroups fetches a slice of groups for the current project.
func (c *Client) ListGroups(ctx context.Context, projectID string) ([]Group, error) {
	// build the URL including query params
	path := fmt.Sprintf("projects/%s/groups", projectID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	// json decode
	var container struct {
		Data []Group `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode list groups")
	}
	return container.Data, nil
}

// GetGroup fetches a single group by id.
func (c *Client) GetGroup(ctx context.Context, projectID, groupID string) (*Group, error) {
	path := fmt.Sprintf("projects/%s/groups/%s", projectID, groupID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	return decodeGroupResponse(res.Body)
}

// DeleteGroup deletes a group by id.
func (c *Client) DeleteGroup(ctx context.Context, projectID, groupID string) error {
	path := fmt.Sprintf("projects/%s/groups/%s", projectID, groupID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodDelete, uri.String(), nil)
	if err != nil {
		return errors.Wrap(err, "http delete request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return decodeAPIError(res.Body)
	}

	return nil
}

// ListMail fetches a list of mail resources.
func (c *Client) ListMail(ctx context.Context, projectID string) ([]Mail, error) {
	// build the URL including query params
	path := fmt.Sprintf("projects/%s/mail", projectID)
	uri := c.buildURL(path, nil)

	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	// json decode
	var container struct {
		Data []Mail `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode list mail")
	}
	return container.Data, nil
}

// ListMailLogs fetches a list of mail log resources for the given mail entry.
func (c *Client) ListMailLogs(ctx context.Context, projectID, mailID string) ([]MailLog, error) {
	// build the URL including query params
	path := fmt.Sprintf("projects/%s/mail/%s/logs", projectID, mailID)
	uri := c.buildURL(path, nil)

	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	// json decode
	var container struct {
		Data []MailLog `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode list mail logs")
	}
	return container.Data, nil
}

// GetMail fetches a single mail resource.
func (c *Client) GetMail(ctx context.Context, projectID, mailID string) (*Mail, error) {
	path := fmt.Sprintf("projects/%s/mail/%s", projectID, mailID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	return decodeMailResponse(res.Body)
}

func decodeGroupResponse(r io.Reader) (*Group, error) {
	var container struct {
		Data *Group `json:"data"`
	}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode get group")
	}
	return container.Data, nil
}

func decodeMailResponse(r io.Reader) (*Mail, error) {
	var container struct {
		Data *Mail `json:"data"`
	}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode get mail")
	}
	return container.Data, nil
}

// CreateTemplate HTTP POST /templates/{id}
func (c *Client) CreateTemplate(ctx context.Context, params *CreateTemplateParams) (*Template, error) {
	path := fmt.Sprintf("projects/%s/templates/%s", params.ProjectID, params.ID)
	uri := c.buildURL(path, nil)

	// request body
	req := createTemplateRequest{
		GroupID: params.GroupID,
		Txt:     params.Txt,
	}
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(&req); err != nil {
		return nil, err
	}

	// post
	res, err := c.request(http.MethodPut, uri.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "http post request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	var container struct {
		Data *Template `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode get template")
	}
	return container.Data, nil
}

// GetTemplate fetches a single template by id.
func (c *Client) GetTemplate(ctx context.Context, projectID, templateID string) (*Template, error) {
	path := fmt.Sprintf("projects/%s/templates/%s", projectID, templateID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return nil, decodeAPIError(res.Body)
	}

	var container struct {
		Data *Template `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode get template")
	}
	return container.Data, nil
}

// ListTemplates fetches a slice of templates for the current project.
func (c *Client) ListTemplates(ctx context.Context, projectID string) ([]Template, error) {
	// build the URL including query params
	path := fmt.Sprintf("projects/%s/templates", projectID)
	uri := c.buildURL(path, nil)
	res, err := c.request(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "http get request failed")
	}
	defer res.Body.Close()

	var container struct {
		Data []Template `json:"data"`
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "json decode list templates")
	}
	return container.Data, nil
}

// DeleteTemplate deletes the template with the given id.
func (c *Client) DeleteTemplate(ctx context.Context, projectID, templateID string) error {
	path := fmt.Sprintf("projects/%s/templates/%s", projectID, templateID)
	uri := c.buildURL(path, nil)
	fmt.Println(uri)

	res, err := c.request(http.MethodDelete, uri.String(), nil)
	if err != nil {
		return errors.Wrapf(err, "http delete template (%s) request failed", templateID)
	}
	defer res.Body.Close()

	// 4xx range
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return decodeAPIError(res.Body)
	}

	return nil
}

func (c *Client) request(method, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, body)
	if method != http.MethodDelete {
		req.Header.Set("Accept", "application/json")
	}
	// req.Header.Set("Authorization", "Bearer "+c.jwt)

	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "do HTTP %s request", req.Method)
	}
	return res, nil
}

func decodeAPIError(r io.Reader) error {
	var apiErr APIError
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&apiErr); err != nil {
		return err
	}
	return &apiErr
}
