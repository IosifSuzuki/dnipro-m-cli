package network

import "dniprom-cli/pkg/jsonx"

type Product struct {
	ID   int64 `json:"id"`
	Name struct {
		RU string `json:"ru"`
		UK string `json:"uk"`
		EN string `json:"en"`
	} `json:"name"`
	PriceNew jsonx.NullableFloat64 `json:"price_new"`
	PriceOld jsonx.NullableFloat64 `json:"price_old"`
}
