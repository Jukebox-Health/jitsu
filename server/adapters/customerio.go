package adapters

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	customerEventsAPIURL = "https://track.customer.io/api/v1/customers/%v/events"
	eventsAPIURL         = "https://track.customer.io/api/v1/events"
)

//CustomerIORequest is a dto for sending requests to CustomerIO
type CustomerIORequest struct {
	SiteID  string                 `json:"site_id"`
	APIKey  string                 `json:"api_key"`
	Payload map[string]interface{} `json:"payload"`
}

//CustomerIOResponse is a dto for receiving response from CustomerIO
type CustomerIOResponse struct {
}

//CustomerIORequestFactory is a factory for building CustomerIO HTTP requests from input events
type CustomerIORequestFactory struct {
	siteID string
	apiKey string
}

//newCustomerIORequestFactory returns configured HTTPRequestFactory instance for CustomerIO requests
func newCustomerIORequestFactory(siteID, apiKey string) (HTTPRequestFactory, error) {
	return &CustomerIORequestFactory{siteID: siteID, apiKey: apiKey}, nil
}

func (arf *CustomerIORequestFactory) lookupCustomerIdentifier(payload map[string]interface{}) string {

	keys := []string{"id", "email", "cio_id"}

	for _, key := range keys {
		if _, exists := payload[key]; exists && payload[key] != nil {
			return fmt.Sprintf("%v", payload[key])
		}
	}

	return ""
}

//Create returns created CustomerIO request
//put empty array in body if object is nil (is used in test connection)
func (arf *CustomerIORequestFactory) Create(payload map[string]interface{}) (*Request, error) {

	var trackURL = eventsAPIURL

	if payload == nil {
		payload = make(map[string]interface{})
	}

	identifier := arf.lookupCustomerIdentifier(payload)

	if identifier != "" {
		trackURL = fmt.Sprintf(customerEventsAPIURL, identifier)
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling CustomerIO request [%v]: %v", payload, err)
	}

	creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", arf.siteID, arf.apiKey)))
	return &Request{
		URL:     trackURL,
		Method:  http.MethodPost,
		Body:    b,
		Headers: map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Basic %v", creds)},
	}, nil
}

func (arf *CustomerIORequestFactory) Close() {
}

//CustomerIOConfig is a dto for parsing CustomerIO configuration
type CustomerIOConfig struct {
	SiteID string `mapstructure:"site_id" json:"site_id,omitempty" yaml:"site_id,omitempty"`
	APIKey string `mapstructure:"api_key" json:"api_key,omitempty" yaml:"api_key,omitempty"`
}

//Validate returns err if invalid
func (ac *CustomerIOConfig) Validate() error {
	if ac == nil {
		return errors.New("customerio config is required")
	}
	if ac.SiteID == "" {
		return errors.New("'site_id' is required parameter")
	}

	if ac.APIKey == "" {
		return errors.New("'api_key' is required parameter")
	}

	return nil
}

//CustomerIO is an adapter for sending HTTP requests to CustomerIO
type CustomerIO struct {
	AbstractHTTP

	config *CustomerIOConfig
}

//NewCustomerIO returns configured CustomerIO adapter instance
func NewCustomerIO(config *CustomerIOConfig, httpAdapterConfiguration *HTTPAdapterConfiguration) (*CustomerIO, error) {
	httpReqFactory, err := newCustomerIORequestFactory(config.SiteID, config.APIKey)
	if err != nil {
		return nil, err
	}

	httpAdapterConfiguration.HTTPReqFactory = httpReqFactory
	httpAdapter, err := NewHTTPAdapter(httpAdapterConfiguration)
	if err != nil {
		return nil, err
	}

	a := &CustomerIO{config: config}
	a.httpAdapter = httpAdapter
	return a, nil
}

//NewTestCustomerIO returns test instance of adapter
func NewTestCustomerIO(config *CustomerIOConfig) *CustomerIO {
	return &CustomerIO{config: config}
}

//TestAccess sends test request (empty POST) to CustomerIO and check if error has occurred
func (a *CustomerIO) TestAccess() error {
	httpReqFactory, err := newCustomerIORequestFactory(a.config.SiteID, a.config.APIKey)
	if err != nil {
		return err
	}

	body := map[string]interface{}{"name": "connection_test", "anonymous_id": "jitsu_test_abc123"}
	r, err := httpReqFactory.Create(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return err
	}

	for k, v := range r.Headers {
		httpReq.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading CustomerIO response body: %v", err)
		}

		response := &CustomerIOResponse{}
		err = json.Unmarshal(responseBody, response)
		if err != nil {
			return fmt.Errorf("Error unmarshalling CustomerIO response body: %v", err)
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("error connecting to customerio")
		}

		return nil
	}

	return err
}

//Type returns adapter type
func (a *CustomerIO) Type() string {
	return "CustomerIO"
}
