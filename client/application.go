// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "krak8s": application Resource Client
//
// Command:
// $ goagen
// --design=krak8s/design
// --out=$(GOPATH)/src/krak8s
// --version=v1.2.0-dirty

package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CreateApplicationPath computes a request path to the create action of application.
func CreateApplicationPath(projectid string) string {
	param0 := projectid

	return fmt.Sprintf("/v1/projects/%s/applications", param0)
}

// Request the creation of an application deployment in the project/namespace
func (c *Client) CreateApplication(ctx context.Context, path string, payload *ApplicationPostBody) (*http.Response, error) {
	req, err := c.NewCreateApplicationRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateApplicationRequest create the request corresponding to the create action endpoint of the application resource.
func (c *Client) NewCreateApplicationRequest(ctx context.Context, path string, payload *ApplicationPostBody) (*http.Request, error) {
	var body bytes.Buffer
	err := c.Encoder.Encode(payload, &body, "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	header.Set("Content-Type", "application/json")
	return req, nil
}

// DeleteApplicationPath computes a request path to the delete action of application.
func DeleteApplicationPath(projectid string, appid string) string {
	param0 := projectid
	param1 := appid

	return fmt.Sprintf("/v1/projects/%s/applications/%s", param0, param1)
}

// Delete the specified application from the project/namespace
func (c *Client) DeleteApplication(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDeleteApplicationRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDeleteApplicationRequest create the request corresponding to the delete action endpoint of the application resource.
func (c *Client) NewDeleteApplicationRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// GetApplicationPath computes a request path to the get action of application.
func GetApplicationPath(projectid string, appid string) string {
	param0 := projectid
	param1 := appid

	return fmt.Sprintf("/v1/projects/%s/applications/%s", param0, param1)
}

// Get the status of the specified application in the project/namespace
func (c *Client) GetApplication(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetApplicationRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetApplicationRequest create the request corresponding to the get action endpoint of the application resource.
func (c *Client) NewGetApplicationRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ListApplicationPayload is the application list action payload.
type ListApplicationPayload struct {
	Namespaceid string `form:"namespaceid" json:"namespaceid" xml:"namespaceid"`
}

// ListApplicationPath computes a request path to the list action of application.
func ListApplicationPath(projectid string) string {
	param0 := projectid

	return fmt.Sprintf("/v1/projects/%s/applications", param0)
}

// Retrieve the collection of all applications in the project/namespace.
func (c *Client) ListApplication(ctx context.Context, path string, payload *ListApplicationPayload) (*http.Response, error) {
	req, err := c.NewListApplicationRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListApplicationRequest create the request corresponding to the list action endpoint of the application resource.
func (c *Client) NewListApplicationRequest(ctx context.Context, path string, payload *ListApplicationPayload) (*http.Request, error) {
	var body bytes.Buffer
	err := c.Encoder.Encode(payload, &body, "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	header.Set("Content-Type", "application/json")
	return req, nil
}
