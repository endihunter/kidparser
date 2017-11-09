package helpers

import (
	"regexp"
	"strings"
	"net/url"
	"golang.org/x/net/html"
	"github.com/thoas/go-funk"
)

func Parse(n *html.Node, deviceID int, log []DeltaLog) []DeltaLog {
	allowed := []string{
		"app", "url", "in_sms", "out_sms", "in_call", "out_call", "chat", "gps-point", "folder", "keystrokes", "idle",
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data != "p" {
			log = Parse(c, deviceID, log)
			continue
		}

		if c.Type == html.ElementNode && c.Data == "p" {
			if cls := NodeClass(c); cls != "" && funk.Contains(allowed, cls) {
				var item DeltaLog

				switch cls {
				case "app":
					item = LogAppFactory(c)
				case "folder":
					item = LogFolderFactory(c)
				case "in_sms", "out_sms":
					item = LogSmsFactory(c, cls)
				case "in_call", "out_call":
					item = LogCallFactory(c, cls)
				case "chat":
					item = LogChatFactory(c)
				case "url":
					item = LogUrlFactory(c)
				case "gps-point":
					item = LogGpsFactory(c)
				case "keystrokes":
					item = LogKeystrokeFactory(c)
				case "clipboard", "system":
					item = LogMessageFactory(c, cls)
				case "idle":
					item = LogIdleFactory(c)
				}

				item["device_id"] = deviceID

				log = append(log, item)
			}
		}
	}

	return log
}

func DetectSearchQuery(href string) string {
	engines := map[string][]string{
		"Google":     {"google", "q"},
		"Bing":       {"bing.com/search", "q"},
		"Yahoo":      {"search.yahoo.", "p"},
		"Youtube":    {"youtube.", "search_query"},
		"Amazon":     {"amazon.", "field-keywords"},
		"Ebay":       {"ebay.", "_nkw"},
		"Facebook":   {"facebook.com/search", "q"},
		"BBC":        {"bbc.co.uk/search", "q"}, // no query string, search phrase is in url: bbc.co.uk/s
		"CNN":        {"cnn.com/search", "q"},
		"Wikipedia":  {"wikipedia.org/wiki", "*"}, // no query string, search phrase is in url: wikipedia
		"CraigsList": {"craigslist.org/search", "query"},
		"AOL":        {"aol.com/aol/search", "q"},
		"ASK.com":    {"ask.com/web", "q"},
		"Lycos":      {"search.lycos.com", "query"},
		"Webcrawler": {"webcrawler.com/search/web", "q"},
		"Info.com":   {"info.com/searchw", "qkw"},
		"Mahalo.com": {"mahalo.com/search", "q"},
		"AllExperts": {"allexperts.com/sitesearch", "terms"},
		"MSN.com":    {"msnbc.msn.com", "q"},
		"Usatoday":   {"usatoday.com/search/results", "q"},
	}

	for _, o := range engines {
		subUrl, param := o[0], o[1]

		if strings.Contains(href, subUrl) {
			urlParts, err := url.Parse(href)
			if err != nil {
				return ""
			}

			args := urlParts.Query()

			if "*" != param {
				if val := args.Get(param); len(val) > 0 {
					return strings.Trim(val, " ")
				}
			} else {
				if reg, err := regexp.Compile(subUrl + "/(.+)"); err == nil {
					result := reg.FindStringSubmatch(href)[1]

					return strings.Trim(result, "")
				}
			}
		}
	}

	return ""
}
