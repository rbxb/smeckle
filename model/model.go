package model

type Object struct {
	Faces []int
}

type Model struct {
	Name string
	Verts, UVs, Norms []float32
	Objects []Object
}