package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/gocolly/colly"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 찾은 링크 하나하나
type (
	Item struct {
		Count     int    `json:"count"`
		StoryURL  string `json:"storyurl"`
		Source    string `json:"source"`
		comments  string
		CrawledAt time.Time
		Comments  string `json:"comments"`
		Title     string `json:"title"`
	}
)

// 어떤 포탈? 어떤 키워드?
const (
	max_item_count = 30
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: search.naver.com
		colly.AllowedDomains("search.naver.com"),
		// Parallelism
		colly.Async(true),
	)

	Items := []Item{}
	count := 0

	currentpagecount := 0
	// <dt><a class="sh_blog_title _sp_each_url _sp_each_title" href="link">
	c.OnHTML("dt a.sh_blog_title", func(e *colly.HTMLElement) {
		if count == max_item_count {
			return
		}
		if currentpagecount == 0 {
			currentpagecount = 1
		}

		item := Item{}
		//temp.comments = e.Attr("a", "class")
		item.StoryURL = e.Attr("href")
		if item.StoryURL == "" {
			return
		}
		item.Source = "https://search.naver.com/"
		item.Title = e.Text
		item.CrawledAt = time.Now()
		item.Count = count + 1

		Items = append(Items, item)
		count = item.Count
	})

	maxpagecount := 0
	c.OnHTML(".paging", func(e *colly.HTMLElement) {
		if maxpagecount != 0 {
			return
		}
		if pagecount := len(e.ChildAttrs("a", "href")); pagecount != 0 {
			maxpagecount = max_item_count / pagecount
		}
	})

	c.OnHTML(".next", func(e *colly.HTMLElement) {
		if maxpagecount == 0 {
			return
		}

		if currentpagecount == maxpagecount {
			return
		}

		currentpagecount = currentpagecount + 1
		c.Visit("https://search.naver.com/search.naver" + e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set max Parallelism and introduce a Random Delay
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/test", func(context echo.Context) error {
		// 재검색을 대비한 주요 객체 초기화
		defer func() {
			Items = nil
			count = 0
			maxpagecount = 0
			currentpagecount = 0
			c.Init()
		}()

		m := echo.Map{}
		if err := context.Bind(&m); err != nil {
			return err
		}

		// 일단은 네이버 블로그만
		// 차후 검색 가능 portal list 별도로 분리
		if m["portal"] != "naverblog" {
			return context.JSON(http.StatusBadRequest, m)
		}

		// 검색할 키워드 체크
		if m["keyword"] == "" {
			return context.JSON(http.StatusBadRequest, m)
		}

		// 아래의 Sprintf 사용하는 방식 대신 m["keyword"].(string) 이렇게 type assertion 사용하는 방식으로 사용함
		//Keyword := fmt.Sprintf("%v", m["keyword"])

		// 복수 키워드 가능 구현할 때 주석 해제
		//for _, keyword := range keywords {
		addressNaverBlog := "https://search.naver.com/search.naver?nso=so:dd,a:t,p:1w&srchby=title&st=date&where=post&query=" + m["keyword"].(string)
		c.Visit(addressNaverBlog)
		//}

		c.Wait()

		return context.JSON(http.StatusOK, Items)
	})

	e.Logger.Fatal(e.Start(":3000"))
}

// // 붕붕카, 최신순, 1주일, 제목만 체크
// https://search.naver.com/search.naver?date_from=20030520&date_option=3&date_to=20200630&dup_remove=1&nso=so%3Add%2Ca%3At%2Cp%3A1w&post_blogurl=&post_blogurl_without=&query=%EB%B6%95%EB%B6%95%EC%B9%B4&sm=tab_pge&srchby=title&st=date&where=post&start=1
// https://search.naver.com/search.naver?date_from=20030520&date_option=3&date_to=20200630&dup_remove=1&nso=so%3Add%2Ca%3At%2Cp%3A1w&post_blogurl=&post_blogurl_without=&query=%EB%B6%95%EB%B6%95%EC%B9%B4&sm=tab_pge&srchby=title&st=date&where=post&start=11

// ?date_from=20030520&date_option=3&date_to=20200630&dup_remove=1&nso=so%3Add%2Ca%3At%2Cp%3A1w&post_blogurl=&post_blogurl_without=&query=%EB%B6%95%EB%B6%95%EC%B9%B4&sm=tab_pge&srchby=title&st=date&where=post&start=11
// https://search.naver.com/search.naver?nso=so:dd,a:t,p:1w&post_blogurl=&post_blogurl_without=&query=붕붕카&sm=tab_pge&srchby=title&st=date&where=post&start=1

// // 아래가 최소 검색 페이지 까지 생략
// https://search.naver.com/search.naver?nso=so:dd,a:t,p:1w&srchby=title&st=date&where=post&query=붕붕카
