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

type dateStr struct {
	Display         bool    `json:"display"`
	Font            string  `json:"font"`
	FontSize        float64 `json:"fontSize"`
	Color           string  `json:"color"`
	BackgroundColor string  `json:"backgroundColor"`
	Position        string  `json:"position"`
	Format          string  `json:"format"`
}

type config struct {
	images     []string
	ScreenSize string        `json:"screenSize"`
	FsPath     []string      `json:"fsPath"`
	Delay      time.Duration `json:"delay"`
	DateStamp  dateStr       `json:"dateStamp"`
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

func desktop(f string, c config) {
	imagick.Initialize()
	defer imagick.Terminate()
	rightNow := time.Now()
	_, err := imagick.ConvertImageCommand([]string{"convert", f, "-resize", c.ScreenSize, "/tmp/tempout.png"})
	if err != nil {
		panic(err)
	}

	// Current coordinates of text
	var dx, dy float64
	// Width of a space in current font/size
	var sx float64

	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	mw.ReadImage("/tmp/tempout.png")
	// UnComment to Get original image size
	//width := mw.GetImageWidth()
	//height := mw.GetImageHeight()

	// Calculate size
	hWidth := uint(1920)
	hHeight := uint(1080)

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}
	mw.SetSize(5760, 1080)
	pw.SetColor(c.DateStamp.BackgroundColor)
	dw.SetFillColor(pw)

	mw.SetBackgroundColor(pw)
	dx = 1
	dw.SetFontSize(72)
	dw.SetFont(c.DateStamp.Font)
	fm := mw.QueryFontMetrics(dw, "M")
	dy = fm.CharacterHeight + fm.Descender + 1000

	// text
	if c.DateStamp.Display {
		dw.SetTextUnderColor(pw)
		draw_setfont(mw, dw, c.DateStamp.Font, c.DateStamp.FontSize, c.DateStamp.Color, &sx)
		draw_metrics(mw, dw, &dx, dy, sx,
			fmt.Sprintf(c.DateStamp.Position, rightNow.Format(c.DateStamp.Format)))
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

func main() {
	var untilHellFreezesOver bool = true
	cmd := string("/usr/bin/gsettings")
	parm1 := []string{"set", "org.gnome.desktop.background", "picture-uri", "file:///tmp/out.png"}
	parm2 := []string{"set", "org.gnome.desktop.background", "picture-options", "centered"}
	c := readConf()

	fmt.Println("STARTED")
	fmt.Printf("dimensions: %s Delay: %d Display Date is %v\n", c.ScreenSize, c.Delay, c.DateStamp.Display)
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

			// convert(f, c)
			desktop(f, c)
			time.Sleep(c.Delay * time.Second)
		}
	}
}
