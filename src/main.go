package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	TGSTAT_IDRK  = ""
	TGSTAT_SIRK  = ""
	CF_CLEARANCE = ""
)

func extractTagFromURL(url string) string {
	parts := strings.Split(url, "/@")
	if len(parts) < 2 {
		return ""
	}
	tagWithPath := strings.Split(parts[1], "/")
	return "t.me/" + strings.ToLower(tagWithPath[0]) + "\n"
}

func processCategories(n []string) {
	file, _ := os.OpenFile("chatlinks.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
	})

	cookies := []*http.Cookie{
		{
			Name:   "tgstat_idrk",
			Value:  "913a83c1e07790d96f448f6ab67c09abdd55d1bb41d4874155ffd29c28b74adea%3A2%3A%7Bi%3A0%3Bs%3A11%3A%22tgstat_idrk%22%3Bi%3A1%3Bs%3A53%3A%22%5B11218702%2C%22b-6SuVjV_-zNV0VwDiG9KcKFKFGYnwgJ%22%2C2592000%5D%22%3B%7D", // замени на реальное значение
			Domain: "tgstat.ru",                                                                                                                                                                                                      // предположительный домен
			Path:   "/",
		},
		{
			Name:   "tgstat_sirk",
			Value:  "gf0efq3q6fgj3ih2628kv236dm",
			Domain: "tgstat.ru",
			Path:   "/",
		},
		{
			Name:   "cf_clearance",
			Value:  "5eRlLOPaQ_nkYFNNiixHfoXTnrar2pvMzbs_C4qlys8-1747746385-1.2.1.1-tFwcXMR1JspRhmpNuMKYMtNV6CAK7EDoxdgrK_4Pjiz.dwI4E4BYP5WzcN8j.baUkpv_TLrQ4q.aetjkZ9dNCLjX8wZ6TiwInc81lymDlXXJISCOlegmQw0qz62Wp_NhGHTjd4gQzDVCVobLT.2hMEIyjcQfRG6XBR3MTS1w6LsxT67io.rbVFb34bLWOsoJKAjHAaQ8yLky38SfCjn9_uE1rhStUuoVnOyBTdohXs7A1OKeHvBRavghrpD4tggMlZqc9zVS0cXaXAHq6waIG7oum8QKC6ba1bIN_Nk8SRV5_9PoBL.nnU2FEbynonYOBzkfwiqL0Nw9rcZnP8tGR7wFQgOVHNPDqksUVZTgMSg",
			Domain: "tgstat.ru",
			Path:   "/",
		},
	}

	c.SetCookies("https://tgstat.ru", cookies)

	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "stat") && strings.Contains(href, "chat/") {
			file.Write([]byte(extractTagFromURL(href)))
		}
	})

	for _, item := range n {
		url := "https://tgstat.ru/ratings/chats" + item + "/public?sort=msgs"
		if err := c.Visit(url); err != nil {
			log.Println("Error:", err)
		}
	}
}

func processChannels() {
	file, _ := os.OpenFile("chatlinks.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	scanner := bufio.NewScanner(file)
	defer file.Close()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	cookies := []*http.Cookie{
		{
			Name:   "tgstat_idrk",
			Value:  TGSTAT_IDRK,
			Domain: "tgstat.ru",
			Path:   "/",
		},
		{
			Name:   "tgstat_sirk",
			Value:  TGSTAT_SIRK,
			Domain: "tgstat.ru",
			Path:   "/",
		},
		{
			Name:   "cf_clearance",
			Value:  CF_CLEARANCE,
			Domain: "tgstat.ru",
			Path:   "/",
		},
	}

	c.SetCookies("https://tgstat.ru", cookies)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
	})

	c.OnHTML("a.btn", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "https://t.me") {
			file.Write([]byte((href + "\n")))
		}
	})

	for scanner.Scan() {
		line := scanner.Text()
		if err := c.Visit(line); err != nil {
			log.Println("Error:", err)
		}
	}

}

func main() {
	list := []string{}
	os.Remove("chatlinks.txt")

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
	})

	c.OnHTML("div.card.border.m-0 a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if !strings.Contains(href, "https://") && !strings.Contains(href, "#") {
			list = append(list, href)
		}
	})

	err := c.Visit("https://tgstat.ru/")
	if err != nil {
		log.Fatal(err)
	}

	processCategories(list)
	processChannels()

}
