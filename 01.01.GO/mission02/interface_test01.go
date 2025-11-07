package main

import "fmt"

//题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
//考察点 ：接口的定义与实现、面向对象编程风格。

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

func main() {
	r := Rectangle{width: 5.0, height: 10.0}
	var s Shape = r
	fmt.Printf("Rectangle is area is %f\n", s.Area())
	fmt.Printf("Rectangle is perimeter is %f\n", s.Perimeter())

	c := Circle{radius: 5.0}
	s = c
	fmt.Printf("Circle is area is %f\n", s.Area())
	fmt.Printf("Circle is perimeter is %f\n", s.Perimeter())
}
