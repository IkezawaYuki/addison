package main

import "fmt"

func main() {
	i := 2

	switch i {
	case 1:
		fmt.Println("i is 1")
	case 2:
	case 3:
		fmt.Println("i is 2 or 3")
	default:
		fmt.Println("i is something else")
	}
	//if len(os.Args) != 2 {
	//	_, _ = fmt.Fprintf(os.Stderr, "imageConvert [-gif | -png | -jpg]")
	//	os.Exit(1)
	//}
	//img, kind, err := image.Decode(os.Stdin)
	//if err != nil {
	//	_, _ = fmt.Fprintf(os.Stderr, "cannot read an image: %v\n", err)
	//	os.Exit(1)
	//}
	//_, _ = fmt.Fprintln(os.Stderr, "input format =", kind)
	//switch os.Args[1] {
	//case "-jpg":
	//	if err := jpeg.Encode(os.Stdout, img, &jpeg.Options{Quality: 95}); err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
	//	}
	//case "-png":
	//	if err := png.Encode(os.Stdout, img); err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
	//	}
	//case "-gif":
	//	if err := gif.Encode(os.Stdout, img, nil); err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "conversion error: %v\n", err)
	//	}
	//default:
	//	_, _ = fmt.Fprintf(os.Stderr, "imageConvert [-gif | -png | -jpg]")
	//	os.Exit(1)
	//}
}
