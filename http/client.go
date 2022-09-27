package http

import (
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

// ListTemplates fetches a slice of templates for the current project.
func (c *Client) ListTemplates(ctx context.Context) ([]Template, error) {
	// build the URL including query params
	query := url.Values{"projectId": []string{"project-1"}}
	uri := c.buildURL("templates", query)

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

func (c *Client) request(method, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, body)
	req.Header.Set("Accept", "application/json")
	// req.Header.Set("Authorization", "Bearer "+c.jwt)

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "do HTTP %s request", req.Method)
	}
	return res, nil
}
