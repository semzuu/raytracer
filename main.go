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
var LightSource = NewSphere(Vec3{-2, 1, 2}, 0.1, Color{1, 1, 1});

func rayColor(ray Ray, objects []Object) Color {
    minScalar := math.MaxFloat64;
    index := -1;
    for i, object := range objects {
        t, ok := object.Hit(ray);
        if ok && t >= 0 && t <= minScalar {
            minScalar = t;
            index = i;
        }
    }
    if index == -1 {
        _, ok := LightSource.Hit(ray);
        if ok {
            return LightSource.Color();
        }
        col := 0x50/255.0;
        return Color{col, col, col};
    } else {
        point := ray.At(minScalar)
        normal := point.Sub(objects[index].Center()).Normalize();
        light := LightSource.Center().Sub(point).Add(point).Normalize();
        sim := light.Dot(normal);
        if sim <= 0 {
            mult := (sim+1)*0.5;
            return objects[index].Color().Scale(mult);
        }
        return objects[index].Color();
    }
}

func main() {
    output := image.NewNRGBA(image.Rect(0, 0, ImageWidth, ImageHeight));
    var objects []Object;
    objects = append(objects, NewSphere(
        Vec3{-1, 0, 1},
        0.5,
        Color{0.7, 0, 0.2},
    ));
    objects = append(objects, NewSphere(
        Vec3{1, 0, 2},
        1,
        Color{0.1, 0, 0.9},
    ));
    objects = append(objects, NewSphere(
        Vec3{0, -100, 50},
        100,
        Color{0, 0.5, 0.2},
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
    filepath := "output/shaded-spheres.png";
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
