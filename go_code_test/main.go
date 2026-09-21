package main

func A() {
	B()
}

func B() {
	C()
}

func C() {
	D()
}

func D() {}

func main() {}
