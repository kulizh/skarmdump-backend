package page

import (
	"html/template"
	"io"
	"strings"
)

type OGData struct {
	Title    string
	ImgURL   string
	PageURL  string
	SiteName string
}

const ogTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<meta property="og:title" content="{{.Title}}">
<meta property="og:image" content="{{.ImgURL}}">
<meta property="og:url" content="{{.PageURL}}">
<meta property="og:type" content="website">
{{if .SiteName}}<meta property="og:site_name" content="{{.SiteName}}">
{{end}}<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.ImgURL}}">
<meta name="twitter:card" content="photo">
</head>
<body>
<img src="{{.ImgURL}}" alt="{{.Title}}" style="max-width:100%;height:auto;">
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
	return strings.TrimRight(domain, "/") + "/" + hash + ".png"
}
