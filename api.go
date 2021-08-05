package main

import (
	"encoding/json"
	"fmt"
	"github.com/comhttp/json/jdb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oknors/okno/app/mod"
	"github.com/oknors/okno/pkg/utl"

	"net/http"
	"sort"
	"strconv"
	"time"
)

// Post contains articles and pages used by the CMS
type Post struct {
	ID      string `schema:"id,required"`
	Title   string
	Slug    string
	Date    time.Time
	Excerpt string
	Active  bool
	Order   int
}
type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Less(i, j int) bool { return p[i].Order < p[j].Order }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Handler() http.Handler {
	r := mux.NewRouter()

	r.StrictSlash(true)

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/{site}/{col}/{id}", viewPost).Methods("GET")
	r.HandleFunc("/{site}/{col}", viewAllPosts).Methods("GET")
	r.HandleFunc("/{site}/{col}/{per}/{page}/{truncate}", viewPosts).Methods("GET")
	r.Headers("Access-Control-Allow-Origin", "*")
	return handlers.CORS()(handlers.CompressHandler(InterceptHandler(r, DefaultErrorHandler)))
}

// HomeHandler handles a request for (?)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The Truth Is Out There..."))
}

func viewPosts(w http.ResponseWriter, r *http.Request) {
	per, _ := strconv.Atoi(mux.Vars(r)["per"])
	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	trn, _ := strconv.Atoi(mux.Vars(r)["truncate"])
	site := mux.Vars(r)["site"]
	col := mux.Vars(r)["col"]

	posts := Posts{}

	postsRaw, err := jdb.JDB.ReadAll("sites/" + site + "/jdb/" + col)
	utl.ErrorLog(err)
	fmt.Println("AAAAA:", postsRaw)
	for _, postInterface := range postsRaw {
		var rawPost mod.Post
		err := json.Unmarshal([]byte(postInterface), &rawPost)
		utl.ErrorLog(err)
		if !rawPost.IsDraft {
			p := Post{
				ID:      rawPost.ID,
				Title:   rawPost.Title,
				Slug:    rawPost.Slug,
				Date:    rawPost.CreatedAt,
				Order:   rawPost.Order,
				Excerpt: utl.Truncate(rawPost.Content, trn),
			}
			posts = append(posts, p)
		}
	}
	sort.Sort(posts)

	pn := len(posts)
	lb := map[string]interface{}{
		"currentPage": page,
		"pageCount":   int(pn) / per,
		"posts":       posts,
		"postsNumber": pn,
	}
	out, err := json.Marshal(lb)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}
	w.Write([]byte(out))
}

func viewAllPosts(w http.ResponseWriter, r *http.Request) {
	site := mux.Vars(r)["site"]
	col := mux.Vars(r)["col"]

	posts := mod.Posts{}

	postsRaw, err := jdb.JDB.ReadAll("sites/" + site + "/jdb/" + col)
	utl.ErrorLog(err)
	for _, postInterface := range postsRaw {
		var rawPost mod.Post
		err := json.Unmarshal([]byte(postInterface), &rawPost)
		utl.ErrorLog(err)
		if !rawPost.IsDraft {
			posts = append(posts, rawPost)
		}
	}
	sort.Sort(posts)

	out, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}
	w.Write([]byte(out))
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	site := mux.Vars(r)["site"]
	col := mux.Vars(r)["col"]
	id := mux.Vars(r)["id"]
	post := mod.Post{}
	err := jdb.JDB.Read("sites/"+site+"/jdb/"+col, id, &post)
	out, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}
	w.Write([]byte(out))
}
