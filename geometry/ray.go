package geometry

type Ray struct {
    Origin Point3;
    Direction Point3;
};

func (self Ray) At(t float64) Point3 {
    return self.Direction.Scale(t).Add(self.Origin);
}
