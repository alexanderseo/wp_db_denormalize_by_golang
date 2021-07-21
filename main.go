package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	start := time.Now()


	terms := GetTerms()
	term_taxonomy := GetTermTaxonomy()
	termmeta := GetTermmeta()
	attachmentsPosts := SetMapAttachmentsBySizes()
	relationships := GetTermRelationships()
	posts := GetPosts()
	fabrics := GetFabrics(*posts)
	postMeta := GetPostmeta()


	//xType := fmt.Sprintf("%T", terms)
	//fmt.Println(xType) // "[]int"
	start_c := time.Now()
	SkCategories(*terms, *term_taxonomy, *termmeta, *attachmentsPosts, *relationships, *posts)
	duration_c := time.Since(start_c)
	fmt.Println("main SkCategories: ", duration_c)

	start_f := time.Now()
	SkFabrics(*fabrics, *postMeta, *terms, *termmeta, *attachmentsPosts)
	duration_f := time.Since(start_f)
	fmt.Println("main SkFabrics: ", duration_f)

	duration := time.Since(start)
	fmt.Println("main: ", duration)
}
