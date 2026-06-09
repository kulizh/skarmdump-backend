package page

import (
	"html/template"
	"io"
	"strings"
)

type OGData struct {
	Title   string
	ImgURL  string
	PageURL string
}

const ogTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="Share screenshot">
<meta property="og:image" content="{{.ImgURL}}">
<meta property="og:url" content="{{.PageURL}}">
<meta property="og:type" content="image">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="630">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="Share screenshot">
<meta name="twitter:image" content="{{.ImgURL}}">
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
