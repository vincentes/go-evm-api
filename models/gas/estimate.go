package models

type GasEstimateResponse struct {
	Slow    string `json:"slow"`
	Average string `json:"average"`
	Fast    string `json:"fast"`
}
