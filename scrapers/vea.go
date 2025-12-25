package scrapers

import (
	"encoding/json"
	"fmt"

	"ratoneando/cores/api"
	"ratoneando/products"
	"ratoneando/utils/logger"
)

// Vea uses the same VTEX structure as Jumbo/Carrefour/Disco
type VeaResponseProduct struct {
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

type VeaProductData struct {
	MeasurementUnitUn string  `json:"MeasurementUnit"`
	UnitMultiplierUn  float64 `json:"UnitMultiplier"`
}

type VeaRawProduct struct {
	VeaResponseProduct
	VeaProductData
}

type VeaResponseStructure []VeaResponseProduct

func Vea(query string) ([]products.Schema, error) {
	return api.Core(api.CoreProps[VeaResponseStructure, VeaRawProduct]{
		Query:         query,
		BaseUrl:       "https://www.vea.com.ar",
		SearchPattern: func(q string) string { return "/api/catalog_system/pub/products/search/?ft=" + q },
		Source:        "vea",
		Normalizer: func(response VeaResponseStructure) []VeaRawProduct {
			var normalizedProducts []VeaRawProduct

			for _, rawProduct := range response {
				var productData VeaProductData

				if len(rawProduct.ProductData) > 0 {
					err := json.Unmarshal([]byte(rawProduct.ProductData[0]), &productData)
					if err != nil {
						logger.LogWarn(fmt.Sprintf("Error unmarshalling Vea product data: %s", err))
					}
				}

				normalizedProducts = append(normalizedProducts, VeaRawProduct{
					VeaResponseProduct: rawProduct,
					VeaProductData:     productData,
				})
			}

			return normalizedProducts
		},
		Extractor: func(rawProduct VeaRawProduct) products.ExtendedSchema {
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
				Source:      "vea",
				Name:        rawProduct.ProductName,
				Link:        fmt.Sprintf("https://www.vea.com.ar/%s/p", rawProduct.LinkText),
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
