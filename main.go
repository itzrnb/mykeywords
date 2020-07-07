package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type item struct {
	StoryURL  string
	Source    string
	comments  string
	CrawledAt time.Time
	Comments  string
	Title     string
}

//var links = []string{}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, links)
}

func main() {
	stories := []item{}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: search.naver.com
		colly.AllowedDomains("search.naver.com"),
		// Parallelism
		colly.Async(true),
	)

	// On every a element which has .top-matter attribute call callback
	// This class is unique to the div that holds all information about a story

	// <div class="top-matter"><p class="title"><a class="title may-blank outbound" data-event-action="title" href=".....
	// On every a element which has .top-matter attribute call callback
	// This class is unique to the div that holds all information about a story
	// c.OnHTML(".top-matter", func(e *colly.HTMLElement) {
	// 	temp := item{}
	// 	temp.StoryURL = e.ChildAttr("a[data-event-action=title]", "href")
	// 	temp.Source = "https://old.reddit.com/r/programming/"
	// 	temp.Title = e.ChildText("a[data-event-action=title]")
	// 	temp.Comments = e.ChildAttr("a[data-event-action=comments]", "href")
	// 	temp.CrawledAt = time.Now()
	// 	stories = append(stories, temp)
	// })

	//<table border="0" class='comment-tree'>
	//		<tr class='athing comtr
	//c.OnHTML(".comment-tree tr.athing", func(e *colly.HTMLElement) {

	count := 0
	// <dt><a class="sh_blog_title _sp_each_url _sp_each_title" href="link">
	c.OnHTML("dt a.sh_blog_title", func(e *colly.HTMLElement) {
		temp := item{}
		//temp.comments = e.Attr("a", "class")
		temp.StoryURL = e.Attr("href")
		if temp.StoryURL == "" {
			return
		}
		// temp.Source = "https://search.naver.com/"
		temp.Title = e.Text
		// temp.Comments = e.ChildAttr("a[data-event-action=comments]", "href")
		// temp.CrawledAt = time.Now()

		//temp.StoryURL = e.Attr("href")
		// temp.Source = "https://search.naver.com/"
		//temp.Title = e.ChildText("a[data-event-action=title]")
		// temp.Comments = e.ChildAttr("a[data-event-action=comments]", "href")
		// temp.CrawledAt = time.Now()
		// matched, _ := regexp.MatchString(`https://blog.naver.com/.*`, e.Attr("href"))
		// if matched == true {
		// 	fmt.Println(e.Attr("href"))
		// 	links = append(links, e.Attr("href")+"\n")
		// }

		stories = append(stories, temp)
		count = count + 1
		fmt.Println("Link:", count, temp.StoryURL)
		fmt.Println("Title:", count, temp.Title)
	})

	// // On every span tag with the class next-button
	//<span class="next-button">
	// c.OnHTML("span.next-button", func(h *colly.HTMLElement) {
	// 	t := h.ChildAttr("a", "href")
	// 	c.Visit(t)
	// })

	// Set max Parallelism and introduce a Random Delay
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if len(os.Args) < 2 {
		log.Println("검색할 키워드를 입력해 주세요.")
		os.Exit(1)
	}

	// Crawl all reddits the user passes in
	//keywords := []string{"https://search.naver.com/search.naver?where=post&sm=tab_jum&query=붕붕카", "https://search.naver.com/search.naver?where=post&sm=tab_jum&query=밥상"}
	keywords := os.Args[1:]
	for _, keyword := range keywords {
		addressNaverBlog := "https://search.naver.com/search.naver?nso=so:dd,a:t,p:1w&srchby=title&st=date&where=post&query=" + keyword
		c.Visit(addressNaverBlog)
		fmt.Println(addressNaverBlog)
	}

	// 	fmt.Println(addressNaverBlog)

	c.Wait()
	//fmt.Println(stories)
	//fmt.Println(links)

	http.HandleFunc("/", handler)
	// 일단 아래 웹서버 비활성화
	//log.Fatal(http.ListenAndServe(":7777", nil))
}

// // 붕붕카, 최신순, 1주일, 제목만 체크
// https://search.naver.com/search.naver?date_from=20030520&date_option=3&date_to=20200630&dup_remove=1&nso=so%3Add%2Ca%3At%2Cp%3A1w&post_blogurl=&post_blogurl_without=&query=%EB%B6%95%EB%B6%95%EC%B9%B4&sm=tab_pge&srchby=title&st=date&where=post&start=1
// https://search.naver.com/search.naver?date_from=20030520&date_option=3&date_to=20200630&dup_remove=1&nso=so%3Add%2Ca%3At%2Cp%3A1w&post_blogurl=&post_blogurl_without=&query=%EB%B6%95%EB%B6%95%EC%B9%B4&sm=tab_pge&srchby=title&st=date&where=post&start=11

// ?date_from=20030520&date_option=3&date_to=20200630&dup_remove=1&nso=so%3Add%2Ca%3At%2Cp%3A1w&post_blogurl=&post_blogurl_without=&query=%EB%B6%95%EB%B6%95%EC%B9%B4&sm=tab_pge&srchby=title&st=date&where=post&start=11
// https://search.naver.com/search.naver?nso=so:dd,a:t,p:1w&post_blogurl=&post_blogurl_without=&query=붕붕카&sm=tab_pge&srchby=title&st=date&where=post&start=1

// // 아래가 최소 검색 페이지 까지 생략
// https://search.naver.com/search.naver?nso=so:dd,a:t,p:1w&srchby=title&st=date&where=post&query=붕붕카

// c.OnHTML(".top-matter", func(e *colly.HTMLElement) {
// 		temp := item{}
// 		temp.StoryURL = e.ChildAttr("a[data-event-action=title]", "href")
// 		temp.Source = "https://old.reddit.com/r/programming/"
// 		temp.Title = e.ChildText("a[data-event-action=title]")
// 		temp.Comments = e.ChildAttr("a[data-event-action=comments]", "href")
// 		temp.CrawledAt = time.Now()
// 		stories = append(stories, temp)
// 	})
