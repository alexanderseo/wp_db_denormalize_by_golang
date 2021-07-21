package main

import (
	"database/sql"
	"fmt"
	"github.com/techleeone/gophp/serialize"
)

type SkFabric struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Hide string `json:"hide"`
	Collection string `json:"collection"`
	Color string `json:"color"`
	Image string `json:"image"`
	Gallery string `json:"gallery"`
	Video string `json:"video"`
	Category string `json:"category"`
	Properties string `json:"properties"`
	Description string `json:"description"`
	Material string `json:"material"`
}

type ItemFabric map[int64]SkFabric

type SkTableFabrics []ItemFabric

func SkFabrics(fabrics []WpFabrics, postMeta map[int64][]Postmeta, terms []WpTerms, termMeta map[int64][]TermmetaItems, attachments map[int64][]sizeAttachments) {
	var skFabrics SkTableFabrics

	for _, fabric := range fabrics {
		for key, f := range fabric {
			thisMeta := setThisMeta(key, postMeta)
			itemFabric := ItemFabric{key: {
				Id: setIdFabric(f),
				Name: setNameFabric(f),
				Slug: setSlugFabric(f),
				Hide: setHideFabric(thisMeta),
				Collection: setCollectionFabric(thisMeta, terms, termMeta, attachments),
			}}

			skFabrics = append(skFabrics, itemFabric)
		}
	}
	//fmt.Println("--------", skFabrics)
}

func setThisMeta(id int64, postMeta map[int64][]Postmeta) []Postmeta {
	return postMeta[id]
}

func setIdFabric(fabric Post) string {
	return fabric.Id
}

func setNameFabric(fabric Post) string {
	return fabric.Post_title
}

func setSlugFabric(fabric Post) string {
	return fabric.Post_name
}

func setHideFabric(thisMeta []Postmeta) string {
	var hide sql.NullString

	for _, meta := range thisMeta {
		if meta.Meta_key == "hide" {
			hide = meta.Meta_value
		}
	}

	return hide.String
}

func setCollectionFabric(thisMeta []Postmeta, terms []WpTerms,termMeta map[int64][]TermmetaItems, attachments map[int64][]sizeAttachments) string {
	var collection sql.NullString
	type CollectionData struct {
		Id string `json:"id"`
		Slug string `json:"slug"`
		Name string `json:"name"`
		Gallery string `json:"gallery"`
		Care_advice string `json:"care_advice"`
	}

	for _, meta := range thisMeta {
		if meta.Meta_key == "collection" {
			collection = meta.Meta_value
		}
	}

	collectionId := collection.String

	setSlugNameCollection := func (id string, terms []WpTerms) (slug, name string) {
		for _, v := range terms {
			idInt64 := strToInt(id)
			if val, ok := v[idInt64]; ok {
				slug = val.Slug
				name = val.Name
			}
		}

		return
	}

	slug, name := setSlugNameCollection(collectionId, terms)

	setGalleryCollection := func (id string, termMeta map[int64][]TermmetaItems, attachments map[int64][]sizeAttachments) (gallery, careAdvice string) {
		idInt64 := strToInt(id)
		var galleryData, careData []sizeAttachments
		if term, ok := termMeta[idInt64]; ok {
			for _, termMap := range term {

				if galleryString, ok := termMap["gallery"]; ok {
					galleryArray, _ := serialize.UnMarshal([]byte(galleryString))
					if galleryArray != nil {
						var Ids []string
						if rec, ok := galleryArray.([]interface{}); ok {
							for _, val := range rec {
								Ids = append(Ids, val.(string))
							}
						}

						if Ids != nil {
							for _, idImg := range Ids {
								attachId := strToInt(idImg)
								if g, ok := attachments[attachId]; ok {
									galleryData = g
								}
							}
						}
					}
				}

				if careString, ok := termMap["care-advice"]; ok {
					if careString != "" {
						careId := strToInt(careString)
						if c, ok := attachments[careId]; ok {
							careData = c
						}
					}
				}
			}

		}

		galleryBytes, _ := serialize.Marshal(galleryData)
		gallery = string(galleryBytes[:])
		careBytes, _ := serialize.Marshal(careData)
		careAdvice = string(careBytes[:])
		return
	}
	gallery, careAdvice := setGalleryCollection(collectionId, termMeta, attachments)

	collectionData := CollectionData{
		Id: collectionId,
		Slug: slug,
		Name: name,
		Gallery: gallery,
		Care_advice: careAdvice,
	}

	fmt.Println("---------", collectionData)
	return collectionId
}


