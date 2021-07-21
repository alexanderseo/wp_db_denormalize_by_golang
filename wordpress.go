package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

var dbs = DbConnect()

func strToInt(id string) int64 {
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		panic(err)
	}

	return idInt
}

type Term struct {
	Term_id string `json:"term_id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Term_group string `json:"term_group"`
}

type WpTerms map[int64]Term

func GetTerms() *[]WpTerms {
	start := time.Now()
	var wpTerms []WpTerms

	rows, err := dbs.Query("SELECT * FROM wp_terms")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var terms []Term
	for rows.Next() {
		t := Term{}
		err := rows.Scan(&t.Term_id, &t.Name, &t.Slug, &t.Term_group)

		if err != nil {
			fmt.Println(err)
			continue
		}

		terms = append(terms, t)
	}

	for _, t := range terms{
		termId := strToInt(t.Term_id)
		termsMap := WpTerms{termId: {t.Term_id, t.Name, t.Slug, t.Term_group}}
		wpTerms = append(wpTerms, termsMap)
	}
	duration := time.Since(start)
	fmt.Println("Terms: ", duration)

	return &wpTerms
}

type TermTaxonomy struct {
	Term_taxonomy_id string `json:"term_taxonomy_id"`
	Term_id string `json:"term_id"`
	Taxonomy string `json:"taxonomy"`
	Description string `json:"description"`
	Parent string `json:"parent"`
	Count string `json:"count"`
}

type WpTermTaxonomy map[int64]TermTaxonomy

func GetTermTaxonomy() *[]WpTermTaxonomy {
	start := time.Now()
	var wpTermTaxonomy []WpTermTaxonomy

	rows, err := dbs.Query("SELECT * FROM wp_term_taxonomy")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var termTaxonomy []TermTaxonomy

	for rows.Next() {
		t := TermTaxonomy{}
		err := rows.Scan(&t.Term_taxonomy_id, &t.Term_id, &t.Taxonomy, &t.Description, &t.Parent, &t.Count)

		if err != nil {
			fmt.Println(err)
			continue
		}

		termTaxonomy = append(termTaxonomy, t)
	}

	for _, t := range termTaxonomy {
		termTaxonomyId := strToInt(t.Term_taxonomy_id)
		taxonomy := WpTermTaxonomy {
			termTaxonomyId: {
				t.Term_taxonomy_id,
				t.Term_id,
				t.Taxonomy,
				t.Description,
				t.Parent,
				t.Count}}
		wpTermTaxonomy = append(wpTermTaxonomy, taxonomy)
	}
	duration := time.Since(start)
	fmt.Println("Taxonomy: ", duration)
	return &wpTermTaxonomy
}

type Termmeta struct {
	Meta_id string `json:"meta_id"`
	Term_id string `json:"term_id"`
	Meta_key string `json:"meta_key"`
	Meta_value string `json:"meta_value"`
}

type TermmetaItems map[string]string

func GetTermmeta() *map[int64][]TermmetaItems {
	start := time.Now()
	rows, err := dbs.Query("SELECT * FROM wp_termmeta")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var termmeta []Termmeta

	for rows.Next() {
		t := Termmeta{}
		err := rows.Scan(&t.Meta_id, &t.Term_id, &t.Meta_key, &t.Meta_value)

		if err != nil {
			fmt.Println(err)
			continue
		}

		termmeta = append(termmeta, t)
	}
	wpTermmeta := make(map[int64][]TermmetaItems)
	for _, t := range termmeta{
		tm := strToInt(t.Term_id)
		tmeta := TermmetaItems{t.Meta_key: t.Meta_value}
		wpTermmeta[tm] = append(wpTermmeta[tm], tmeta)
	}
	duration := time.Since(start)
	fmt.Println("Termmeta: ", duration)
	return &wpTermmeta
}

type AttachmentsPost struct {
	Id string `json:"id"`
	Post_title string `json:"post_title"`
	Post_parent string `json:"post_parent"`
	Guid string `json:"guid"`
}

type WpPostsAttachments map[int64]AttachmentsPost

func GetAttachments() *[]WpPostsAttachments {
	start := time.Now()
	rows, err := dbs.Query("SELECT ID, post_title, post_parent, guid FROM wp_posts WHERE post_type = 'attachment';")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var attachmentRows []AttachmentsPost

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
		attachmentsItem := WpPostsAttachments{
			attachmentId: {
				Id: a.Id,
				Post_title: a.Post_title,
				Post_parent: a.Post_parent,
				Guid: a.Guid}}
		attachmentsPosts = append(attachmentsPosts, attachmentsItem)
	}
	duration := time.Since(start)
	fmt.Println("Attachments: ", duration)
	return &attachmentsPosts
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

func SetMapAttachmentsBySizes() *map[int64][]sizeAttachments {
	start := time.Now()
	suffixes := [6]string{"--w_100", "--w_150", "--w_300", "--w_400", "--w_500", ""}
	mimeTypes := []string{"png", "jpg", "jpeg"}
	attachmentsAll := GetAttachments()
	dataList := make(map[int64][]sizeAttachments)

	for _, attach := range *attachmentsAll {
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
	duration := time.Since(start)
	fmt.Println("AttachmentsBySizes: ", duration)
	return &dataList
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

type TermRelationships struct {
	Object_id string `json:"object_id"`
	Term_taxonomy_id string `json:"term_taxonomy_id"`
}

type WpTableTermRelationships map[int64][]TermRelationships

func GetTermRelationships() *WpTableTermRelationships {
	start := time.Now()
	rows, err := dbs.Query("SELECT object_id, term_taxonomy_id FROM wp_term_relationships")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var relationships []TermRelationships

	for rows.Next() {
		r := TermRelationships{}
		err := rows.Scan(&r.Object_id, &r.Term_taxonomy_id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		relationships = append(relationships, r)
	}

	wp_term_relationships := make(WpTableTermRelationships)

	for _, rl := range relationships {
		ri := strToInt(rl.Object_id)
		tmp_r := TermRelationships{rl.Object_id, rl.Term_taxonomy_id}
		wp_term_relationships[ri] = append(wp_term_relationships[ri], tmp_r)
	}
	duration := time.Since(start)
	fmt.Println("TermRelationships: ", duration)
	fmt.Println("LEN: ", len(wp_term_relationships))
	return &wp_term_relationships
}

type Post struct {
	Id string `json:"id"`
	Post_content string `json:"post_content"`
	Post_title string `json:"post_title"`
	Post_status string `json:"post_status"`
	Post_parent string `json:"post_parent"`
	Guid string `json:"guid"`
	Post_name string `json:"post_name"`
	Post_date string `json:"post_date"`
	Post_modified string `json:"post_modified"`
	Post_type string `json:"post_type"`
	Menu_order string `json:"menu_order"`
}

type WpPosts map[int64]Post

func GetPosts() *[]WpPosts {
	start := time.Now()
	rows, err := dbs.Query("SELECT ID, post_content, post_title, post_status, post_parent, guid, post_name, post_date, post_modified, post_type, menu_order FROM wp_posts WHERE post_status = 'publish' OR post_status = 'inherit' OR post_status = 'future'")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		p := Post{}
		err := rows.Scan(
			&p.Id,
			&p.Post_content,
			&p.Post_title,
			&p.Post_status,
			&p.Post_parent,
			&p.Guid,
			&p.Post_name,
			&p.Post_date,
			&p.Post_modified,
			&p.Post_type,
			&p.Menu_order)

		if err != nil {
			fmt.Println(err)
			continue
		}
		posts = append(posts, p)
	}

	var wpPostList []WpPosts
	for _, post := range posts {
		pi := strToInt(post.Id)
		wpPosts := WpPosts{pi: {
			Id:post.Id,
			Post_content: post.Post_content,
			Post_title: post.Post_title,
			Post_status: post.Post_status,
			Post_parent: post.Post_parent,
			Guid: post.Guid,
			Post_name: post.Post_name,
			Post_date: post.Post_date,
			Post_modified: post.Post_modified,
			Post_type: post.Post_type,
			Menu_order: post.Menu_order}}
		wpPostList = append(wpPostList, wpPosts)
	}
	duration := time.Since(start)
	fmt.Println("GetPosts: ", duration)
	fmt.Println("LEN: ", len(wpPostList))
	return &wpPostList
}

type WpFabrics map[int64]Post

func GetFabrics(fabrics []WpPosts) *[]WpFabrics {
	start := time.Now()
	fb := fabrics

	var wp_fabrics []WpFabrics
	for _, f := range fb {
		for _, fabric := range f {
			if fabric.Post_type == "fabric" {
				fi := strToInt(fabric.Id)
				wp_fabric := WpFabrics{fi: {
					Id: fabric.Id,
					Post_content: fabric.Post_content,
					Post_parent: fabric.Post_parent,
					Post_name: fabric.Post_name,
					Post_title: fabric.Post_title,
					Post_status: fabric.Post_status,
					Post_date: fabric.Post_date,
					Post_modified: fabric.Post_modified,
					Post_type: fabric.Post_type,
					Guid: fabric.Guid,
					Menu_order: fabric.Menu_order}}
				wp_fabrics = append(wp_fabrics, wp_fabric)
			}
		}
	}
	duration := time.Since(start)
	fmt.Println("GetFabrics: ", duration)
	fmt.Println("LEN: ", len(wp_fabrics))
	return &wp_fabrics
}

type Postmeta struct {
	Post_id string `json:"post_id"`
	Meta_key string `json:"meta_key"`
	Meta_value sql.NullString `json:"meta_value"`
}

func GetPostmeta() *map[int64][]Postmeta {
	start := time.Now()
	rows, err := dbs.Query("SELECT post_id, meta_key, meta_value FROM wp_postmeta")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var postMeta []Postmeta

	for rows.Next() {
		pm := Postmeta{}
		err := rows.Scan(&pm.Post_id, &pm.Meta_key, &pm.Meta_value)
		if err != nil {
			fmt.Println(err)
			continue
		}

		postMeta = append(postMeta, pm)
	}

	wpPostMeta := make(map[int64][]Postmeta)
	for _, m := range postMeta {
		mi := strToInt(m.Post_id)
		wpm := Postmeta{
			Post_id: m.Post_id,
			Meta_key: m.Meta_key,
			Meta_value: m.Meta_value,
		}
		wpPostMeta[mi] = append(wpPostMeta[mi], wpm)
	}
	duration := time.Since(start)
	fmt.Println("PostMeta: ", duration)
	fmt.Println("LEN: ", len(wpPostMeta))
	return &wpPostMeta
}