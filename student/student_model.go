package student

type Student struct {
	Name string `bson:"name"`
	Age  string `bson:"age"`
}
