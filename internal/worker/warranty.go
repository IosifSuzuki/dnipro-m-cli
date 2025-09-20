package worker

import (
	"dniprom-cli/internal/client"
	"dniprom-cli/internal/container"
	"dniprom-cli/internal/model/app"
	"dniprom-cli/internal/model/network"
	"dniprom-cli/pkg/logger"
	"errors"
)

type Warranty struct {
	container    container.Container
	dniproClient client.DniproClient
}

func NewWarrantyWorker(container container.Container, dniproClient client.DniproClient) *Warranty {
	return &Warranty{
		container:    container,
		dniproClient: dniproClient,
	}
}

func (w *Warranty) FetchByCode(code string) (*app.Warranty, error) {
	log := w.container.GetLogger()
	const defaultMissingValue = "unknown"
	warranty := app.Warranty{
		ID:           -1,
		Code:         code,
		Title:        defaultMissingValue,
		WarrantyText: defaultMissingValue,
	}

	product, err := w.dniproClient.FetchAutocompleteProduct(code)
	if err != nil {
		log.Error("fail to fetch autocomplete product", logger.FError(err))
		return &warranty, err
	}

	if product == nil {
		return &warranty, errors.New("product not found")
	}
	warranty.Title = GetProductName(product)
	log.Debug(
		"success to fetch autocomplete product",
		logger.F("code", code),
	)

	warranty.ID = product.ID
	warrantyText, err := w.dniproClient.GetWarranty(product.ID)
	if err != nil {
		log.Error(
			"fail to fetch warranty's product",
			logger.F("code", code),
			logger.FError(err),
		)
	}
	if warrantyText == "" {
		warrantyText = defaultMissingValue
	}
	warranty.WarrantyText = warrantyText
	return &warranty, nil
}

func GetProductName(product *network.Product) string {
	const defaultProductTitle = "unknown"
	if product == nil {
		return defaultProductTitle
	}
	if product.Name.UK != "" {
		return product.Name.UK
	} else if product.Name.RU != "" {
		return product.Name.RU
	} else if product.Name.EN != "" {
		return product.Name.EN
	}
	return defaultProductTitle
}
