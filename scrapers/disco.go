package scrapers

import (
	"encoding/json"
	"fmt"

	"ratoneando/cores/api"
	"ratoneando/products"
	"ratoneando/utils/logger"
)

// Disco uses the same VTEX structure as Jumbo/Carrefour
type DiscoResponseProduct struct {
	ProductId   string   `json:"productId"`
	ProductName string   `json:"productName"`
	Link        string   `json:"link"`
	LinkText    string   `json:"linkText"`
	ProductData []string `json:"ProductData"`
	Items       []struct {
		Images []struct {
			ImageUrl string `json:"imageUrl"`
		} `json:"images"`
		Sellers []struct {
			CommertialOffer struct {
				Price                float64 `json:"Price"`
				ListPrice            float64 `json:"ListPrice"`
				PriceWithoutDiscount float64 `json:"PriceWithoutDiscount"`
				AvailableQuantity    int     `json:"AvailableQuantity"`
				IsAvailable          bool    `json:"IsAvailable"`
			} `json:"commertialOffer"`
		} `json:"sellers"`
	} `json:"items"`
}

type DiscoProductData struct {
	MeasurementUnitUn string  `json:"MeasurementUnit"`
	UnitMultiplierUn  float64 `json:"UnitMultiplier"`
}

type DiscoRawProduct struct {
	DiscoResponseProduct
	DiscoProductData
}

type DiscoResponseStructure []DiscoResponseProduct

func Disco(query string) ([]products.Schema, error) {
	return api.Core(api.CoreProps[DiscoResponseStructure, DiscoRawProduct]{
		Query:         query,
		BaseUrl:       "https://www.disco.com.ar",
		SearchPattern: func(q string) string { return "/api/catalog_system/pub/products/search/?ft=" + q },
		Source:        "disco",
		Normalizer: func(response DiscoResponseStructure) []DiscoRawProduct {
			var normalizedProducts []DiscoRawProduct

			for _, rawProduct := range response {
				var productData DiscoProductData

				if len(rawProduct.ProductData) > 0 {
					err := json.Unmarshal([]byte(rawProduct.ProductData[0]), &productData)
					if err != nil {
						logger.LogWarn(fmt.Sprintf("Error unmarshalling Disco product data: %s", err))
					}
				}

				normalizedProducts = append(normalizedProducts, DiscoRawProduct{
					DiscoResponseProduct: rawProduct,
					DiscoProductData:     productData,
				})
			}

			return normalizedProducts
		},
		Extractor: func(rawProduct DiscoRawProduct) products.ExtendedSchema {
			// Bounds checking
			if len(rawProduct.Items) == 0 {
				return products.ExtendedSchema{}
			}
			if len(rawProduct.Items[0].Images) == 0 {
				return products.ExtendedSchema{}
			}
			if len(rawProduct.Items[0].Sellers) == 0 {
				return products.ExtendedSchema{}
			}

			return products.ExtendedSchema{
				ID:          rawProduct.ProductId,
				Source:      "disco",
				Name:        rawProduct.ProductName,
				Link:        fmt.Sprintf("https://www.disco.com.ar/%s/p", rawProduct.LinkText),
				Image:       rawProduct.Items[0].Images[0].ImageUrl,
				Unavailable: !rawProduct.Items[0].Sellers[0].CommertialOffer.IsAvailable,
				Price:       rawProduct.Items[0].Sellers[0].CommertialOffer.Price,
				ListPrice:   rawProduct.Items[0].Sellers[0].CommertialOffer.ListPrice,
				Unit:        rawProduct.MeasurementUnitUn,
				UnitFactor:  rawProduct.UnitMultiplierUn,
			}
		},
	})
}
