package camera

import (
	"math"
	"math/rand"
	. "raytracer/geometry"
    . "raytracer/utils"
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

const eps = 1e-9;
var reflectance = 0.5;
const samples = 20;
const bounces = 100;
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
        if ok && t >= eps && t <= hit.t {
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

func reflectLambertian(hit HitInfo) Ray {
    random := randPoint().Normalize().Add(hit.normal).Add(hit.point);
    return NewRay(
        hit.point,
        random,
    );
}

func reflectRandom(hit HitInfo) Ray {
    random := randPoint().Normalize();
    if random.Dot(hit.normal) < 0 {
        random.Scale(-1);
    }
    return NewRay(
        hit.point,
        random,
    );
}

func (self Camera) cast(ray Ray, scene []Object, bounces int) Vec3 {
    if bounces <= 0 {
        return NewVec3(0, 0, 0);
    }

    hit := hitWorld(ray, scene);
    if hit.isHit {
        nray := reflectLambertian(hit);
        return self.cast(nray, scene, bounces-1).Scale(reflectance);
    }

    unit := ray.Direction.Normalize();
    p := (unit.Y + 1) * 0.5;
    return cyan.Scale(p).Add(white.Scale(1-p));
}

func (self Camera) getRay(u, v float64) Ray {
    deviX, deviY := rng.Float64() - 0.5, rng.Float64() - 0.5;
    origin := self.pixel00.Add(self.du.Scale(u+deviX)).Add(self.dv.Scale(v+deviY));
    return NewRay(
        origin,
        origin.Sub(self.Center),
    );
}

func (self Camera) Trace(u, v float64, scene []Object) Vec3 {
    var color Vec3;
    reflectance = 0.2 * math.Floor(u/640/0.2) + 0.1;
    for range samples {
        ray := self.getRay(u, v);
        color = color.Add(self.cast(ray, scene, bounces));
    }
    color = GammaCorrect(color.Scale(1/float64(samples)));
    return NewVec3(
        Clamp(color.X, 0, 1),
        Clamp(color.Y, 0, 1),
        Clamp(color.Z, 0, 1),
    );
}
