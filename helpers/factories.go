package helpers

import (
	"time"
	"regexp"
	"net/url"
	"strings"
	"golang.org/x/net/html"
	"github.com/uniplaces/carbon"
)

type DeltaLog map[string]interface{}

func LogFactory(n *html.Node, tp string) DeltaLog {
	log := map[string]interface{}{
		"log_type": tp,
		"name":     "",
		"title":    "",
		"time":     "",
		"duration": LogDuration(n),
	}

	if name := NodeValue(n, "name"); len(name) > 0 {
		log["name"] = name
	} else {
		log["name"] = NodeText(n, 50)
	}

	if title := NodeValue(n, "title"); len(title) > 0 {
		log["title"] = title
	}

	if tm := NodeValue(n, "time"); len(tm) > 0 {
		log["time"] = LogTime(tm)
	}

	return log
}

func LogAppFactory(n *html.Node) DeltaLog {
	log := LogFactory(n, "app")

	log["title"] = NodeText(n, 150)

	return log
}

func LogUrlFactory(n *html.Node) DeltaLog {
	log := LogFactory(n, "url")
	//log["title"] = NodeText(n, 150)

	href := NodeValue(n, "href")

	if isIp, _ := regexp.MatchString("^\\d{1,3}\\.", href); false == isIp {
		if "http" != href[:4] {
			href = "http://" + href
		}

		if lnk, err := url.Parse(href); err == nil {
			log["name"] = strings.Replace(lnk.Host, "www.", "", 1)
			href = lnk.String()

			if sq := DetectSearchQuery(href); len(sq) > 0 {
				log["name"] = sq
				log["log_type"] = "search-query"
			}
		}
	}

	return log
}

func LogGpsFactory(n *html.Node) DeltaLog {
	log := LogFactory(n, "gps-point")

	lng := NodeValue(n, "longitude")
	lat := NodeValue(n, "latitude")

	if lng != "" && lat != "" {
		log["name"] = lat + "," + lng
	}

	return log
}

func LogIdleFactory(n *html.Node) DeltaLog {
	return LogFactory(n, "idle")
}

func LogFolderFactory(n *html.Node) DeltaLog {
	log := LogFactory(n, "folder")
	log["name"] = NodeText(n, 0)

	return log
}

func LogSmsFactory(n *html.Node, lt string) DeltaLog {
	log := LogFactory(n, lt)
	val := NodeText(n, 0)
	reg, _ := regexp.Compile(`^sms (?:to|from)(.+)$`)
	result := reg.FindStringSubmatch(val)[1]

	parts := strings.Split(strings.Trim(result, " "), ":")

	var phone, name string
	name = "Unknown contact"

	if len(parts) > 1 && len(parts[0]) > 0 {
		name = StripNewLines(parts[0])
	}

	if len(parts) > 1 {
		phone = StripNewLines(parts[1])
	}

	log["name"] = name
	log["title"] = phone

	return log
}

func LogCallFactory(n *html.Node, lt string) DeltaLog {
	log := LogFactory(n, lt)

	val := NodeText(n, 0)
	parts := strings.Split(val, ":")

	name := "Unknown contact"

	if len(parts) > 1 {
		name = StripNewLines(
			strings.Trim(parts[1], " "),
		)
	}

	log["name"] = name

	return log
}

func LogChatFactory(n *html.Node) DeltaLog {
	log := LogFactory(n, "chat")

	var name, title string

	msg := strings.Trim(NodeText(n, 0), " ")
	parts := strings.Split(msg, ":")

	name = "Unknown contact"
	if len(parts) > 1 && len(parts[0]) > 0 {
		name = strings.Trim(parts[0], " ")
	}
	title = strings.Trim(strings.Join(parts, ", "), ", ")

	log["name"] = name
	log["title"] = title

	return log
}

func LogKeystrokeFactory(n *html.Node) DeltaLog {
	log := LogFactory(n, "keystrokes")
	txt := NodeText(n, 0)

	log["title"] = txt
	log["duration"] = float64(len(txt))

	return log
}

func LogMessageFactory(n *html.Node, lt string) DeltaLog {
	log := LogFactory(n, lt)
	log["name"] = NodeText(n, 0)

	return log
}

func LogTime(timeStr string) time.Time {
	t, err := time.Parse("15:04:05", timeStr)

	if err != nil {
		t, _ = time.Parse("15:04", timeStr)
	}

	logTime, _ := carbon.CreateFromTime(t.Hour(), t.Minute(), t.Second(), 0, "")

	return logTime.Time
}

func StripNewLines(name string) string {
	name = strings.Replace(name, "\n\r", "", -1)
	name = strings.Replace(name, "\n", "", -1)

	return name
}

func LogDuration(n *html.Node) float64 {
	var d float64

	if ds := NodeValue(n, "dur"); len(ds) > 0 {
		dur, _ := time.ParseDuration(ds + "s")
		d = dur.Seconds()
	}

	return d
}
