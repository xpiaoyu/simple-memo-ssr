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
	RouteArticleList   = "/list"
	RouteGetArticle    = "/get/*id"
	RouteGetArticleOld = "/get"
	RoutePostArticle   = "/post"
	RouteCreateArticle = "/create"
	RouteAssets        = "/assets/:p"
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
	log.Println("Root path:", rootPath)
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
	// Show article directory
	router.GET(RouteIndex, getDir)
	//
	router.GET(RouteDir, getDir)

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

func createArticle(c *fasthttp.RequestCtx) {
	if string(c.Method()) == "OPTIONS" {
		c.SetStatusCode(204)
		c.Response.Header.Set("access-control-allow-headers", "content-type")
		return
	}
	t := new(struct {
		Id string `json:"id"`
	})
	if err := json.Unmarshal(c.PostBody(), t); err != nil {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println(err)
	}
	articleId := t.Id
	if len(articleId) < 1 {
		c.SetStatusCode(400)
		if _, err := c.WriteString("article id invalid"); err != nil {
			log.Println("[ERROR]", err)
		}
		return
	}
	filename := "article/" + articleId + ".md"
	if canCreateFile(filename) {
		a := new(Article)
		a.Id = articleId
		a.Markdown = "# " + articleId + "\n"
		err := ioutil.WriteFile(filename, []byte(a.Markdown), os.ModePerm)
		if err != nil {
			log.Println("[error] can't write file err msg:", err)
			c.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		fi, err := os.Stat(filename)
		if err != nil {
			c.SetStatusCode(fasthttp.StatusInternalServerError)
			log.Println(err)
			return
		}
		a.Timestamp = fi.ModTime().UnixNano() / 1e6
		ArticleMap[articleId] = a
		ArticleList = append(ArticleList, a)
		sort.Sort(ArticleList)
		c.SetStatusCode(fasthttp.StatusOK)
		if _, err := c.WriteString("success"); err != nil {
			log.Println("[ERROR]", err)
		}
	} else {
		c.SetStatusCode(fasthttp.StatusOK)
		if _, err := c.WriteString("existed"); err != nil {
			log.Println("[ERROR]", err)
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
	log.Println("Path:", p)
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
	log.Println("article id:", articleId, "key:", b2s(key))
	//article, ok := ArticleMap[articleId]
	//if !ok {
	//	c.SetStatusCode(fasthttp.StatusNotFound)
	//	if _, err := c.WriteString("Article Not Found"); err != nil {
	//		log.Println("[ERROR]", err)
	//	}
	//	return
	//}
	_, html, err := getMarkdownAndHtml("article" + articleId)
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
			log.Println("Article name:", v.Name())
			t := new(Article)
			t.Id = strings.Replace(v.Name(), ".md", "", -1)
			if len(t.Id) < 1 {
				log.Println("[error] article id length invalid")
				os.Exit(1)
			}
			t.Timestamp = v.ModTime().UnixNano() / 1e6
			md, html, err := getMarkdownAndHtml("article/" + v.Name())
			if err != nil {
				log.Println("[error] getSummaryAndMarkdown err:", err)
				os.Exit(1)
			}
			t.Markdown = b2s(md)
			t.Html = html
			ArticleList = append(ArticleList, t)
			ArticleMap[t.Id] = t
		}
	}
	sort.Sort(ArticleList)
	log.Println("Scan article directory successfully")
	return
}

func getMarkdownAndHtml(filename string) (markdown []byte, html []byte, err error) {
	b := cacheGet(filename)
	if b != nil {
		// Hit cache
		markdown = b
		html = b
		return
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("[ERROR]", "读取", filename, "失败")
		return
	}
	log.Println("读取", filename, "成功")
	//markdown = string(bytes)
	markdown = bytes
	html = MarkdownToHtml(bytes)
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
