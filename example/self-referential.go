package main

type Calc struct {
	number float64
	add    func(float64) float64
	minus  func(float64) float64
	times  func(float64) float64
	divide func(float64) float64
}

func newCalc() *Calc {
	var calc Calc
	calc.number = 0
	calc.add = func(x float64) float64 {
		calc.number += x
		return calc.number
	}
	calc.minus = func(x float64) float64 {
		calc.number -= x
		return calc.number
	}
	calc.times = func(x float64) float64 {
		calc.number *= x
		return calc.number
	}
	calc.divide = func(x float64) float64 {
		calc.number /= x
		return calc.number
	}
	return &calc
}

func main() {
	calc := newCalc()
	calc.add(2)
	calc.times(4)
	calc.divide(8)
	calc.minus(2)
	println(calc.number)
}
