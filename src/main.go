package main

import (
	"image"
	"image/png"
	"io"
	"log"
	"os"
);

const ImageWidth = 256;
const ImageHeight = 256;

func main() {
    output := image.NewNRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))
    
    for y := 0; y < ImageHeight; y++ {
        log.Printf("Scanlines Remaining: %d\n", y);
        for x := 0; x < ImageWidth; x++ {
            base := (y-output.Rect.Min.Y)*output.Stride + (x-output.Rect.Min.X)*4;
            output.Pix[base]     = uint8(x);
            output.Pix[base + 1] = uint8(y);
            output.Pix[base + 2] = 0;
            output.Pix[base + 3] = 255;
        }
    }
    filepath := "output/first-image.png";
    err := pngExport(filepath, output);
    if err != nil {
        log.Fatalln(err);
    }
}

func pngExport(filepath string, image image.Image) error {
    log.Printf("Writing to file: %s\n", filepath);
    file, err := os.Create(filepath);
    if err != nil {
        return err;
    }
    writer := io.Writer(file);
    err = png.Encode(writer, image);
    if err != nil {
        return err;
    }
    log.Println("Done");
    return nil;
}
