package helpers

import (
	"io"
	"bytes"
	"strings"
	"golang.org/x/net/html"
	//"github.com/JalfResi/GoTidy"
	"github.com/microcosm-cc/bluemonday"
)

func NodeValue(node *html.Node, key string) string {
	attributes := node.Attr
	for i := 0; i < len(attributes); i++ {
		if attributes[i].Key == key {
			return attributes[i].Val
		}
	}

	return ""
}

func NodeClass(node *html.Node) string {
	return NodeValue(node, "class")
}

func NodeText(node *html.Node, max int) string {
	cnt := node.FirstChild.Data

	sb := node.FirstChild.NextSibling
	if sb != nil && "img" == sb.Data {
		cnt += NodeContent(sb)
	}

	if max > 0 && len(cnt) > max {
		cnt = cnt[:max]
	}

	return strings.Trim(cnt, " ")
}

func NodeContent(node *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, node)

	return buf.String()
}

func StripTags(content string) string {
	st := bluemonday.NewPolicy()
	st.AllowImages()
	st.AllowAttrs("class", "name", "longitude", "latitude", "href", "dur", "time", "move", "title").OnElements("p")
	st.AllowAttrs("src").OnElements("img")

	return st.Sanitize(content)
}

// @todo: Fix fucking Tidy cleaner
func CleanContentWithTidy(fileContent string) (string, error) {
	return fileContent, nil

	//t := tidy.New()
	//defer t.Free()
	//
	//t.OutputXml(true)
	//t.Clean(true)
	//t.HideComments(true)
	//t.FixBadComments(true)
	//t.DropEmptyParas(true)
	//t.DropFontTags(true)
	//t.FixUri(true)
	//t.ShowBodyOnly(tidy.True)
	//t.OutputEncoding(tidy.Utf8)
	//t.CharEncoding(tidy.Utf8)
	//tidyContent, err := t.Tidy(fileContent)
	//
	//return tidyContent, err
}
