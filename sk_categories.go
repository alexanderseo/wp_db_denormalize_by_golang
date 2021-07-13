package main

import (
	"fmt"
	"github.com/techleeone/gophp/serialize"
	"strings"
)

type SkCategory struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Parent_id string `json:"parent_id"`
	Recommended_categories string `json:"recommended_categories"`
	Dative_title string `json:"dative_title"`
	Nominative_title string `json:"nominative_title"`
	Has_fabric string `json:"has_fabric"`
	Attributes_product_comparison string `json:"attributes_product_comparison"`
	Attributes_filter_list string `json:"attributes_filter_list"`
}

type ItemCategory map[int64]SkCategory

type SkTableCategories []ItemCategory

func SkCategories(terms []WpTableTerms, termTaxonomy []WpTableTermTaxonomy, termMeta map[int64][]TermmetaItems)  {
	var categoriesData SkTableCategories

	for _, itemTaxonomy := range termTaxonomy{
		for keyTaxonomy, valueTaxonomy := range itemTaxonomy {
			if valueTaxonomy.Taxonomy == "product_cat" {
				for _, itemTerm := range terms{
					for keyItemTerm, valueItemTerm := range itemTerm{
						if keyTaxonomy == keyItemTerm {

							thisTermMeta := setThisTermMeta(keyTaxonomy, termMeta)

							categoryItems := ItemCategory {
								keyTaxonomy: SkCategory{
									Id: setId(valueItemTerm),
									Name: setName(valueItemTerm),
									Slug: setSlug(valueItemTerm),
									Parent_id: setParentId(valueTaxonomy),
									Recommended_categories: setRecommendedCategories(thisTermMeta),
									Dative_title: setDativeTitle(thisTermMeta),
									Nominative_title: setNominativeTitle(thisTermMeta),
									Has_fabric: setHasFabric(thisTermMeta),
									Attributes_product_comparison: setAttributesProductComparison(thisTermMeta),
									Attributes_filter_list: setAttributesFilterList(thisTermMeta),
								},
							}

							categoriesData = append(categoriesData, categoryItems)
						}
					}
				}
			}
		}
	}

	fmt.Println(categoriesData)
}

func setThisTermMeta(id int64, termMeta map[int64][]TermmetaItems) []TermmetaItems {
	var thisTermMeta []TermmetaItems

	for key, value := range termMeta {
		if key == id {
			thisTermMeta = value
			break
		}
	}

	return thisTermMeta
}

func setId(term Terms) string {
	return term.Term_id
}

func setName(term Terms) string {
	return term.Name
}

func setSlug(term Terms) string {
	return term.Slug
}

func setParentId(taxonomies TermTaxonomy) string {
	return taxonomies.Parent
}

func setRecommendedCategories(termMeta []TermmetaItems) string {
	recommended_categories := "0"

	for _, v := range termMeta {
		if val, ok := v["recommended_categories"]; ok {
			recommended_categories = val
		}
	}

	out, _ := serialize.UnMarshal([]byte(recommended_categories))

	if out != nil {
		var Ids []string
		if rec, ok := out.([]interface{}); ok {
			for _, val := range rec {
				Ids = append(Ids, val.(string))
			}
		}

		recommended_categories = strings.Join(Ids, ", ")
	}

	return recommended_categories
}

func setDativeTitle(termMeta []TermmetaItems) string {
	dative_title := "0"

	for _, v := range termMeta {
		if val, ok := v["dative-title"]; ok {
			dative_title = val
		}
	}

	return dative_title
}

func setNominativeTitle(termMeta []TermmetaItems) string {
	nominative_title := "0"

	for _, v := range termMeta {
		if val, ok := v["nominative-title"]; ok {
			nominative_title = val
		}
	}

	return nominative_title
}

func setHasFabric(termMeta []TermmetaItems) string {
	has_fabric := "0"

	for _, v := range termMeta {
		if val, ok := v["has-fabric"]; ok {
			has_fabric = val
		}
	}

	return has_fabric
}

func setAttributesProductComparison(termMeta []TermmetaItems) string {
	attributes_product_comparison := "0"

	for _, v := range termMeta {
		if val, ok := v["attributes-product-comparison"]; ok {
			attributes_product_comparison = val
		}
	}

	return attributes_product_comparison
}

func setAttributesFilterList(termMeta []TermmetaItems) string {
	var attributes_filter_list string
	var switchExpressions string
	var indexArray string
	type data map[string]string
	dataList := make(map[string][]data)
	var nameFilter string

	for _, v := range termMeta {
		for key, value := range v {
			if strings.Contains(key, "attributes-filter-list_") {
				words := strings.Split(key, "_")
				if words[0] != "" {
					indexArray = words[1]
					switchExpressions = words[2]
				} else {
					indexArray = words[2]
					switchExpressions = words[3]
				}

				if !strings.Contains(value, "field_") {
					nameFilter = value

					switch switchExpressions {
					case "attribute-item":
						inx := indexArray
						dataItem := data{"item": nameFilter}
						dataList[inx] = append(dataList[inx], dataItem)

					case "attribute-type":
						inx := indexArray
						dataItem := data{"type": nameFilter}
						dataList[inx] = append(dataList[inx], dataItem)

					}
				}
			}
		}
	}
	type reformattedDataList map[string]string
	var reformattedData []reformattedDataList
	if len(dataList) != 0 {
		for _, value := range dataList {
			for _, item := range value {
				itemData := reformattedDataList{item["item"]: item["type"]}
				reformattedData = append(reformattedData, itemData)
			}
		}
	}
	if len(reformattedData) != 0 {
		filterBytes, _ := serialize.Marshal(reformattedData)
		attributes_filter_list = string(filterBytes[:])
	} else {
		filterBytes, _ := serialize.Marshal("0")
		attributes_filter_list = string(filterBytes[:])
	}

	return attributes_filter_list
}