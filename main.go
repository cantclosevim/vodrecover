package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var vodDomains = []string{
	"https://vod-secure.twitch.tv",
	"https://vod-metro.twitch.tv",
	"https://vod-pop-secure.twitch.tv",
	"https://d2e2de1etea730.cloudfront.net",
	"https://dqrpb9wgowsf5.cloudfront.net",
	"https://ds0h3roq6wcgc.cloudfront.net",
	"https://d2nvs31859zcd8.cloudfront.net",
	"https://d2aba1wr3818hz.cloudfront.net",
	"https://d3c27h4odz752x.cloudfront.net",
	"https://dgeft87wbj63p.cloudfront.net",
	"https://d1m7jfoe9zdc1j.cloudfront.net",
	"https://d3vd9lfkzbru3h.cloudfront.net",
	"https://d2vjef5jvl6bfs.cloudfront.net",
	"https://d1ymi26ma8va5x.cloudfront.net",
	"https://d1mhjrowxxagfy.cloudfront.net",
	"https://ddacn6pr5v0tl.cloudfront.net",
	"https://d3aqoihi2n8ty8.cloudfront.net",
	"https://d1xhnb4ptk05mw.cloudfront.net",
	"https://d6tizftlrpuof.cloudfront.net",
	"https://d36nr0u3xmc4mm.cloudfront.net",
	"https://d1oca24q5dwo6d.cloudfront.net",
	"https://d2um2qdswy1tb0.cloudfront.net",
	"https://d1w2poirtb3as9.cloudfront.net",
	"https://d6d4ismr40iw.cloudfront.net",
	"https://d1g1f25tn8m2e6.cloudfront.net",
	"https://dykkng5hnh52u.cloudfront.net",
	"https://d2dylwb3shzel1.cloudfront.net",
	"https://d2xmjdvx03ij56.cloudfront.net",
}

func main() {
	var verbose bool
	var trackerUrl string

	flag.BoolVar(&verbose, "v", false, "Display verbose output")
	flag.BoolVar(&verbose, "verbose", false, "Display verbose output")
	flag.StringVar(&trackerUrl, "u", "", "TwitchTracker url to be used")
	flag.StringVar(&trackerUrl, "url", "", "TwitchTracker url to be used")
	flag.Parse()

	if trackerUrl == "" {
		fmt.Print("Enter a TwitchTracker VOD url: ")
		fmt.Scanln(&trackerUrl)
	}

	req, err := http.NewRequest("GET", trackerUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	startTimestamp := doc.Find("div.stream-timestamp-dt").First().Text()
	endTimestamp := doc.Find("div.text-right div.stream-timestamp-dt").Text()

	t, err := time.Parse("2006-01-02 15:04:05", startTimestamp)
	if err != nil {
		log.Fatal(err)
	}

	channel := regexp.MustCompile(`\/([a-zA-Z0-9_]+)\/streams`).FindStringSubmatch(trackerUrl)[1]
	streamId := regexp.MustCompile(`\/streams\/([0-9]+)`).FindStringSubmatch(trackerUrl)[1]
	formattedString := fmt.Sprintf("%s_%s_%d", channel, streamId, t.Unix())

	if verbose {
		fmt.Printf("%-20s %-20s %-20s %-20s\n%-20s %-20s %-20s %-20s\n",
			"[channel]", "[stream]", "[start]", "[end]",
			channel, streamId, startTimestamp, endTimestamp)
	}

	h := sha1.New()
	h.Write([]byte(formattedString))
	hash := hex.EncodeToString(h.Sum(nil))

	for _, domain := range vodDomains {
		vodUrl := domain + "/" + hash[:20] + "_" + formattedString + "/chunked/index-dvr.m3u8"
		resp, err := http.Get(vodUrl)
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println(vodUrl)
			break
		}
	}
}
