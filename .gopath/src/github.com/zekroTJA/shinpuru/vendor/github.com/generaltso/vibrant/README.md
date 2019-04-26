# vibrant

[![godoc reference](https://godoc.org/github.com/generaltso/vibrant?status.png)](https://godoc.org/github.com/generaltso/vibrant)


Extract prominent colors from images. Go port of the [Android awesome Palette class](https://android.googlesource.com/platform/frameworks/support/+/b14fc7c/v7/palette/src/android/support/v7/graphics/) aka [Vibrant.js](https://github.com/jariz/vibrant.js).

![screenshot of app.go](https://u.teknik.io/Rv3r3.png)

# install

```
go get github.com/generaltso/vibrant
```

# usage

```go
// example: create css stylesheet from image file
checkErr := func(err error) { if err != nil { panic(err) } }

f, err := os.Open("some_image.jpg")
checkErr(err)
defer f.Close()

img, _, err := image.Decode(f)
checkErr(err)

palette, err := vibrant.NewPaletteFromImage(img)
checkErr(err)

for name, swatch := range palette.ExtractAwesome() {
  fmt.Printf("/* %s (population: %d) */\n%s\n\n", name, swatch.Population, swatch)
}
```

###### output:
```css
/* LightMuted (population: 253) */
.lightmuted{background-color:#cbc0a2;color:#000000;}

/* DarkMuted (population: 11069) */
.darkmuted{background-color:#5b553f;color:#ffffff;}

/* Vibrant (population: 108) */
.vibrant{background-color:#dfd013;color:#000000;}

/* LightVibrant (population: 87) */
.lightvibrant{background-color:#f4ed7d;color:#000000;}

/* DarkVibrant (population: 2932) */
.darkvibrant{background-color:#917606;color:#ffffff;}

/* Muted (population: 4098) */
.muted{background-color:#a58850;color:#000000;}

```

See [godoc reference](https://godoc.org/github.com/generaltso/vibrant) for full API.

# bonus round

## reference implementation/command line tool
```
go get github.com/generaltso/vibrant/cmd/vibrant
```

```
usage: vibrant [options] file

options:
  -compress
    	Strip whitespace from output.
  -css
    	Output results in CSS.
  -json
    	Output results in JSON.
  -lowercase
    	Use lowercase only for all output. (default true)
  -rgb
    	Output RGB instead of HTML hex, e.g. #ffffff.
```

## webapp

```
go get github.com/generaltso/vibrant
cd $GOPATH/src/github.com/generaltso/vibrant
go run webapp.go
# open http://localhost:8080/ in a browser
```


# thanks

https://github.com/Infinity/Iris

[This Google I/O 2014 presentation](https://www.youtube.com/watch?v=ctzWKRlTYHQ?t=451)

https://github.com/jariz/vibrant.js

https://github.com/akfish/node-vibrant
