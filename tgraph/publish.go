package tgraph

import (
	"fmt"
	"html"
	"github.com/indes/flowerss-bot/log"
	"math/rand"
	"time"
)

func PublishHtml(sourceTitle string, title string, rawLink string, htmlContent string) (string, error) {
	//html = fmt.Sprintf(
	//	"<p>本文章由 干物妹！小埋 抓取自RSS，版權歸<a href=\"\">源站點</a>所有。</p><hr>",
	//) + html + fmt.Sprintf(
	//	"<hr><p>本文章由 干物妹！小埋 抓取自RSS，版權歸<a href=\"\">源站點</a>所有。</p><p>查看原文：<a href=\"%s\">%s - %s</p>",
	//	rawLink,
	//	title,
	//	sourceTitle,
	//)

	htmlContent = html.UnescapeString(htmlContent) + fmt.Sprintf(
		"<hr><p>本文章由 干物妹！小埋 抓取自RSS，版權歸<a href=\"\">源站點</a>所有。</p><p>查看原文：<a href=\"%s\">%s - %s</p>",
		rawLink,
		title,
		sourceTitle,
	)
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	client := clientPool[rand.Intn(len(clientPool))]

	if page, err := client.CreatePageWithHTML(title+" - "+sourceTitle, sourceTitle, rawLink, htmlContent, true); err == nil {
		log.Printf("Created telegraph page url: %s", page.URL)
		return page.URL, err
	} else {
		log.Printf("Create telegraph page error: %s", err)
		return "", nil
	}
}
