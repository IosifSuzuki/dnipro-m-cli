package worker

import (
	"dniprom-cli/internal/client"
	"dniprom-cli/internal/container"
	"dniprom-cli/internal/model/app"
	"dniprom-cli/internal/model/network"
	"dniprom-cli/pkg/logger"
	"errors"
	"fmt"
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

func (w *Warranty) FetchByCode(code string) (*app.ProductWarranty, error) {
	log := w.container.GetLogger()
	const defaultMissingValue = "unknown"
	productWarranty := app.ProductWarranty{
		ID:           -1,
		Code:         code,
		Title:        defaultMissingValue,
		WarrantyText: defaultMissingValue,
		OldPrice:     defaultMissingValue,
		NewPrice:     defaultMissingValue,
	}

	productResponse, err := w.dniproClient.FetchAutocompleteProduct(code)
	if err != nil {
		log.Error("fail to fetch autocomplete product", logger.FError(err))
		return &productWarranty, err
	}

	if productResponse == nil {
		return &productWarranty, errors.New("product not found")
	}
	productWarranty.Title = GetProductName(productResponse)
	log.Debug(
		"success to fetch autocomplete product",
		logger.F("code", code),
	)

	productWarranty.ID = productResponse.ID
	warrantyText, err := w.dniproClient.GetWarranty(productResponse.ID)
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
	productWarranty.WarrantyText = warrantyText
	if oldPrice := productResponse.PriceOld.Value; oldPrice != nil {
		productWarranty.OldPrice = getFormattedPrice(*oldPrice)
	}
	if newPrice := productResponse.PriceNew.Value; newPrice != nil {
		productWarranty.NewPrice = getFormattedPrice(*newPrice)
	}
	return &productWarranty, nil
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

func getFormattedPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}
