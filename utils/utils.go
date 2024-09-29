package utils

import (
	"image"
	"image/color"
	"image/png"
	"io"
    "math"
	"log"
	"os"
	. "raytracer/geometry"
);

type Color = Vec3;

func NewColor(r, g, b float64) Color {
    return NewVec3(r, g, b);
}

func GammaCorrect(color Color) Color {
    return NewColor(
        math.Sqrt(color.X),
        math.Sqrt(color.Y),
        math.Sqrt(color.Z),
    );
}

func ConvertColor(c Color) color.RGBA {
    c = c.Scale(255);
    return color.RGBA{
        uint8(c.X),
        uint8(c.Y),
        uint8(c.Z),
        255,
    };
}

func Clamp(value, min, max float64) float64 {
    if value > max {
        return max;
    }
    if value < min {
        return min;
    }
    return value;
}

func PngExport(filepath string, image image.Image) error {
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
