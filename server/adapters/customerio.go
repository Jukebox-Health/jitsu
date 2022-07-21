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
	customerioTrackAPIURL = "https://track.customer.io/api/v1/"
)

//CustomerIORequest is a dto for sending requests to CustomerIO
type CustomerIORequest struct {
	SiteID string                   `json:"site_id"`
	APIKey string                   `json:"api_key"`
	Events []map[string]interface{} `json:"events"`
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

//Create returns created CustomerIO request
//put empty array in body if object is nil (is used in test connection)
func (arf *CustomerIORequestFactory) Create(object map[string]interface{}) (*Request, error) {
	//empty array is required. Otherwise nil will be sent (error)
	var eventsArr []map[string]interface{}
	if object != nil {
		eventsArr = append(eventsArr, object)
	}

	req := CustomerIORequest{SiteID: arf.siteID, APIKey: arf.apiKey, Events: eventsArr}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling CustomerIO request [%v]: %v", req, err)
	}

	creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", arf.siteID, arf.apiKey)))
	return &Request{
		URL:     customerioTrackAPIURL,
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

	body := map[string]interface{}{"event_type": "connection_test", "user_id": "hello@jitsu.com"}
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
