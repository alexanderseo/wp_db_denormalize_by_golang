package main

import (
	"database/sql"
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

func SkFabrics(fabrics []WpFabrics, postMeta map[int64][]Postmeta) {
	var skFabrics SkTableFabrics

	for _, fabric := range fabrics {
		for key, f := range fabric {
			thisMeta := setThisMeta(key, postMeta)
			itemFabric := ItemFabric{key: {
				Id: setIdFabric(f),
				Name: setNameFabric(f),
				Slug: setSlugFabric(f),
				Hide: setHideFabric(thisMeta),
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
