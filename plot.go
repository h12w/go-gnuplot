// Copyright 2012, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gnuplot

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"
    "path/filepath"
    "math"
)

var GNUPLOT = `gnuplot`

const (
    // line types
    DOTTED_LINE = "3"

    // Point type
    CIRCLE_POINT = "6"
    ROUND_POINT  = "7"
)

type Gnuplot struct {
    f           io.Writer
    indentLevel int
    Style
}

type Style struct {
    Color      string
    LineType   string
    LineWidth  float64
    PointType  string
    PointSize  float64
}

func NewGnuplot(f io.Writer, imgFile string, width, height int) *Gnuplot {
    p := &Gnuplot{f, 0, 
        Style{
            Color:    "-1",
            LineType: "-1",
            LineWidth: 0.5,
            PointType: "7",
            PointSize: 0.5}}

    p.setOutput(imgFile, width, height, `background rgb '#FFFFFF' font 'Cambria Math,8'`)
    p.Set("termoption dashed")
    return p
}

func (p *Gnuplot) PointStyle() string {
    return fmt.Sprintf("pointtype %s pointsize %g linecolor %s", p.PointType, p.PointSize, p.Color)
}

func (p *Gnuplot) LineStyle() string {
    return fmt.Sprintf("linetype %s linewidth %g linecolor %s", p.LineType, p.LineWidth, p.Color)
}

func (p *Gnuplot) Emit(cmd string) {
    _, err := io.WriteString(p.f, strings.Repeat("    ", p.indentLevel) + cmd + "\n")
    if err != nil {
        panic(err)
    }
}

func (p *Gnuplot) Plot(f func()) {
    p.Emit(`plot \`)
    p.indentLevel++
    f()
    p.Empty()
    p.indentLevel--
}

func (p *Gnuplot) Set(cmd string) {
    p.Emit("set " + cmd)
}

func (p *Gnuplot) Unset(cmd string) {
    p.Emit("unset " + cmd)
}

func (p *Gnuplot) data(style string, f func()) {
    p.Emit(`'-' ` + style)
    p.indentLevel++
    f()
    p.indentLevel--
    p.Emit(`e,\`)
}

func (p *Gnuplot) Margin(m int) {
    p.Set(fmt.Sprintf("rmargin %d", m))
    p.Set(fmt.Sprintf("lmargin %d", m))
    p.Set(fmt.Sprintf("tmargin %d", m))
    p.Set(fmt.Sprintf("bmargin %d", m))
}

// HACK: Plot an invisible dot
func (p *Gnuplot) Empty() {
    p.Emit(`'-' linetype bgnd`)
    p.Emit(fmt.Sprint(math.MaxFloat64, math.MaxFloat64))
    p.Emit(`e`)
}

func (p *Gnuplot) PlotBorder() {
    // plot border only once
    p.Set("tics")
    p.Set("border")
    // draw the border
    p.Plot(func () {})
}

func (p *Gnuplot) With(option string, f func()) {
    p.Set(option)
    p.indentLevel++
    f()
    p.indentLevel--
    p.Unset(option)
}

func (p *Gnuplot) Multiplot(f func()) {
    p.With("multiplot", func() {
        p.NoTitle()
        p.Unset("tics")
        p.Unset("border")
        p.Unset("raxis")
        f()
        p.PlotBorder()
    })
}

func (p *Gnuplot) Points(points ...string) {
    for _, point := range points {
        p.Emit(`"<echo ` + point + `" with points ` + p.PointStyle() + `,\`)
    }
}

// Plot y = ax + b
func (p *Gnuplot) Lines(lines ...string) {
    for _, line := range lines {
        p.Emit(line + " with lines " + p.LineStyle() + `,\`)        
    }
}

func (p *Gnuplot) Circle(a, b, r float64) {
    p.Set(fmt.Sprintf("object circle at %g,%g size %g fillcolor rgb '#000000' fillstyle empty linewidth 0.2", a, b, r))
}

func (p *Gnuplot) setOutput(file string, width, height int, style string) {
    ext := filepath.Ext(file)
    switch(ext) {
    case ".svg":
        p.setTerminal("svg", width, height, style)
    case ".png":
        p.setTerminal("pngcairo", width, height, style)
    case ".htm", ".html":
        p.setTerminal("canvas", width, height, style)
    // Ugly formats that is worthless supporting
    // case ".pdf":
    //     p.setTerminal("pdfcairo", width, height, style)
    // case ".emf":
    //     p.setTerminal("emf", width, height, style)
    // case ".jpg", ".jpeg":
    //     p.setTerminal("jpeg", width, height, style)
    default:
        fmt.Println("Unsupported Format", ext)
        os.Exit(1)
    }
    p.Emit(fmt.Sprintf("set output '%s'", file))
}

func (p *Gnuplot) setTerminal(name string, width, height int, style string) {
    p.Emit(fmt.Sprintf("set terminal %s size %d,%d %s", name, width, height, style))
    
}

func (p *Gnuplot) NoTitle() {
    p.Unset("key")
}

func (p * Gnuplot) LockRatio() {
    p.Set("size square")
}

func (p * Gnuplot) Xtics(style string) {
    p.Set("xtics " + style)
}

func (p * Gnuplot) Ytics(style string) {
    p.Set("ytics " + style)
}

func (p * Gnuplot) Xrange(from, to float64) {
    p.Set(fmt.Sprintf("xrange [%g:%g]", from, to))
}

func (p * Gnuplot) Yrange(from, to float64) {
    p.Set(fmt.Sprintf("yrange [%g:%g]", from, to))
}

func (p *Gnuplot) Quit() {
    p.Emit("quit")
}

func (p *Gnuplot) Test() {
    p.Emit("test")
}

func RunPlot(script string) {
    err := RunBatchCmd(GNUPLOT, script)
    if err != nil {
        panic(err)
    }
}

// Run a program without stdin
func RunBatchCmd(name string, arg ...string) *exec.ExitError {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run() // Start and Wait
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            return exitError
        } else {
            panic(err)
        }
    }
    return nil
}
