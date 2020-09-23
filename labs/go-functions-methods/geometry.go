// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//Point defined as a structure with two float values
type Point struct{ x, y float64 }

//Distance traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

//Distance same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.
type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

//X returns y value of Point
func (p Point) X() float64 {
	return p.x
}

//Y returns y value of Point
func (p Point) Y() float64 {
	return p.y
}
func isSegment(p, q, r Point) bool {
	if (q.X() <= math.Max(p.X(), r.X())) &&
		(q.X() >= math.Min(p.X(), r.X())) &&
		(q.Y() <= math.Max(q.Y(), r.Y())) &&
		(q.Y() >= math.Min(p.Y(), r.Y())) {
		return true
	}
	return false
}

func direction(p, q, r Point) int {
	dir := (q.Y()-p.Y())*(r.X()-q.X()) - (q.X()-p.X())*(r.Y()-q.Y())
	if dir == 0 {
		return 0
	}
	if dir > 0 {
		return 1
	}
	return 2
}

func (path Path) intersection() bool {
	if len(path) == 3 {
		return false
	}
	itdoes := false
	for i := 0; i < len(path)-3; i++ {
		for j := 0; j < i+1; j++ {
			a := direction(path[j], path[j+1], path[len(path)-2])
			b := direction(path[j], path[j+1], path[len(path)-1])
			c := direction(path[len(path)-2], path[len(path)-1], path[j])
			d := direction(path[len(path)-2], path[len(path)-1], path[j+1])
			itdoes = false
			if ((a != b) && (c != d)) ||
				((a == 0) && isSegment(path[j], path[len(path)-2], path[j+1])) ||
				((b == 0) && isSegment(path[j], path[len(path)-1], path[j+1])) ||
				((c == 0) && isSegment(path[len(path)-2], path[j], path[len(path)-1])) ||
				((d == 0) && isSegment(path[len(path)-2], path[j+1], path[len(path)-1])) {
				itdoes = true
			}
		}
	}
	return itdoes
}

func random() float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	return math.Round(((random.Float64()*200)-100)*100) / 100
}

func createFigure(sides int) Path {
	figure := Path{}
	for i := 0; i < sides; i++ {
		figure = append(figure, Point{random(), random()})
		for figure.intersection() {
			figure[i] = Point{random(), random()}
		}
	}
	return figure
}

func main() {
	if len(os.Args) > 0 {
		sides, _ := strconv.Atoi(os.Args[1])
		if sides < 3 {
			fmt.Println("WRONG INPUT: Number of sides must be over 3 to make a figure")
			return
		}
		figure := createFigure(sides)
		fmt.Printf("- Generating a [%d] sides figure\n- Figure's vertices\n", sides)
		for i := 0; i < sides; i++ {
			fmt.Printf("\t- ( %.1f, %.1f)\n", figure[i].X(), figure[i].Y())
		}
		fmt.Println("- Figure's Perimeter")
		var per = 0.0
		fmt.Print("\t- ")
		i := 0
		for ; i < sides-1; i++ {
			per = per + figure[i].Distance(figure[(i+1)%sides])
			fmt.Printf("%.2f + ", figure[i].Distance(figure[(i+1)%sides]))
		}
		per = per + figure[i].Distance(figure[(i)%sides])
		fmt.Printf("%.2f = %.2f\n", figure[i].Distance(figure[(i)%sides]), per)
	} else {
		fmt.Println("No argument introduced")
	}
	return
}

//!-path
