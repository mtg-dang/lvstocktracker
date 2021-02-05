package lvapi

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// A RegionURL represents a region object containing the relevant data
// for a Louis Vuitton region identifier.
type RegionURL struct {
	code string // Region Code
	url  string // Main landing page URL for specified region code
}

// A CategoryURL represents a subcategory object containing a subcategory name
// and the corresponding route to the subcategory page.
type CategoryURL struct {
	name string // Subcategory name
	url  string // Subcategory URL/route
}

// A ProductRoute represents a product page url and its corresponding product name.
type ProductRoute struct {
	name  string // Product name
	route string // Product route
}

// A ProductImage represents a product name with its matching product image url.
type ProductImage struct {
	name string // Product name
	url  string // Product image URL
}

// A ProductAvailability represents a product identifier sku with its online availability
type ProductAvailability struct {
	Sku       string `json:"Sku"`       // Product identifier
	Available bool   `json:"Available"` // Product availability
}

// createCollyCollector creates a new gocolly collector and assigns random user agent.
// It returns the created gocolly collector.
func createCollyCollector() *colly.Collector {
	// Init colly collector
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	// Random UA on each access to prevent blacklisting
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	return c
}

// GetLVRegionCodesAndURLs sends a request to the Louis Vuitton landing page for crawling.
// It crawls the page for region URLs.
// It returns a map keyed by region codes with corresponding URL value for region.
func GetLVRegionCodesAndURLs() []RegionURL {
	// Array for holding RegionURL structs which contain region code and corresponding url
	var regionCodesAndURLs []RegionURL
	// Init colly collector
	c := createCollyCollector()
	// Find within the children of each li tag anything with class .lvdispatch-link
	// Use the value of that link as well as the region code contained within the link
	// to create a a RegionURL struct. Then append the created struct to a slice for return.
	c.OnHTML("li", func(e *colly.HTMLElement) {
		pageDom := e.DOM
		pageDom.Find(".lvdispatch-link").Each(func(i int, s *goquery.Selection) {
			link, linkExists := s.Attr("href")
			if linkExists {
				regionCodeAndURL := RegionURL{code: strings.Split(link, "/")[3], url: link}
				regionCodesAndURLs = append(regionCodesAndURLs, regionCodeAndURL)
			}
		})
	})
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response Handler
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Body)
	})
	// Send visit request to colly collector
	c.Visit("https://www.louisvuitton.com/dispatch/?noDRP=true")

	return regionCodesAndURLs
}

// GetLVMainCategories sends a request to url for crawling.
// It crawls url for the main nav bar item names.
// It returns a slice of strings which are the category names.
func GetLVMainCategories(url string) []string {
	// Slice to hold category names
	var mainCategories []string
	// Init colly collector
	c := createCollyCollector()
	// Find within the children of each li tag anything with class .lv-header-main-nav__item
	// Append the text of that span to categories array
	c.OnHTML("li", func(e *colly.HTMLElement) {
		pageDom := e.DOM
		pageDom.Find(".lv-header-main-nav__item").Each(func(i int, s *goquery.Selection) {
			mainCategories = append(mainCategories, s.Find("span").Text())
		})
	})
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response Handler
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Body)
	})
	// Send visit request to colly collector
	c.Visit(url)
	return mainCategories
}

// GetLVSubCategoriesRoutes sends a request to url for crawling.
// It crawls url based on mainCategory for subcateogries within the nav.
// It returns a map containing the route of all subcategories under mainCategory with corresponding label as the key.
func GetLVSubCategoriesRoutes(mainCategory string, url string) []CategoryURL {
	// Slice to hold subcategory structs
	var subCategories []CategoryURL
	// Init colly collector
	c := createCollyCollector()
	// Find within the children of each li tag anything with class .lv-header-main-nav__item.
	// Finds the span within the found class that matches mainCategory.
	// Finds the corresponding subcategories within the parent to find all routes and subcategory labels.
	// Creates a struct for each subcategory label/route into subCategories slice
	c.OnHTML("li[role=presentation]", func(e *colly.HTMLElement) {
		pageDom := e.DOM
		pageDom.Find(".lv-header-main-nav__item").Each(func(i int, s *goquery.Selection) {
			if s.Find("span").Text() == mainCategory {
				s.Parent().
					Find(".lv-header-main-nav-panel").
					Find(".lv-header-main-nav-child__item").
					Each(func(i int, s *goquery.Selection) {
						href, hrefExists := s.Find(".lv-header-main-nav-child__link").Attr("href")
						if hrefExists {
							subCategory := CategoryURL{name: s.Find(".lv-header-main-nav-child__link").Text(), url: href}
							subCategories = append(subCategories, subCategory)
						}
					})
			}
		})
	})
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response Handler
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Body)
	})
	// Send visit request to colly collector
	c.Visit(url)
	return subCategories
}

// GetLVProductPageRoutes sends a request to url for crawling.
// It crawls for each product contained in url which is the subcategory url.
// It returns a slice of ProductRoute structures each containing the product name
// and product route.
func GetLVProductPageRoutes(url string) []ProductRoute {
	// Slice containing ProductRoute objects
	// Each product route is obtained from a subcategory page
	var productPages []ProductRoute
	// Init colly collector
	c := createCollyCollector()
	// Find within the children of ul tag with class lv-list.
	// Finds each .lv-product-card within the list
	// Creates a ProductRoute structure using the product name and the product href into productPages.
	c.OnHTML("ul[class=lv-list]", func(e *colly.HTMLElement) {
		pageDom := e.DOM
		pageDom.Find(".lv-product-card").Each(func(i int, s *goquery.Selection) {
			productRoute, productRouteExists := s.Attr("href")
			if productRouteExists {
				p := strings.NewReader(s.Text())
				productText, _ := goquery.NewDocumentFromReader(p)
				productText.Find("img").Each(func(i int, el *goquery.Selection) {
					//productImageSrc, productImageSrcExists := el.Attr("src")
					el.Remove()
					productRouteStruct := ProductRoute{name: strings.TrimSpace(productText.Text()), route: productRoute}
					productPages = append(productPages, productRouteStruct)
				})
			}
		})
	})
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response Handler
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Body)
	})
	// Send visit request to colly collector
	c.Visit(url)
	return productPages
}

// GetLVProductImages sends a request to url for crawling.
// It crawls for each product contained in url which is the subcategory url.
// It returns a slice of ProductImage structures which contain the product name and product image url.
func GetLVProductImages(url string) []ProductImage {
	// Slice to hold ProductImage structs
	var productImages []ProductImage
	// Init colly collector
	c := createCollyCollector()
	// Find within the children of ul tag with class lv-list.
	// Finds each .lv-product-card within the list
	// Creates a ProductImage using the product name and the product image url into productPages productImages.
	c.OnHTML("ul[class=lv-list]", func(e *colly.HTMLElement) {
		pageDom := e.DOM
		pageDom.Find(".lv-product-card").Each(func(i int, s *goquery.Selection) {
			p := strings.NewReader(s.Text())
			productText, _ := goquery.NewDocumentFromReader(p)
			productText.Find("img").Each(func(i int, el *goquery.Selection) {
				productImageSrc, productImageSrcExists := el.Attr("src")
				el.Remove()
				if productImageSrcExists {
					productImage := ProductImage{name: strings.TrimSpace(productText.Text()), url: productImageSrc}
					productImages = append(productImages, productImage)
				}
			})
		})
	})
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response Handler
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Body)
	})
	// Send visit request to colly collector
	c.Visit(url)
	return productImages
}

// getLVProductJSONBodyBySKU sends a request to https://api.louisvuitton.com/api/eng-ca/catalog/skus/{sku}
// It crawls the REST API endpoint and returns the response body.
// The response body should be a JSON string if the REST API was successfully loaded.
// gocolly is used to extract the JSON from the REST API endpoint as access via HTTP requests is denied.
// gocolly allows us to access the end point by randomizing our user agent.
func getLVProductJSONBodyBySKU(sku string) string {
	// REST API endpoint for LV SKU catalog
	endpoint := "https://api.louisvuitton.com/api/eng-ca/catalog/skus/" + sku
	// JSON output
	jsonString := ""
	// Init colly collector
	c := createCollyCollector()
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response body contains the JSON string from API endpoint.
	// Extract the JSON string for return
	c.OnResponse(func(r *colly.Response) {
		jsonString = string(r.Body)
	})
	// Send visit request to colly collector
	c.Visit(endpoint)
	return jsonString
}

// GetLVProductPageURLBySKU sends a request to getLVProductJSONBodyBySKU.
// It retrieves a JSON string for the corresponding sku.
// The JSON string is parsed and the product page URL is returned.
func GetLVProductPageURLBySKU(sku string) string {
	// Output URL
	url := ""
	// Call to retrieve JSON string from REST API endpoint
	jsonString := getLVProductJSONBodyBySKU(sku)
	// Checks if the JSON response has a list size of greater than zero to
	// ensure that SKU is valid. If valid, then the JSON string is parsed
	// and the product page URL is extracted from the JSON string.
	// If invalid, then an error is returned.
	if strings.Contains(jsonString, "\"skuListSize\":0") {
		url = "Invalid SKU"
	} else {
		var result map[string]interface{}
		json.Unmarshal([]byte(jsonString), &result)
		for _, item := range result["skuList"].([]interface{}) {
			url = fmt.Sprintf("%v", item.(map[string]interface{})["url"])
		}
	}
	return url
}

// getNextMapLevelDown checks rec for key, and returns the value of key in rec.
// Returns a reflect.Zero value if rec is not a Map or if key is not found
func getNextMapLevelDown(rec reflect.Value, key string) reflect.Value {
	if rec.Kind() == reflect.Map {
		for _, k := range rec.MapKeys() {
			v := rec.MapIndex(k)
			if k.Interface() == key {
				return reflect.ValueOf(v.Interface())
			}
		}
	}
	return reflect.Zero(reflect.TypeOf(0))
}

// GetLVProductPageAPIEndPointBySKU sends a request to getLVProductJSONBodyBySKU.
// It retrieves a JSON string for the corresponding sku.
// The JSON string is parsed and the product API endpoint is returned.
func GetLVProductPageAPIEndPointBySKU(sku string) string {
	// Output endpoint
	endpoint := ""
	// Call to retrieve JSON string from REST API endpoint
	jsonString := getLVProductJSONBodyBySKU(sku)
	// Checks if the JSON response has a list size of greater than zero to
	// ensure that SKU is valid. If valid, then the JSON string is parsed
	// and the product page API endpoint is extracted from the JSON string.
	// If invalid, then an error is returned.
	if strings.Contains(jsonString, "\"skuListSize\":0") {
		endpoint = "Invalid SKU"
	} else {
		var result map[string]interface{}
		json.Unmarshal([]byte(jsonString), &result)
		for _, item := range result["skuList"].([]interface{}) {
			if rec, ok := item.(map[string]interface{}); ok {
				for key, val := range rec {
					if key == "_links" {
						v := reflect.ValueOf(val)
						selfMap := getNextMapLevelDown(v, "self")
						if selfMap.IsValid() {
							hrefMap := getNextMapLevelDown(selfMap, "href")
							endpoint = fmt.Sprintf("%v", hrefMap)
						}
					}
				}
			}
		}
	}
	return endpoint
}

// GetLVProductAvailabilityBySKU sends a request to the product API page for sku:
// 		'https://api.louisvuitton.com/api/eng-ca/catalog/product/sku'
// It crawls and retrieves the JSON string from the endpoint.
// The JSON string is parsed into a map, and then proccessed to extract availability
// for sku based on the value of backOrderDisclaimer for the sku.
// It returns true if the product sku is available, false if not.
// gocolly is used to extract the JSON from the REST API endpoint as access via HTTP requests is denied.
// gocolly allows us to access the end point by randomizing our user agent.
func GetLVProductAvailabilityBySKU(sku string) ProductAvailability {
	// REST API endpoint for LV SKU catalog
	endpoint := "https://api.louisvuitton.com/api/eng-ca/catalog/product/" + sku
	isProductAvailable := false
	// Init colly collector
	c := createCollyCollector()
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response body contains the JSON string from API endpoint.
	// Parse the JSON string from the response to extract the backOrderDisclaimer.
	// There may be multiple backOrderDisclaimer fields depending on if there is
	// related skus to the search sku. Match backOrderDisclaimer using identifier.
	c.OnResponse(func(r *colly.Response) {
		jsonString := string(r.Body)
		if strings.Contains(jsonString, "errorCode") {
			isProductAvailable = false
		} else {
			var result map[string]interface{}
			json.Unmarshal([]byte(jsonString), &result)
			for _, item := range result["model"].([]interface{}) {
				if item.(map[string]interface{})["identifier"] == sku {
					propertyMapSlice := reflect.ValueOf(item.(map[string]interface{})["additionalProperty"])
					if propertyMapSlice.Kind() == reflect.Slice {
						for i := 0; i < propertyMapSlice.Len(); i++ {
							propertyMap := reflect.ValueOf(propertyMapSlice.Index(i))
							if strings.Contains(fmt.Sprintf("%v", propertyMap.Interface()), "name:backOrderDisclaimer value:false") {
								isProductAvailable = true
								break
							}
						}
					}
				}
			}
		}
	})
	// Send visit request to colly collector
	c.Visit(endpoint)
	return ProductAvailability{Sku: sku, Available: isProductAvailable}
}

// GetLVAlternativeStyleProductIndentifierAndAvailabilityForSKU sends a request to the product API page for sku:
// 		'https://api.louisvuitton.com/api/eng-ca/catalog/product/sku'
// It crawls and retrieves the JSON string from the endpoint.
// The JSON string is parsed into a map, and then proccessed to extract availability
// for sku based on the value of backOrderDisclaimer for the sku. Then searches the rest of the JSON,
// if there is alternative styles for their backOrderDisclaimer and their sku.
// It returns a slice of structs each containing a sku number, and the availability.
// gocolly is used to extract the JSON from the REST API endpoint as access via HTTP requests is denied.
// gocolly allows us to access the end point by randomizing our user agent.
func GetLVAlternativeStyleProductIndentifierAndAvailabilityForSKU(sku string) []ProductAvailability {
	// Output slice
	var productAvailabilitySlice []ProductAvailability
	// REST API endpoint for LV SKU catalog
	endpoint := "https://api.louisvuitton.com/api/eng-ca/catalog/product/" + sku
	// Init colly collector
	c := createCollyCollector()
	// Request Handler
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Response body contains the JSON string from API endpoint.
	// Parse the JSON string from the response to extract the backOrderDisclaimer.
	// There may be multiple backOrderDisclaimer fields depending on if there is
	// related skus to the search sku. Create a ProductAvailability struct for each product
	c.OnResponse(func(r *colly.Response) {
		jsonString := string(r.Body)
		if strings.Contains(jsonString, "errorCode") {
			productAvailabilitySlice = nil
		} else {
			var result map[string]interface{}
			json.Unmarshal([]byte(jsonString), &result)
			for _, item := range result["model"].([]interface{}) {
				identifier := item.(map[string]interface{})["identifier"]
				if identifier != nil {
					propertyMapSlice := reflect.ValueOf(item.(map[string]interface{})["additionalProperty"])
					if propertyMapSlice.Kind() == reflect.Slice {
						for i := 0; i < propertyMapSlice.Len(); i++ {
							propertyMap := reflect.ValueOf(propertyMapSlice.Index(i))
							if strings.Contains(fmt.Sprintf("%v", propertyMap.Interface()), "name:backOrderDisclaimer value:false") {
								product := ProductAvailability{Sku: fmt.Sprintf("%v", identifier), Available: true}
								productAvailabilitySlice = append(productAvailabilitySlice, product)
								break
							} else if strings.Contains(fmt.Sprintf("%v", propertyMap.Interface()), "name:backOrderDisclaimer value:true") {
								product := ProductAvailability{Sku: fmt.Sprintf("%v", identifier), Available: false}
								productAvailabilitySlice = append(productAvailabilitySlice, product)
								break
							}
						}
					}
				}
			}
		}
	})
	// Send visit request to colly collector
	c.Visit(endpoint)
	return productAvailabilitySlice
}
