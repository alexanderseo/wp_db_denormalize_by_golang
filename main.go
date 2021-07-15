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
	//xType := fmt.Sprintf("%T", terms)
	//fmt.Println(xType) // "[]int"

	SkCategories(terms, term_taxonomy, termmeta, attachmentsPosts, relationships, posts)

	duration := time.Since(start)
	fmt.Println(duration)
}
