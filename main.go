package main

import (
	"image"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	. "raytracer/geometry"
);

type Color = Vec3;

const AspectRatio float64 = 16/9.0;
const ImageWidth = 640;
const ImageHeight = int(ImageWidth/AspectRatio);
const ViewportHeight float64 = 2;
const ViewportWidth = ViewportHeight*AspectRatio;
const FocalLength = 3;

func rayColor(ray Ray, objects []Object) Color {
    minScalar := math.MaxFloat64;
    index := -1;
    for i, object := range objects {
        t, ok := object.Hit(ray);
        if ok && t <= minScalar {
            minScalar = t;
            index = i;
        }
    }
    if index == -1 {
        col := 0x18/255.0;
        return Color{col, col, col};
    } else {
        normal := ray.At(minScalar).Sub(objects[index].Center()).Normalize();
        return normal.Add(Vec3{1, 1, 1}).Scale(0.5);
    }
}

func main() {
    output := image.NewNRGBA(image.Rect(0, 0, ImageWidth, ImageHeight));
    var objects []Object;
    objects = append(objects, NewSphere(
        Vec3{-1, 0, 2},
        0.5,
    ));
    objects = append(objects, NewSphere(
        Vec3{0, -100, 50},
        100,
    ));
    objects = append(objects, NewSphere(
        Vec3{1, 1, 2},
        1,
    ));

    cameraCenter := Point3{0, 0, -FocalLength};
    Vu, Vv := Vec3{ViewportWidth, 0, 0}, Vec3{0, -ViewportHeight, 0};
    Du, Dv := Vu.Scale(1/float64(ImageWidth)), Vv.Scale(1/float64(ImageHeight));
    viewportUpperLeft := cameraCenter.Add(Vec3{ViewportWidth*-0.5, ViewportHeight*0.5, FocalLength});
    pixel00 := viewportUpperLeft.Add(Du.Scale(0.5)).Add(Dv.Scale(0.5));
    
    for y := 0; y < ImageHeight; y++ {
        log.Printf("Scanlines Remaining: %d\n", ImageHeight-y);
        for x := 0; x < ImageWidth; x++ {
            pixelCenter := pixel00.Add(Du.Scale(float64(x))).Add(Dv.Scale(float64(y)));
            ray := Ray{
                pixelCenter,
                pixelCenter.Sub(cameraCenter),
            };

            color := rayColor(ray, objects).Scale(255);
            base := (y-output.Rect.Min.Y)*output.Stride + (x-output.Rect.Min.X)*4;
            output.Pix[base]     = uint8(color.X);
            output.Pix[base + 1] = uint8(color.Y);
            output.Pix[base + 2] = uint8(color.Z);
            output.Pix[base + 3] = 255;
        }
    }
    filepath := "output/multiple-circles.png";
    err := pngExport(filepath, output);
    if err != nil {
        log.Fatalln(err);
    }
}

func pngExport(filepath string, image image.Image) error {
    log.Printf("Writing to file: %s\n", filepath);
    file, err := os.Create(filepath);
    defer file.Close();
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
