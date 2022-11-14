package sunshine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseUrl   string
	appId     string
	apiKeyId  string
	appSecret string
	client    *http.Client
}

func NewClient(baseUrl, appId, apiKeyId, appSecret string) *Client {
	client := http.DefaultClient
	return &Client{
		baseUrl,
		appId,
		apiKeyId,
		appSecret,
		client,
	}
}

type CustomIntegration struct {
	Id          string    `json:"id"`
	Status      string    `json:"status"`
	DisplayName string    `json:"displayName"`
	Webhooks    []Webhook `json:"webhooks"`
}

type Webhook struct {
	Id       string   `json:"id"`
	Version  string   `json:"version"`
	Target   string   `json:"target"`
	Triggers []string `json:"triggers"`
	Secret   string   `json:"secret"`
}

type createCustomIntegrationRequest struct {
	Type        string          `json:"type"`
	DisplayName string          `json:"displayName"`
	Webhooks    []createWebhook `json:"webhooks"`
}

type createWebhook struct {
	Target            string   `json:"target"`
	Triggers          []string `json:"triggers"`
	IncludeFullUser   bool     `json:"includeFullUser"`
	IncludeFullSource bool     `json:"includeFullSource"`
}

type customIntegrationWrapper struct {
	Integration CustomIntegration `json:"integration"`
}

func (client *Client) CreateCustomIntegration(displayName string, target string, triggers []string, includeFullUser bool, includeFullSource bool) (*CustomIntegration, error) {
	requestBody := createCustomIntegrationRequest{
		Type:        "custom",
		DisplayName: displayName,
		Webhooks: []createWebhook{
			{
				Target:            target,
				Triggers:          triggers,
				IncludeFullUser:   includeFullUser,
				IncludeFullSource: includeFullSource,
			},
		},
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/apps/%s/integrations", client.baseUrl, client.appId), bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(client.apiKeyId, client.appSecret)

	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var wrapper customIntegrationWrapper
	err = json.Unmarshal(body, &wrapper)

	return &wrapper.Integration, err
}

type updateCustomIntegration struct {
	DisplayName string `json:"displayName"`
}

func (client *Client) UpdateCustomIntegration(integrationId string, displayName string) error {
	requestBody := updateCustomIntegration{DisplayName: displayName}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%sv2/apps/%s/integrations/%s", client.baseUrl, client.appId, integrationId), bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return err
	}

	request.SetBasicAuth(client.apiKeyId, client.appSecret)

	_, err = client.client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteCustomIntegration(integrationId string) error {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%sv2/apps/%s/integrations/%s", client.baseUrl, client.appId, integrationId), bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	request.SetBasicAuth(client.apiKeyId, client.appSecret)

	_, err = client.client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

type customIntegrationsWrapper struct {
	Integrations []CustomIntegration `json:"integrations"`
}

func (client *Client) GetAllCustomIntegrations() (*[]CustomIntegration, error) {
	url := fmt.Sprintf("%sv2/apps/%s/integrations?filter[types]=custom", client.baseUrl, client.appId)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(client.apiKeyId, client.appSecret)

	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var wrapper customIntegrationsWrapper
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return nil, err
	}

	return &wrapper.Integrations, nil
}

func (client *Client) GetCustomIntegration(integrationId string) (*CustomIntegration, error) {
	url := fmt.Sprintf("%sv2/apps/%s/integrations/%s", client.baseUrl, client.appId, integrationId)
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(client.apiKeyId, client.appSecret)

	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var wrapper customIntegrationWrapper
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return nil, err
	}

	return &wrapper.Integration, nil
}
