package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "imageConvert [-gif | -png | -jpg]")
		os.Exit(1)
	}
	img, kind, err := image.Decode(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot read an image: %v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintln(os.Stderr, "input format =", kind)
	switch os.Args[1] {
	case "-jpg":
		if err := jpeg.Encode(os.Stdout, img, &jpeg.Options{Quality: 95}); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
		}
	case "-png":
		if err := png.Encode(os.Stdout, img); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
		}
	case "-gif":
		if err := gif.Encode(os.Stdout, img, nil); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "imageConvert [-gif | -png | -jpg]")
		os.Exit(1)
	}
}
