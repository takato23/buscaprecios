package scrapers

import (
	"fmt"

	"ratoneando/cores/api"
	"ratoneando/products"
	"ratoneando/utils/logger"
)

// DiscoProduct matches the VTEX REST API response structure
type DiscoProduct struct {
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

type DiscoResponse []DiscoProduct

func Disco(query string) ([]products.Schema, error) {
	return api.Core(api.CoreProps[DiscoResponse, DiscoProduct]{
		Query:         query,
		BaseUrl:       "https://www.disco.com.ar",
		SearchPattern: func(q string) string { return "/api/catalog_system/pub/products/search/?ft=" + q },
		Source:        "disco",
		Normalizer: func(response DiscoResponse) []DiscoProduct {
			return response
		},
		Extractor: func(p DiscoProduct) products.ExtendedSchema {
			if len(p.Items) == 0 || len(p.Items[0].Sellers) == 0 {
				logger.LogWarn(fmt.Sprintf("Disco: Product %s has no items/sellers", p.ProductId))
				return products.ExtendedSchema{Unavailable: true}
			}
			
			imageUrl := ""
			if len(p.Items[0].Images) > 0 {
				imageUrl = p.Items[0].Images[0].ImageUrl
			}

			return products.ExtendedSchema{
				ID:          p.ProductId,
				Source:      "disco",
				Name:        p.ProductName,
				Link:        fmt.Sprintf("https://www.disco.com.ar/%s/p", p.LinkText),
				Image:       imageUrl,
				Unavailable: !p.Items[0].Sellers[0].CommertialOffer.IsAvailable,
				Price:       p.Items[0].Sellers[0].CommertialOffer.Price,
				ListPrice:   p.Items[0].Sellers[0].CommertialOffer.ListPrice,
			}
		},
	})
}
