/*
Extract prominent colors from images. Go port of the Android Palette class aka Vibrant.js

https://android.googlesource.com/platform/frameworks/support/+/b14fc7c/v7/palette/src/android/support/v7/graphics/

https://github.com/jariz/vibrant.js

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
		fmt.Printf("/* %s (population: %d) *\/\n%s\n\n", name, swatch.Population, swatch)
	}

output:

	/* LightMuted (population: 253) *\/
	.lightmuted{background-color:#cbc0a2;color:#000000;}

	/* DarkMuted (population: 11069) *\/
	.darkmuted{background-color:#5b553f;color:#ffffff;}

	/* Vibrant (population: 108) *\/
	.vibrant{background-color:#dfd013;color:#000000;}

	/* LightVibrant (population: 87) *\/
	.lightvibrant{background-color:#f4ed7d;color:#000000;}

	/* DarkVibrant (population: 2932) *\/
	.darkvibrant{background-color:#917606;color:#ffffff;}

	/* Muted (population: 4098) *\/
	.muted{background-color:#a58850;color:#000000;}

*/
package vibrant
