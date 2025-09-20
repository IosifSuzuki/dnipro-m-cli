package client

import (
	"dniprom-cli/internal/container"
	"dniprom-cli/internal/model/network"
	"dniprom-cli/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	SearchAPIEndpoint   = "search/drop-down/"
	WarrantyAPIEndpoint = "shop/catalog/get-product-service-maintenance/"
)

type DniproClient interface {
	FetchAutocompleteProduct(code string) (*network.Product, error)
	GetWarranty(id int64) (string, error)
}

type dniproClient struct {
	container container.Container
	client    *http.Client
}

func NewDniproClient(container container.Container) DniproClient {
	return &dniproClient{
		container: container,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d *dniproClient) FetchAutocompleteProduct(code string) (*network.Product, error) {
	log := d.container.GetLogger()
	fullPath := d.GetPath(SearchAPIEndpoint)

	u, err := url.Parse(fullPath)
	if err != nil {
		log.Error("error parsing url", logger.FError(err))
		return nil, err
	}
	q := u.Query()
	q.Set("q", code)
	u.RawQuery = q.Encode()

	log.Debug(
		"request path",
		logger.F("path", u.String()),
	)
	req, err := d.buildRequest(u)
	if err != nil {
		log.Error("fail to build request", logger.FError(err))
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("fail to make request", logger.FError(err))
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Error("fail to decode response", logger.FError(err))
		return nil, err
	}

	products, ok := data["products"].([]interface{})
	if !ok {
		log.Error("products are missing", logger.F("code", code))
		return nil, nil
	}

	if len(products) < 1 {
		log.Error("products not found", logger.F("code", code))
		return nil, nil
	}

	productJSON, err := json.Marshal(products[0])
	if err != nil {
		log.Error("fail to marshal product", logger.FError(err))
		return nil, nil
	}
	var autocompleteProduct network.Product
	if err := json.Unmarshal(productJSON, &autocompleteProduct); err != nil {
		log.Error("fail to unmarshal product", logger.FError(err))
		return nil, err
	}
	return &autocompleteProduct, nil
}

func (d *dniproClient) GetWarranty(id int64) (string, error) {
	log := d.container.GetLogger()
	fullPath := d.GetPath(WarrantyAPIEndpoint)

	u, err := url.Parse(fullPath)
	if err != nil {
		log.Error("error parsing url", logger.FError(err))
		return "", err
	}
	q := u.Query()
	q.Set("productId", fmt.Sprintf("%d", id))
	u.RawQuery = q.Encode()

	log.Debug(
		"request path",
		logger.F("path", u.String()),
	)
	req, err := d.buildRequest(u)
	if err != nil {
		log.Error("fail to build request", logger.FError(err))
		return "", err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("fail to make request", logger.FError(err))
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Error("fail to decode response", logger.FError(err))
		return "", err
	}
	warrantyRaw, ok := data["warranty"].([]interface{})
	if !ok {
		log.Error("warranty is missing")
		return "", nil
	}
	if len(warrantyRaw) < 1 {
		log.Error("warranty not found")
		return "", nil
	}

	warrantyRawItem, ok := warrantyRaw[0].(map[string]interface{})
	if !ok {
		log.Error("warranty item is missing")
		return "", nil
	}
	text, ok := warrantyRawItem["warranty"].(string)
	if !ok {
		log.Error("warranty text is missing")
		return "", nil
	}
	return text, nil
}

func (d *dniproClient) GetPath(endpoint string) string {
	return fmt.Sprintf(
		"%s%s",
		d.container.GetConfig().BaseURL,
		endpoint,
	)
}

func (d *dniproClient) buildRequest(u *url.URL) (*http.Request, error) {
	log := d.container.GetLogger()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Error("fail to create request", logger.FError(err))
		return nil, err
	}
	req.Header.Set("User-Agent", "GoClient/1.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	return req, nil
}
