package scrapers

import (
	"fmt"

	"ratoneando/cores/api"
	"ratoneando/products"
	"ratoneando/utils/logger"
)

// CarrefourProduct matches the VTEX REST API response structure
type CarrefourProduct struct {
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	LinkText    string `json:"linkText"`
	Items       []struct {
		Images []struct {
			ImageUrl string `json:"imageUrl"`
		} `json:"images"`
		Sellers []struct {
			CommertialOffer struct {
				Price       float64 `json:"Price"`
				ListPrice   float64 `json:"ListPrice"`
				IsAvailable bool    `json:"IsAvailable"`
			} `json:"commertialOffer"`
		} `json:"sellers"`
	} `json:"items"`
}

type CarrefourResponse []CarrefourProduct

func Carrefour(query string) ([]products.Schema, error) {
	return api.Core(api.CoreProps[CarrefourResponse, CarrefourProduct]{
		Query:         query,
		BaseUrl:       "https://www.carrefour.com.ar",
		SearchPattern: func(q string) string { return "/api/catalog_system/pub/products/search/?ft=" + q },
		Source:        "carrefour",
		Normalizer: func(response CarrefourResponse) []CarrefourProduct {
			return response
		},
		Extractor: func(p CarrefourProduct) products.ExtendedSchema {
			if len(p.Items) == 0 || len(p.Items[0].Sellers) == 0 {
				logger.LogWarn(fmt.Sprintf("Carrefour: Product %s has no items/sellers", p.ProductId))
				return products.ExtendedSchema{Unavailable: true}
			}
			
			imageUrl := ""
			if len(p.Items[0].Images) > 0 {
				imageUrl = p.Items[0].Images[0].ImageUrl
			}

			return products.ExtendedSchema{
				ID:          p.ProductId,
				Source:      "carrefour",
				Name:        p.ProductName,
				Link:        fmt.Sprintf("https://www.carrefour.com.ar/%s/p", p.LinkText),
				Image:       imageUrl,
				Unavailable: !p.Items[0].Sellers[0].CommertialOffer.IsAvailable,
				Price:       p.Items[0].Sellers[0].CommertialOffer.Price,
				ListPrice:   p.Items[0].Sellers[0].CommertialOffer.ListPrice,
			}
		},
	})
}
