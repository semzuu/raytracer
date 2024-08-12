package camera

import (
	"math"
	"math/rand"
	. "raytracer/geometry"
	"time"
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

const bounces = 200;
var rng = rand.New(rand.NewSource(time.Now().Unix()));
var white = NewVec3(1, 1, 1);
var cyan = NewVec3(0.5, 0.7, 1);
type HitInfo struct {
    object *Object;
    isHit bool;
    t float64;
    normal Vec3;
    point Vec3;
};

func randPoint() Vec3 {
    min := -1.0;
    max := 1.0;
    return NewVec3(
        rng.Float64()*(max-min) + min,
        rng.Float64()*(max-min) + min,
        rng.Float64()*(max-min) + min,
    );
}

func hitWorld(ray Ray, scene []Object) HitInfo {
    var hit HitInfo;
    hit.t = math.MaxFloat64;
    for _, object := range scene {
        t, ok := object.Hit(ray);
        if ok && t >= 0 && t <= hit.t {
            hit.t = t;
            hit.isHit = ok;
            hit.object = &object;
        }
    }
    if hit.isHit {
        obj := *(hit.object);
        hit.point = ray.At(hit.t)
        hit.normal = hit.point.Sub(obj.Center()).Normalize();
    }
    return hit;
}

func (self Camera) cast(ray Ray, scene []Object, bounces int) Vec3 {
    if bounces <= 0 {
        return NewVec3(0, 0, 0);
    }

    hit := hitWorld(ray, scene);
    if hit.isHit {
        random := randPoint().Normalize();
        if random.Dot(hit.normal) < 0 {
            random.Scale(-1);
        }
        nray := NewRay(
            hit.point,
            random,
        );
        return self.cast(nray, scene, bounces-1).Scale(0.5);
    }
    return cyan;
}

func (self Camera) Trace(u, v float64, scene []Object) Vec3 {
    pixelCenter := self.pixel00.Add(self.du.Scale(u)).Add(self.dv.Scale(v));
    ray := NewRay(
        pixelCenter,
        pixelCenter.Sub(self.Center),
    );

    return self.cast(ray, scene, bounces);
}
