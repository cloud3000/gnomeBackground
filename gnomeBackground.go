package main

/*
#include "MagickWand/MagickWand.h"
*/

import "C"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type config struct {
	images     []string
	ScreenSize string        `json:"screenSize"`
	FsPath     []string      `json:"fsPath"`
	Delay      time.Duration `json:"delay"`
	DateStamp  bool          `json:"dateStamp"`
	FontSize   float64       `json:"fontSize"`
}

func draw_setfont(mw *imagick.MagickWand, dw *imagick.DrawingWand, font string, size float64, colour string, sx *float64) {
	sflag := false

	if len(font) > 0 {
		dw.SetFont(font)
		sflag = true
	}

	if len(colour) > 0 {
		pw := imagick.NewPixelWand()
		pw.SetColor(colour)
		dw.SetFillColor(pw)
		pw.Destroy()
	}

	if size > 0 {
		dw.SetFontSize(size)
		sflag = true
	}

	if sflag {
		fm := mw.QueryFontMetrics(dw, " ")
		*sx = fm.TextWidth
	}
}

func draw_metrics(mw *imagick.MagickWand, dw *imagick.DrawingWand, dx *float64, dy, sx float64, text string) {
	mw.AnnotateImage(dw, *dx, dy, 0, text)
	mw.DrawImage(dw)

	// get the font metrics
	fm := mw.QueryFontMetrics(dw, text)
	if fm != nil {
		// Adjust the new x coordinate
		*dx += fm.TextWidth + sx
	}
}

func addText(t string, c config) {
	imagick.Initialize()
	defer imagick.Terminate()

	// Current coordinates of text
	var dx, dy float64
	// Width of a space in current font/size
	var sx float64

	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()
	//mw.SetImageAlphaChannel(imagick.CHANNEL_ALPHA)
	mw.ReadImage("/tmp/tempout.png")
	mw.SetSize(5760, 1080)
	dx = 200
	dw.SetFontSize(72)
	dw.SetFont("Times-New-Roman")
	pw.SetColor("#FFFFFF")
	dw.SetFillColor(pw)
	fm := mw.QueryFontMetrics(dw, "M")
	dy = fm.CharacterHeight + fm.Descender + 1000

	// text
	if c.DateStamp {
		draw_setfont(mw, dw, "Times-New-Roman", c.FontSize, "#000000", &sx)
		draw_metrics(mw, dw, &dx, dy, sx, t)
		mw.DrawImage(dw)
	}

	// Now write the magickwand image
	mw.WriteImage("/tmp/out.png")
}
func readConf() config {
	// json data
	var obj config

	// read file
	data, err := ioutil.ReadFile("./gnomeBackground.json")
	if err != nil {
		fmt.Print(err)
		return obj
	}

	// Unmarshal json.
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
		return obj
	}

	return obj
}

func List_dir(wildcard string) []string {
	files, _ := filepath.Glob(wildcard)
	return files
}

func Exec_command(cmd string, args []string) []byte {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Println(err)
	}
	return out
}

func convert(f string, c config) {
	imagick.Initialize()
	defer imagick.Terminate()

	// mw := imagick.NewMagickWand()
	// dw := imagick.NewDrawingWand()
	// pw := imagick.NewPixelWand()

	_, err := imagick.ConvertImageCommand([]string{"convert", f, "-resize", c.ScreenSize, "/tmp/tempout.png"})
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Metadata:\n%s\n", ret.Meta)
	// mw.NewImage()
	// mw.ReadImage("/tmp/out.png")
}
func main() {
	var untilHellFreezesOver bool = true
	cmd := string("/usr/bin/gsettings")
	parm1 := []string{"set", "org.gnome.desktop.background", "picture-uri", "file:///tmp/out.png"}
	parm2 := []string{"set", "org.gnome.desktop.background", "picture-options", "centered"}
	c := readConf()

	fmt.Println("STARTED")
	fmt.Printf("dimensions: %s Delay: %d Display Date is %v\n", c.ScreenSize, c.Delay, c.DateStamp)
	fmt.Println("Getting image files from the following path(s):")
	Exec_command(cmd, parm1)
	Exec_command(cmd, parm2)

	for _, xpath := range c.FsPath {
		fmt.Printf("%s\n", xpath)
	}

	if len(c.FsPath) > 1 {
		for _, wildcard := range c.FsPath {
			flist := List_dir(wildcard)
			c.images = append(c.images, flist...)
		}
	} else {
		c.images = List_dir(c.FsPath[0])
	}
	for untilHellFreezesOver {
		for _, f := range c.images {
			rightNow := time.Now()
			convert(f, c)
			addText(fmt.Sprintf("%v", rightNow.Format(time.RFC1123)), c)
			time.Sleep(c.Delay * time.Second)
		}
	}
}
