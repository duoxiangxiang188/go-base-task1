package main

import (
	"fmt"
	"math"
)

// 形状接口
type Shape interface {
	Area() float64      //面积
	Perimeter() float64 //周长

}

// 矩形结构体
type Rectangle struct {
	Length float64
	Width  float64
	Name   string
}

// 圆形结构体
type Circle struct {
	Radius float64
	Name   string
}

func (r Rectangle) Area() float64 {
	return r.Length * r.Width
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Length + r.Width)
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func main() {
	rect := Rectangle{Length: 2, Width: 2, Name: "矩形"}
	circle := Circle{Radius: 2, Name: "圆形"}
	shapes := []Shape{rect, circle}
	for i, s := range shapes {
		var name string
		switch v := s.(type) {
		case Rectangle:
			name = v.Name
		case Circle:
			name = v.Name
		default:
			name = "位置形状"
		}
		fmt.Printf("【形状 %d  %s】 面积: %.2f ,周长: %.2f \n", i+1, name, s.Area(), s.Perimeter())
	}
}
