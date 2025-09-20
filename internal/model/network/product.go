package network

type Product struct {
	ID   int64 `json:"id"`
	Name struct {
		RU string `json:"ru"`
		UK string `json:"uk"`
		EN string `json:"en"`
	} `json:"name"`
}
