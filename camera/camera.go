package camera

import (
    . "raytracer/geometry"
    "math"
);

type Camera struct {
    Center Point3;
    FocalLength float64;
    viewportUpperLeft Point3;
    vu, vv Vec3;
    du, dv Vec3;
    pixel00 Point3;
};


func NewCamera(center Point3, focalLength float64, ImageWidth float64, ImageHeight float64) Camera {
    aspectRatio := ImageWidth/ImageHeight;
    viewportHeight := 2.0;
    viewportWidth := viewportHeight*aspectRatio;
    viewportUpperLeft := center.Add(NewVec3(-0.5*viewportWidth, 0.5*viewportHeight, focalLength));
    pixelDelta := viewportWidth/ImageWidth;
    du, dv := NewVec3(pixelDelta, 0, 0), NewVec3(0, -pixelDelta, 0);
    return Camera{
        center,
        focalLength,
        viewportUpperLeft,
        NewVec3(viewportWidth, 0, 0), NewVec3(0, -viewportHeight, 0),
        du, dv,
        viewportUpperLeft.Add(du.Scale(0.5)).Add(dv.Scale(0.5)),
    };
}

func (self Camera) Trace(u, v float64, scene []Object) Vec3 {
    pixelCenter := self.pixel00.Add(self.du.Scale(u)).Add(self.dv.Scale(v));
    ray := NewRay(
        pixelCenter,
        pixelCenter.Sub(self.Center),
    );

    minScalar := math.MaxFloat64;
    index := -1;
    for i, object := range scene {
        t, ok := object.Hit(ray);
        if ok && t >= 0 && t <= minScalar {
            minScalar = t;
            index = i;
        }
    }
    if index == -1 {
        col := 0x50/255.0;
        return NewVec3(col, col, col);
    } else {
        point := ray.At(minScalar)
        normal := point.Sub(scene[index].Center()).Normalize();
        _ = normal;
        return scene[index].Color();
    }
}
