package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var dbs = DbConnect()

func strToInt(id string) int64 {
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		panic(err)
	}

	return idInt
}

type Terms struct {
	Term_id string `json:"term_id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Term_group string `json:"term_group"`
}

type WpTableTerms map[int64]Terms

func GetTerms() []WpTableTerms {

	var wp_terms []WpTableTerms

	rows, err := dbs.Query("SELECT * FROM wp_terms")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	terms := []Terms{}

	for rows.Next() {
		t := Terms{}
		err := rows.Scan(&t.Term_id, &t.Name, &t.Slug, &t.Term_group)

		if err != nil {
			fmt.Println(err)
			continue
		}

		terms = append(terms, t)
	}

	for _, t := range terms{
		term_id := strToInt(t.Term_id)
		terms_map := WpTableTerms{term_id: {t.Term_id, t.Name, t.Slug, t.Term_group}}
		wp_terms = append(wp_terms, terms_map)
	}

	return wp_terms
}

type TermTaxonomy struct {
	Term_taxonomy_id string `json:"term_taxonomy_id"`
	Term_id string `json:"term_id"`
	Taxonomy string `json:"taxonomy"`
	Description string `json:"description"`
	Parent string `json:"parent"`
	Count string `json:"count"`
}

type WpTableTermTaxonomy map[int64]TermTaxonomy

func GetTermTaxonomy() []WpTableTermTaxonomy {

	var wp_term_taxonomy []WpTableTermTaxonomy

	rows, err := dbs.Query("SELECT * FROM wp_term_taxonomy")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	term_taxonomy := []TermTaxonomy{}

	for rows.Next() {
		t := TermTaxonomy{}
		err := rows.Scan(&t.Term_taxonomy_id, &t.Term_id, &t.Taxonomy, &t.Description, &t.Parent, &t.Count)

		if err != nil {
			fmt.Println(err)
			continue
		}

		term_taxonomy = append(term_taxonomy, t)
	}

	for _, t := range term_taxonomy{
		term_taxonomy_id := strToInt(t.Term_taxonomy_id)
		taxonomy := WpTableTermTaxonomy{term_taxonomy_id: {t.Term_taxonomy_id, t.Term_id, t.Taxonomy, t.Description, t.Parent, t.Count}}
		wp_term_taxonomy = append(wp_term_taxonomy, taxonomy)
	}

	return wp_term_taxonomy
}

type Termmeta struct {
	Meta_id string `json:"meta_id"`
	Term_id string `json:"term_id"`
	Meta_key string `json:"meta_key"`
	Meta_value string `json:"meta_value"`
}

type TermmetaItems map[string]string

func GetTermmeta() map[int64][]TermmetaItems {

	rows, err := dbs.Query("SELECT * FROM wp_termmeta")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	termmeta := []Termmeta{}

	for rows.Next() {
		t := Termmeta{}
		err := rows.Scan(&t.Meta_id, &t.Term_id, &t.Meta_key, &t.Meta_value)

		if err != nil {
			fmt.Println(err)
			continue
		}

		termmeta = append(termmeta, t)
	}
	wp_termmeta := make(map[int64][]TermmetaItems)
	for _, t := range termmeta{
		tm := strToInt(t.Term_id)
		tmeta := TermmetaItems{t.Meta_key: t.Meta_value}
		wp_termmeta[tm] = append(wp_termmeta[tm], tmeta)
	}

	return wp_termmeta
}

type AttachmentsPost struct {
	Id string `json:"id"`
	Post_title string `json:"post_title"`
	Post_parent string `json:"post_parent"`
	Guid string `json:"guid"`
}

type WpPostsAttachments map[int64]AttachmentsPost

func GetAttachments() []WpPostsAttachments {

	rows, err := dbs.Query("SELECT ID, post_title, post_parent, guid FROM wp_posts WHERE post_type = 'attachment';")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	attachmentRows := []AttachmentsPost{}

	for rows.Next() {
		a := AttachmentsPost{}
		err := rows.Scan(&a.Id, &a.Post_title, &a.Post_parent, &a.Guid)

		if err != nil {
			fmt.Println(err)
			continue
		}
		attachmentRows = append(attachmentRows, a)
	}

	var attachmentsPosts []WpPostsAttachments

	for _, a := range attachmentRows {
		attachmentId := strToInt(a.Id)
		attachmentsItem := WpPostsAttachments{attachmentId: {a.Id, a.Post_title, a.Post_parent, a.Guid}}
		attachmentsPosts = append(attachmentsPosts, attachmentsItem)
	}

	return attachmentsPosts
}

type Attachment struct {
	Original string `json:"original"`
	W100 string `json:"w100"`
	W150 string `json:"w150"`
	W300 string `json:"w300"`
	W400 string `json:"w400"`
	W500 string `json:"w500"`
}

type sizeAttachments map[string]string

func SetMapAttachmentsBySizes() map[int64][]sizeAttachments {
	suffixes := [6]string{"--w_100", "--w_150", "--w_300", "--w_400", "--w_500", ""}
	mimeTypes := []string{"png", "jpg", "jpeg"}
	attachmentsAll := GetAttachments()
	dataList := make(map[int64][]sizeAttachments)

	for _, attach := range attachmentsAll {
		for attachId, attachItem := range attach {
			url := strings.Split(attachItem.Guid, ".")
			count := len(url)
			typeImg := url[len(url)-1]

			if Contains(mimeTypes, typeImg) {
				for _, suffix := range suffixes {
					var sb strings.Builder
					var newSlice []string

					newSlice = make([]string, len(url))
					copy(newSlice, url)

					sb.WriteString(newSlice[count-2])
					sb.WriteString(suffix)
					newSlice[count-2] = sb.String()
					sb.Reset()

					w := getSuffix(suffix)

					urlTmp := sizeAttachments{w: strings.Join(newSlice, ".")}
					dataList[attachId] = append(dataList[attachId], urlTmp)
					}
				} else {
					otherType := sizeAttachments{"url": attachItem.Guid}
					dataList[attachId] = append(dataList[attachId], otherType)
			}
			}
		}

	return dataList
}

func getSuffix(suffix string) string {
	switch suffix {
	case "--w_100":
		w := "w100"
		return w
	case "--w_150":
		w := "w150"
		return w
	case "--w_300":
		w := "w300"
		return w
	case "--w_400":
		w := "w400"
		return w
	case "--w_500":
		w := "w500"
		return w
	case "":
		w := "original"
		return w

	default:
		return ""
	}
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}