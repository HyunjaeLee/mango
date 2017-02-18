package parse

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"regexp"
	"net/url"
	"fmt"
	"sync"
)

type Series struct {
	Id int
	Title string
	ThumbnailUrl string
	Timestamp int
	Day string
}

type Episode struct {
	Id int
	Title string
	ThumbnailUrl string
	Timestamp int
	SeriesId int
}

func parseId(url string) int {
	index := strings.LastIndex(url, "-") + 1
	id, _ := strconv.Atoi(url[index:])
	return id
}

func ParseIndex() []Series {
	days := []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
	seriez := []Series{}
	doc, _ := goquery.NewDocument("https://anigod.com/")
	doc.Find(".index-table-container").Each(func(i int, s *goquery.Selection) {
		if i < 7 {
			s.Find(".index-image-container.badge").Each(func(_ int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				id := parseId(href)
				title, _ := s.Attr("title")
				thumbnailUrl, _ := s.Find(".index-image").First().Attr("src")
				timestampString, _ := s.Attr("timestamp")
				timestamp, _ := strconv.Atoi(timestampString)

				series := Series{
					Id: id,
					Title: title,
					ThumbnailUrl: thumbnailUrl,
					Timestamp: timestamp,
					Day: days[i],
				}

				seriez = append(seriez, series)
			})
		}
	})
	return seriez
}

func ParseFinale() {
	c := make(chan Series)
	wg := new(sync.WaitGroup)
	wg.Add(9)
	for i := 1; i < 10;  i++ {
		go func() {
			ParseFinaleByPage(i, c)
			wg.Done()
		}()
	}
	go func() {
		for series := range c {
			fmt.Println(series)
		}
	}()

	wg.Wait()
}

func ParseFinaleByPage(page int, c chan<- Series) {
	doc, _ := goquery.NewDocument("https://anigod.com/animations/finale/title/asc/" + strconv.Itoa(page))
	doc.Find(".table-image-container").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id := parseId(href)
		title, _ := s.Attr("title")
		thumbnailUrl, _ := s.Find(".lazy").First().Attr("data-original")
		timestampString, _ := s.Attr("timestamp")
		timestamp, _ := strconv.Atoi(timestampString)

		series := Series{
			Id: id,
			Title: title,
			ThumbnailUrl: thumbnailUrl,
			Timestamp: timestamp,
		}

		c <- series
	})
}

func ParseEpisode(seriesId int, page int) []Episode {
	episodes := []Episode{}
	doc, _ := goquery.NewDocument("https://anigod.com/animation/" + strconv.Itoa(seriesId) + "/" + strconv.Itoa(page))
	doc.Find(".table-image-container").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id := parseId(href)
		title, _ := s.Attr("title")
		thumbnailUrl, _ := s.Find(".lazy").First().Attr("data-original")
		timestampString, _ := s.Attr("timestamp")
		timestamp, _ := strconv.Atoi(timestampString)

		episode := Episode{
			Id: id,
			Title: title,
			ThumbnailUrl: thumbnailUrl,
			Timestamp: timestamp,
			SeriesId: seriesId,
		}

		episodes = append(episodes, episode)
	})
	return episodes
}

func ParseVideoUrl(episodeId int) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://anigod.com/episode/" + strconv.Itoa(episodeId), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "http;//sh.st/")
	res, _ := client.Do(req)
	html, _ := ioutil.ReadAll(res.Body)
	pattern, _ := regexp.Compile("var videoID = '(.*?)'")
	s := pattern.FindStringSubmatch(string(html))
	videoUrl := s[1]
	videoUrl = strings.Replace(videoUrl, "\\x2b", "+", -1)
	videoUrl = strings.Replace(videoUrl, "\\", "", -1)
	escaped := url.QueryEscape(videoUrl)
	return "https://anigod.com/video?id=" + escaped
}