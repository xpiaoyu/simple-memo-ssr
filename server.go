package main

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	FasthttpAddr       = ":8083"
	RouteIndex         = "/"
	RouteDir           = "/dir/*path"
	RouteGetArticle    = "/get/*id"
	RouteGetArticleOld = "/get"
	RoutePostArticle   = "/post"
	RouteCreateArticle = "/create"
	RouteAssets        = "/assets/*p"
	RouteEdit          = "/edit/*p"
	ContentTypeJson    = "application/json"
	ContentTypeHtml    = "text/html"
)

type Article struct {
	Id        string `json:"id"`
	Markdown  string `json:"markdown"`
	Html      []byte `json:"html"`
	Timestamp int64  `json:"timestamp"`
}

type UploadPost struct {
	Md  string `json:"md"`
	Sum string `json:"sum"`
	Id  string `json:"id"`
}

type TplArticle struct {
	Title string
	Html  template.HTML
}

type ArticlePointArray []*Article

type FileAndDir struct {
	Name    string
	IsDir   bool
	ModTime string
	Size    string
}

type TplIndex struct {
	FileList []FileAndDir
	Prefix   string
}

type TplEdit struct {
	Markdown template.HTML
	Path     string
}

func (c ArticlePointArray) Len() int {
	return len(c)
}
func (c ArticlePointArray) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ArticlePointArray) Less(i, j int) bool {
	return c[i].Timestamp > c[j].Timestamp
}

var ArticleList ArticlePointArray
var ArticleMap map[string]*Article
var tpl *template.Template
var rootPath string
var fileSizeLevel = []string{"", "K", "M", "G", "T"}

func init() {
	rootPath = filepath.Dir(os.Args[0])
	tpl = template.Must(template.ParseGlob("template/*.html"))
}

func main() {
	log.Println("[I] Root path:", rootPath)
	ArticleMap = make(map[string]*Article)
	scanArticleDir()
	router := fasthttprouter.New()
	// Static files
	router.GET(RouteAssets, func(c *fasthttp.RequestCtx) {
		filePath := c.UserValue("p").(string)
		filePath = strings.Replace(filePath, "..", "", -1)
		c.SendFile("assets/" + filePath)
	})
	// Show article with keyword highlight
	router.GET(RouteGetArticle, getArticle)
	// Compatible with old router
	router.GET(RouteGetArticleOld, getArticle)
	// Show index article directory
	router.GET(RouteIndex, getDir)
	// Show article directory
	router.GET(RouteDir, getDir)
	// Edit markdown
	router.GET(RouteEdit, editHandler)
	// Post markdown
	router.POST(RoutePostArticle, postHandler)
	// Create new markdown
	router.POST(RouteCreateArticle, createHandler)

	/*firstHandler := func(c *fasthttp.RequestCtx) {
		c.Response.Header.Add("Access-Control-Allow-Origin", "*")
		switch string(c.Path()) {
		case RouteAssets:
			filePath := string(c.QueryArgs().Peek("p"))
			filePath = strings.Replace(filePath, "..", "", -1)
			c.SendFile("assets/" + filePath)
		case RouteArticleList:
			getArticleList(c)
		case RouteGetArticle:
			getArticle(c)
		case RoutePostArticle:
			postArticle(c)
		case RouteCreateArticle:
			createArticle(c)
		default:
			c.SetStatusCode(401)
			if _, err := c.WriteString("Unrecognized request."); err != nil {
				log.Println("[ERROR]", err)
			}
		}
	}*/
	log.Fatal(fasthttp.ListenAndServe(FasthttpAddr, router.Handler))
}

func postHandler(c *fasthttp.RequestCtx) {
	path := c.FormValue("path")
	md := c.FormValue("markdown")
	log.Println("[I] 修改文件path:", string(path))
	filename := "article" + b2s(path)
	if err := ioutil.WriteFile(filename, md, os.ModePerm); err != nil {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("[E]", err)
		return
	}
	if _, err := c.WriteString("ok"); err != nil {
		log.Println("[E]", err)
	}
}

func editHandler(c *fasthttp.RequestCtx) {
	p := c.UserValue("p")
	if p == nil {
		c.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	path := p.(string)
	path = strings.Replace(path, "..", "", -1)
	md, err := getMarkdown("article" + path)
	if err != nil {
		if _, err = c.WriteString("404 Not Found"); err != nil {
			log.Println("[E]")
		}
		c.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	c.SetContentType(ContentTypeHtml)
	if err := tpl.ExecuteTemplate(c, "edit.html", TplEdit{
		Markdown: template.HTML(md),
		Path:     path,
	}); err != nil {
		log.Println("[ERROR]", err)
	}
}

func createHandler(c *fasthttp.RequestCtx) {
	/*if string(c.Method()) == "OPTIONS" {
		c.SetStatusCode(204)
		c.Response.Header.Set("access-control-allow-headers", "content-type")
		return
	}*/
	t := c.FormValue("title")
	p := c.FormValue("path")
	articleId := string(t)
	articleId = strings.Replace(articleId, "..", "", -1)
	path := string(p)
	path = strings.Replace(path, "..", "", -1)
	if len(articleId) < 1 {
		c.SetStatusCode(400)
		if _, err := c.WriteString("Invalid title."); err != nil {
			log.Println("[E]", err)
		}
		return
	}
	filename := "article/" + path + "/" + articleId + ".md"
	if canCreateFile(filename) {
		err := ioutil.WriteFile(filename, []byte(articleId), os.ModePerm)
		if err != nil {
			log.Println("[E] can't write file err msg:", err)
			c.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		_, err = os.Stat(filename)
		if err != nil {
			c.SetStatusCode(fasthttp.StatusInternalServerError)
			log.Println(err)
			return
		}
		if _, err := c.WriteString("ok"); err != nil {
			log.Println("[E]", err)
		}
	} else {
		c.SetStatusCode(fasthttp.StatusOK)
		if _, err := c.WriteString("existed"); err != nil {
			log.Println("[E]", err)
		}
	}
}

func listDirectory(path string) (ret []FileAndDir, err error) {
	files, err := ioutil.ReadDir("article" + path)
	if err != nil {
		return nil, err
	}
	dir := make([]FileAndDir, 0, 16)
	ret = make([]FileAndDir, 0, 8)
	for _, v := range files {
		if v.IsDir() {
			dir = append(dir, FileAndDir{
				Name:    v.Name(),
				IsDir:   v.IsDir(),
				ModTime: v.ModTime().Format("2006-01-02 15:04"),
				Size:    "-",
			})
		} else {
			ret = append(ret, FileAndDir{
				Name:    v.Name(),
				IsDir:   v.IsDir(),
				ModTime: v.ModTime().Format("2006-01-02 15:04"),
				Size:    getFileSizeString(v.Size()),
			})
		}
	}
	ret = append(dir, ret...)
	return
}

func getFileSizeString(size int64) string {
	sizeF := float32(size)
	level := 0
	for {
		if sizeF >= 1024 && level < 4 {
			sizeF /= 1024
			level++
		} else {
			break
		}
	}
	if level == 0 {
		return fmt.Sprintf("%.0f", sizeF)
	} else {
		return fmt.Sprintf("%.1f%s", sizeF, fileSizeLevel[level])
	}
}

func getDir(c *fasthttp.RequestCtx) {
	path := c.UserValue("path")
	if path == nil {
		path = "/"
	}
	p := path.(string)
	if !strings.HasSuffix(p, "/") {
		c.Redirect(b2s(c.Path())+"/", 301)
		return
	}
	p = strings.Replace(p, "..", "", -1)
	log.Println("[I] Path:", p)
	c.SetContentType(ContentTypeHtml)
	fl, err := listDirectory(p)
	if err != nil {
		c.NotFound()
		return
	}
	if err := tpl.ExecuteTemplate(c, "index.html", TplIndex{
		FileList: fl,
		Prefix:   p,
	}); err != nil {
		log.Println("[ERROR]", err)
	}
}

func postArticle(c *fasthttp.RequestCtx) {
	if string(c.Method()) == "OPTIONS" {
		c.SetStatusCode(204)
		c.Response.Header.Set("access-control-allow-headers", "content-type")
		return
	}
	//markdown := string(c.PostArgs().Peek("md"))
	//summary := string(c.PostArgs().Peek("sum"))
	//body := string(c.PostBody())
	upload := new(UploadPost)
	err := json.Unmarshal(c.PostBody(), upload)
	if err != nil {
		c.SetStatusCode(400)
		return
	}
	t, ok := ArticleMap[upload.Id]
	if !ok {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("can't find article in map, id:", upload.Id)
		return
	}
	bytes := []byte(upload.Md)
	filename := "article/" + upload.Id + ".md"
	if err := ioutil.WriteFile(filename, bytes, os.ModePerm); err != nil {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("[ERROR]", err)
		return
	}
	t.Id = upload.Id
	t.Markdown = upload.Md
	//t.Html = MarkdownToHtml(upload.Md)
	fi, err := os.Stat(filename)
	if err != nil {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("[ERROR]", err)
		return
	}
	t.Timestamp = fi.ModTime().UnixNano() / 1e6
	sort.Sort(ArticleList)
	if _, err := c.WriteString("success"); err != nil {
		log.Println("[ERROR]", err)
	}
}

func getArticle(c *fasthttp.RequestCtx) {
	//articleId := b2s(c.QueryArgs().Peek("id"))
	var articleId string
	id := c.UserValue("id")
	if id != nil {
		articleId = id.(string)
	} else {
		articleId = b2s(c.QueryArgs().Peek("id"))
		c.Redirect(string(c.Path())+"/"+articleId, 301)
		return
	}
	articleId = strings.Replace(articleId, "..", "", -1)

	key := c.QueryArgs().Peek("k")
	log.Println("[I] article id:", articleId, "key:", b2s(key))
	//article, ok := ArticleMap[articleId]
	//if !ok {
	//	c.SetStatusCode(fasthttp.StatusNotFound)
	//	if _, err := c.WriteString("Article Not Found"); err != nil {
	//		log.Println("[ERROR]", err)
	//	}
	//	return
	//}
	html, err := getHtml("article" + articleId)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	c.SetContentType(ContentTypeHtml)
	output := html
	if utf8.RuneCount(key) >= 2 {
		output = HighlightKeywordBytes(output, key)
	}
	if err := tpl.ExecuteTemplate(c, "article.html",
		TplArticle{
			Title: articleId,
			Html:  template.HTML(output),
		}); err != nil {
		log.Println("[ERROR]", err)
	}
}

func scanArticleDir() {
	ArticleList = *new(ArticlePointArray)
	files, err := ioutil.ReadDir("article")
	if err != nil {
		log.Println("[error]", err)
		os.Exit(1)
	}
	for _, v := range files {
		if strings.HasSuffix(v.Name(), ".md") {
			log.Println("[I] Article name:", v.Name())
			t := new(Article)
			t.Id = strings.Replace(v.Name(), ".md", "", -1)
			if len(t.Id) < 1 {
				log.Println("[error] article id length invalid")
				os.Exit(1)
			}
			t.Timestamp = v.ModTime().UnixNano() / 1e6
			html, err := getHtml("article/" + v.Name())
			if err != nil {
				log.Println("[error] getSummaryAndMarkdown err:", err)
				os.Exit(1)
			}
			t.Markdown = b2s(html)
			t.Html = html
			ArticleList = append(ArticleList, t)
			ArticleMap[t.Id] = t
		}
	}
	sort.Sort(ArticleList)
	log.Println("[I] Scan article directory successfully")
	return
}

func getMarkdown(filename string) (markdown []byte, err error) {
	markdown, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Println("[E] 读取文件失败", err)
		return
	}
	log.Println("[I] 读取", filename, "成功")
	return
}

func getHtml(filename string) (html []byte, err error) {
	b := cacheGet(filename)
	if b != nil {
		// Hit cache
		html = b
		return
	}
	markdown, err := getMarkdown(filename)
	if err != nil {
		log.Printf("[E] 读取失败 name:%s error:%s", filename, err)
		return
	}
	html = MarkdownToHtml(markdown)
	cacheSet(filename, html, 5*1000)
	return
}

func getArticleList(c *fasthttp.RequestCtx) {
	c.SetContentType(ContentTypeHtml)
	/*c.SetContentType(ContentTypeJson)
	if err := json.NewEncoder(c).Encode(ArticleList); err != nil {
		log.Println("[ERROR]", err)
	}*/
}

func canCreateFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// path is not existed
			return true
		} else {
			// unknown error
			return false
		}
	}
	return false
}
