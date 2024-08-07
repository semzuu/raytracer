package geometry

import "math"

type Vec3 struct{
    X, Y, Z float64;
};

type Point3 = Vec3;

func (self Vec3) Clone() Vec3 {
    return Vec3{
        self.X,
        self.Y,
        self.Z,
    };
}

func (self Vec3) Add(other Vec3) Vec3 {
    return Vec3{
        self.X+other.X,
        self.Y+other.Y,
        self.Z+other.Z,
    };
}

func (self Vec3) Sub(other Vec3) Vec3 {
    return Vec3{
        self.X-other.X,
        self.Y-other.Y,
        self.Z-other.Z,
    };
}

func (self Vec3) Scale(factor float64) Vec3 {
    return Vec3{
        self.X*factor,
        self.Y*factor,
        self.Z*factor,
    };
}

func (self Vec3) Normalize() Vec3 {
    return self.Scale(1/self.Length());
}

func (self Vec3) Dot(other Vec3) float64 {
    return self.X*other.X + self.Y*other.Y + self.Z*other.Z;
}

func (self Vec3) Length() float64 {
    return math.Sqrt(self.Dot(self))
}
