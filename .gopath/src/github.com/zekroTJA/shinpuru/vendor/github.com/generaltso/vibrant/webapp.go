// +build ignore

/*
	sample webapp to demo the functionality of github.com/generaltso/vibrant

	base64 encoding the images provides the benefit of not having to manage
	uploaded files or do any javascript fanciness but comes with a performance penalty

	DO NOT USE THIS IN A PRODUCTION ENVIRONMENT
*/
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"text/template"

	"github.com/generaltso/vibrant"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("<-", r.Method, r.URL)
	setStatus := func(status int) {
		w.WriteHeader(status)
		statusText := ""
		switch status {
		case 200:
			statusText = "OK"
		case 400:
			statusText = "Bad Request"
		case 500:
			statusText = "Internal Server Error"
		}
		log.Println("->", status, statusText)
	}
	w.Header().Set("Content-Type", "text/html")

	data := &struct {
		Response         bool
		Error            string
		Selection        string
		Stylesheet       string
		DataURI          string
		Missing          bool
		Benchmark        time.Duration
		DefaultMaxColors int
		MaxColors        int
	}{DefaultMaxColors: vibrant.DEFAULT_CALCULATE_NUMBER_COLORS,
		MaxColors: vibrant.DEFAULT_CALCULATE_NUMBER_COLORS}

	defer func() {
		t, err := template.New("").Parse(tpl)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		if err := t.Execute(buf, data); err != nil {
			panic(err)
		}
		io.Copy(w, buf)
	}()

	if r.Method != "POST" {
		setStatus(200)
		return
	}

	file, header, err := r.FormFile("test")
	if err != nil {
		setStatus(400)
		data.Error = "400 Bad Request"
		return
	}

	switch header.Header["Content-Type"][0] {
	case "image/jpeg":
	case "image/png":
	default:
		setStatus(400)
		data.Error = "JPG/PNG only plz."
		return
	}

	if max := r.FormValue("maxColors"); max != "" {
		n, err := strconv.Atoi(max)
		if err != nil || n < 1 {
			setStatus(500)
			data.Error = err.Error()
			return
		}
		data.MaxColors = n
	}

	img, _, err := image.Decode(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = err.Error()
		return
	}
	file.Seek(0, 0)
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	file.Close()
	data.DataURI = fmt.Sprintf("data:%s;base64,%s", header.Header["Content-Type"][0], base64.StdEncoding.EncodeToString(buf.Bytes()))

	start := time.Now()
	palette, err := vibrant.NewPalette(img, data.MaxColors)
	data.Benchmark = time.Since(start)
	if err != nil {
		setStatus(500)
		data.Error = err.Error()
		return
	}

	awesome := palette.ExtractAwesome()
	stylesheet := ""
	for _, sw := range awesome {
		stylesheet += fmt.Sprintf("%s\n", sw)
		if sw.Name == "Vibrant" {
			vendorPrefixingIsAWESOME := fmt.Sprintf("{\n    background-color: %s;\n   color: %s;\n}\n", sw.Color.RGBHex(), sw.Color.TitleTextColor())
			data.Selection = fmt.Sprintf("::selection %s::-moz-selection %s::-webkit-selection %s", vendorPrefixingIsAWESOME, vendorPrefixingIsAWESOME, vendorPrefixingIsAWESOME)
		}
	}
	data.Stylesheet = stylesheet
	data.Missing = len(awesome) != 6
	setStatus(200)
	data.Response = true
	return
}

func main() {
	http.HandleFunc("/", index)
	log.Println("Listening on :8080...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

const tpl = `
<!doctype html>
<html>
    <head>
        <meta charset='utf-8'>
        <title>Go Vibrant!</title>
        <style>
* {
    box-sizing: border-box;
    margin: 0;
}
*:focus {
    outline: none;
}
*::-moz-focus-inner, *::-moz-focus-outer {
    border: 0;
}
body {
    font-family: sans-serif;
    max-width: 800px;
    margin: auto;
}
h3 {
    color: #800;
}
hr {
    margin: 1em 0;
    border: 0;
    border-top: 1px ridge;
}
h1, form {
    display: inline-block;
}
form {
    margin-left: 1em;
}
button {
    background: #9cf;
    color: #fff;
    border: 0;
    padding: .5em 1em;
    border-radius: 3px;
    text-transform: uppercase;
    letter-spacing: .5px;
    box-shadow: 0 0 5px #ccc;
    transition: all 125ms cubic-bezier(.8,0,.2,1);
    cursor: pointer;
    z-index: 1;
    position: relative;
}
button:hover {
    box-shadow: none;
}
button:active {
    box-shadow: 0 0 0 100vmax rgba(153,204,255,1);
}
input[type="text"] {
    margin-right: 1em;
}
        </style>
    </head>
    <body>
        <h1>choose an image:</h1>
        <form action='/' method='post' enctype='multipart/form-data'>
            <input type='file' name='test' accept='image/*' required>
            <input type='number' min='1' size='11' maxlength='10' name='maxColors' value='{{.MaxColors}}' title='maxColors (default: {{.DefaultMaxColors}})' placeholder='maxColors (default: {{.DefaultMaxColors}})'>
            <button type='submit' name='vibrant' value='q'>Go Vibrant!</button>
        </form>
        {{ if .Error }}
        <hr>
        <h3>Error: {{.Error}}</h3>
        {{ end }}
        {{ if .Response }}
        <hr>
        <style>
figure {
    box-shadow: 0 0 10px #ccc;
}
img {
    max-width: 100%;
    display: block;
    margin: auto;
}
figcaption {
    display: flex;
    flex-wrap: wrap;
}
figcaption div {
    flex: 1 1 16.6667%;
    padding: 3vw 0;
    text-align: center;
    display: inline-block;
}
textarea {
    margin: 1em auto;
    width: 100%;
    display: block;
    height: 12em;
    resize: none;
    border: 0;
}
h2 {
    text-align: center;
}
{{.Selection}}
{{.Stylesheet}}
        </style>
        <figure>
            <img src='{{.DataURI}}' alt='something happened'>
            <figcaption>
                <div class="vibrant">Vibrant</div>
                <div class="lightvibrant">LightVibrant</div>
                <div class="darkvibrant">DarkVibrant</div>
                <div class="muted">Muted</div>
                <div class="lightmuted">LightMuted</div>
                <div class="darkmuted">DarkMuted</div>
            </figcaption>
        </figure>
        <textarea readonly onclick='this.select()'>{{.Stylesheet}}</textarea>
        {{ if .Missing }}
        <h3>If color swatches are missing, try increasing <code>maxColors</code> in the text field above.</h3>
        {{ end }}
        <h2>{{.Benchmark}}</h2>
        {{ end }}
    </body>
</html>`
