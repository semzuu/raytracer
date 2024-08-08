package geometry

type Object interface{
    Hit(Ray) (float64, bool);
    Center() Point3;
    Color() Vec3;
}
