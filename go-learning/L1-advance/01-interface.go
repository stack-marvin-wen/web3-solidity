package main

import (
	"fmt"
)

type Shape interface {
	Area() float32
	Perimeter() float32
}

type Rectangle struct {
	width, height float32
}

func (r Rectangle) Area() float32 {
	return float32(r.width * r.height)
}
func (r Rectangle) Perimeter() float32 {
	return float32(2 * (r.width + r.height))
}

type Circle struct {
	radius float32
}

func (c Circle) Area() float32 {
	return float32(3.14 * c.radius * c.radius)
}
func (c Circle) Perimeter() float32 {
	return float32(2 * 3.14 * c.radius)
}

type Reader interface {
	Read(p []byte) (n int, err error)
}
type Writer interface {
	Write(p []byte) (n int, err error)
}
type ReadWriter interface {
	Reader
	Writer
}
type File struct {
	name string
}

func (f *File) Reader(data []byte) (n int, err error) {
	return 0, nil
}
func (f *File) Write(data []byte) (n int, err error) {
	return 0, nil
}
func emptyInterfaceDemo() {
	var i interface{}
	i = 42
	// 类型断言
	v, ok := i.(int) // 类型断言，检查i是否是int类型
	if ok {
		fmt.Printf("i的值是: %d\n", v)
	} else {
		fmt.Println("i不是int类型")
	}
	// 在switch中使用i.(type)判断i的具体类型
	switch v := i.(type) {
	case int:
		fmt.Printf("i是int类型，值为: %d\n", v)
	case string:
		fmt.Printf("i是string类型，值为: %s\n", v)
	default:
		fmt.Printf("i是其他类型，值为: %v\n", v)
	}

}

func main() {
	var s Shape
	s = Rectangle{10, 20}
	fmt.Printf("Area: %.2f , Perimeter: %.2f\n", s.Area(), s.Perimeter())
	s = Circle{10}
	fmt.Printf("Area: %.2f , Perimeter: %.2f\n", s.Area(), s.Perimeter())
}
