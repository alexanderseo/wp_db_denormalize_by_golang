package main

import (
	"github.com/techleeone/gophp/serialize"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

type SkCategory struct {
	Id string `json:"id"`
	Redis_key string `json:"redis_key"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Thumbnail string `json:"thumbnail"`
	Parent_id string `json:"parent_id"`
	Recommended_categories string `json:"recommended_categories"`
	Dative_title string `json:"dative_title"`
	Nominative_title string `json:"nominative_title"`
	Has_fabric string `json:"has_fabric"`
	Attributes_product_comparison string `json:"attributes_product_comparison"`
	Attributes_filter_list string `json:"attributes_filter_list"`
	Enable_comparison string `json:"enable_comparison"`
	Detail_description string `json:"detail_description"`
}

type ItemCategory map[int64]SkCategory

type SkTableCategories []ItemCategory

func SkCategories(terms []WpTerms,
	termTaxonomy []WpTermTaxonomy,
	termMeta map[int64][]TermmetaItems,
	attachments map[int64][]sizeAttachments,
	relationships WpTableTermRelationships,
	posts []WpPosts)  {

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
									Redis_key: setRedisKey(valueItemTerm.Term_id, relationships, posts),
									Name: setName(valueItemTerm),
									Slug: setSlug(valueItemTerm),
									Thumbnail: setThumbnail(thisTermMeta, attachments),
									Parent_id: setParentId(valueTaxonomy),
									Recommended_categories: setRecommendedCategories(thisTermMeta),
									Dative_title: setDativeTitle(thisTermMeta),
									Nominative_title: setNominativeTitle(thisTermMeta),
									Has_fabric: setHasFabric(thisTermMeta),
									Attributes_product_comparison: setAttributesProductComparison(thisTermMeta),
									Attributes_filter_list: setAttributesFilterList(thisTermMeta),
									Enable_comparison: setEnableComparison(thisTermMeta),
									Detail_description: setDetailDescription(thisTermMeta),
								},
							}
							//fmt.Println("================", categoryItems)
							categoriesData = append(categoriesData, categoryItems)
						}
					}
				}
			}
		}
	}

	//fmt.Println(categoriesData)
}

func setThisTermMeta(id int64, termMeta map[int64][]TermmetaItems) []TermmetaItems {
	return termMeta[id]
}

func setId(term Term) string {
	return term.Term_id
}

func setRedisKey(id string, relationships WpTableTermRelationships, posts []WpPosts) string {
	var ids []string

	for _, relations := range relationships {
		for _, term := range relations {
			if term.Term_taxonomy_id == id {
				ids = append(ids, term.Object_id)
			}
		}
	}

	var products []string
	dataLayout := "2006-01-02 15:04:05"
	for _, i := range ids {
		for _, post := range posts {
			for _, p := range post {
				if i == p.Id {
					var sb strings.Builder
					tTime, _ := time.Parse(dataLayout, p.Post_modified)
					sb.WriteString(p.Id)
					sb.WriteString("_")
					strData := strconv.Itoa(int(tTime.Unix()))
					sb.WriteString(strData)
					products = append(products, sb.String())
					sb.Reset()
				}
			}
		}
	}


	productsString := strings.Join(products[:], "_")
	algorithm := fnv.New64a()
	algorithm.Write([]byte(productsString))
	hash := algorithm.Sum64()

	strHash := strconv.FormatUint(hash, 10)
	return strHash
}

func setName(term Term) string {
	return term.Name
}

func setSlug(term Term) string {
	return term.Slug
}

func setThumbnail(termMeta []TermmetaItems, attachments map[int64][]sizeAttachments) string {
	var thumbnails []sizeAttachments
	var thumbnail string

	for _, v := range termMeta {
		if val, ok := v["thumbnail_id"]; ok {
			i, _ := strconv.ParseInt(val, 10, 64)
			thumbnails = attachments[i]
		}
	}

	thumbnailsBytes, _ := serialize.Marshal(thumbnails)
	thumbnail = string(thumbnailsBytes[:])

	return thumbnail
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

func setEnableComparison(termMeta []TermmetaItems) string {
	enable_product_comparison := "0"

	for _, v := range termMeta {
		if val, ok := v["enable-product-comparison"]; ok {
			enable_product_comparison = val
		}
	}

	return enable_product_comparison
}

func setDetailDescription(termMeta []TermmetaItems) string {
	search_meta_keys := [3]string{"details_0_details_row", "details_1_details_row", "details_2_details_row"}
	type description map[string]string
	descriptions := make(map[int][]description)

	for key, meta := range search_meta_keys {
		for _, v := range termMeta {
			if val, ok := v[meta]; ok {
				tempDescription := description{meta: val}
				descriptions[key] = append(descriptions[key], tempDescription)

			}
		}
	}

	descriptionsBytes, _ := serialize.Marshal(descriptions)
	descriptionsString := string(descriptionsBytes[:])

	return descriptionsString
}