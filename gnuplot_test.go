// Copyright 2012, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gnuplot

import (
    "testing"
    "os"
    "math"
    "math/rand"
    "fmt"
    "time"
)

func TestPlot(t *testing.T) {
    scriptFile := "test.plt"
    imgFile := "abc.svg"
    f, err := os.Create(scriptFile)
    if err != nil {
        panic(err)
    }
    defer func () {
        f.Close()
        os.Remove(scriptFile)
    }()

    p := NewGnuplot(f, imgFile, 600, 600)
    p.Xrange(-3.9, 3.9)
    p.Yrange(-3.9, 3.9)
    p.Margin(3)
    p.LockRatio()
    p.Multiplot(func() {
        p.Circle(0, 0, 3)
        p.With("parametric", func() {
            p.Plot(func() {
                lineType := p.LineType
                p.LineType = DOTTED_LINE
                p.Lines("t, 0", "0, t")
                p.LineType = lineType
                
                points := make([]string, 10)
                lines := make([]string, len(points))
                for i := range points {
                    c := PolarPoint(
                        2*math.Pi*rand.Float64(),
                        3*math.Sqrt(rand.Float64())) // Uniform distribution on the circle

                    points[i] = c.String()
                    lines[i] = Line{c}.AxPlusB()
                }    

                p.Lines(lines...)
                p.Points(points...)
            })
        })
    })

    p.Quit()

    RunPlot(scriptFile)
}

type Point struct {
    X, Y float64
}

func PolarPoint(theta, radius float64) Point {
    return Point{radius*math.Cos(theta), radius*math.Sin(theta)}
}

func (p Point) String() string {
    return fmt.Sprint(p.X, p.Y)
}


type Line struct {
    C Point // normal vector and also the point the line passes
}

func (l Line) AxPlusB() string {
    p := l.C
    x, y := p.X, p.Y
    rr := x*x + y*y
    return fmt.Sprintf("t,%g*t%+g", -x/y, rr/y)
}

func init() {
    rand.Seed(time.Now().UnixNano())
}