package scrapers

import (
	"encoding/json"
	"fmt"

	"ratoneando/cores/api"
	"ratoneando/products"
	"ratoneando/utils/logger"
)

// Carrefour uses the same VTEX structure as Jumbo
type CarrefourResponseProduct struct {
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

type CarrefourProductData struct {
	MeasurementUnitUn string  `json:"MeasurementUnit"`
	UnitMultiplierUn  float64 `json:"UnitMultiplier"`
}

type CarrefourRawProduct struct {
	CarrefourResponseProduct
	CarrefourProductData
}

type CarrefourResponseStructure []CarrefourResponseProduct

func Carrefour(query string) ([]products.Schema, error) {
	return api.Core(api.CoreProps[CarrefourResponseStructure, CarrefourRawProduct]{
		Query:         query,
		BaseUrl:       "https://www.carrefour.com.ar",
		SearchPattern: func(q string) string { return "/api/catalog_system/pub/products/search/?ft=" + q },
		Source:        "carrefour",
		Normalizer: func(response CarrefourResponseStructure) []CarrefourRawProduct {
			var normalizedProducts []CarrefourRawProduct

			for _, rawProduct := range response {
				var productData CarrefourProductData

				if len(rawProduct.ProductData) > 0 {
					err := json.Unmarshal([]byte(rawProduct.ProductData[0]), &productData)
					if err != nil {
						logger.LogWarn(fmt.Sprintf("Error unmarshalling Carrefour product data: %s", err))
					}
				}

				normalizedProducts = append(normalizedProducts, CarrefourRawProduct{
					CarrefourResponseProduct: rawProduct,
					CarrefourProductData:     productData,
				})
			}

			return normalizedProducts
		},
		Extractor: func(rawProduct CarrefourRawProduct) products.ExtendedSchema {
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
				Source:      "carrefour",
				Name:        rawProduct.ProductName,
				Link:        fmt.Sprintf("https://www.carrefour.com.ar/%s/p", rawProduct.LinkText),
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

