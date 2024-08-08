package geometry

import "math";

type Sphere struct{
    center Point3;
    radius float64;
    color Vec3;
};

func NewSphere(center Point3, radius float64, color Vec3) Sphere {
    return Sphere{
        center,
        radius,
        color,
    };
}

func (self Sphere) Center() Point3 {
    return self.center;
}

func (self Sphere) Radius() float64 {
    return self.radius;
}

func (self Sphere) Color() Vec3 {
    return self.color;
}

func (self Sphere) Hit(ray Ray) (float64, bool) {
    temp := self.center.Sub(ray.Origin);
    a := ray.Direction.Dot(ray.Direction);
    h := ray.Direction.Dot(temp);
    c := temp.Dot(temp) - self.radius*self.radius;

    delta := (h*h - a*c);
    if 4*delta < 0 {
        return 0, false;
    }
    return (h-math.Sqrt(delta))/a, true;
}
