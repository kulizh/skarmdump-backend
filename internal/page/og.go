package page

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

type OGData struct {
	Title    string
	ImgURL   string
	PageURL  string
	SiteName string
	ReImgURL string
}

const ogTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<meta property="og:title" content="{{.Title}}">
<meta property="og:image" content="{{.ImgURL}}">
<meta property="og:image:type" content="image/png">
<meta property="og:url" content="{{.PageURL}}">
<meta property="og:type" content="website">
{{if .SiteName}}<meta property="og:site_name" content="{{.SiteName}}">
{{end}}<meta name="twitter:title" content="{{.Title}}"><meta name="twitter:image" content="{{.ImgURL}}">
<meta name="twitter:card" content="summary_large_image">
</head>
<body>
<img src="{{.ReImgURL}}" alt="{{.Title}}" style="max-width:100%;height:auto;">
</body>
</html>`

var tpl = template.Must(template.New("og").Parse(ogTemplate))

func RenderOGPage(w io.Writer, data OGData) error {
	return tpl.Execute(w, data)
}

func OGPageURL(domain, hash string) string {
	return strings.TrimRight(domain, "/") + "/" + hash
}

func OGImgURL(domain, hash string) string {
	return fmt.Sprintf("%s/%s.og.png", strings.TrimRight(domain, "/"), hash)
}
