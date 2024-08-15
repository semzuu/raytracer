package main

import (
	"image"
	"log"
	. "raytracer/geometry"
    . "raytracer/camera"
);

const AspectRatio float64 = 16/9.0;
const ImageWidth = 640;
const ImageHeight = int(ImageWidth/AspectRatio);
const FocalLength = 3;

func main() {
    var scene []Object;
    scene = append(scene, NewSphere(
        NewVec3(-1, 0, 1),
        0.5,
        NewColor(0.7, 0, 0.2),
    ));
    scene = append(scene, NewSphere(
        NewVec3(1, 0, 2),
        1,
        NewColor(0.1, 0, 0.9),
    ));
    scene = append(scene, NewSphere(
        NewVec3(0, -100, 50),
        100,
        NewColor(0, 0.5, 0.2),
    ));

    cameraCenter := NewPoint3(0, 0, -FocalLength);
    camera := NewCamera(
        cameraCenter,
        FocalLength,
        ImageWidth,
        float64(ImageHeight),
    );
    
    output := image.NewNRGBA(image.Rect(0, 0, ImageWidth, ImageHeight));
    for y := 0; y < ImageHeight; y++ {
        log.Printf("Scanlines Remaining: %d\n", ImageHeight-y);
        for x := 0; x < ImageWidth; x++ {
            color := camera.Trace(float64(x), float64(y), scene);
            output.Set(x, y, convertColor(color));
        }
    }
    filepath := "output/anti-aliasing.png";
    err := pngExport(filepath, output);
    if err != nil {
        log.Fatalln(err);
    }
}
