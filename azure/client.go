package azure

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Resource struct {
	SubscriptionID string `json:"subscriptionId"`
}

type Resources struct {
	Resources []Resource `json:"value"`
}

type Client struct {
	URL string
}

func NewClient(serverURL string) *Client {
	return &Client{
		URL: serverURL,
	}
}

func (c Client) Resources() (Resources, error) {
	resp, err := http.Get(c.URL + "/subscriptions/")
	if err != nil {
		return Resources{}, err
	}
	if resp.StatusCode != http.StatusOK {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Resources{}, err
	}

	resources := &Resources{}
	err = json.Unmarshal(body, resources)
	if err != nil {
		return Resources{}, err
	}
	return *resources, nil
}
