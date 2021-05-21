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

type screenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

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
	ScreenSize screenSize    `json:"screenSize"`
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

func process(f string, c config) {
	imagick.Initialize()
	defer imagick.Terminate()

	rightNow := time.Now()
	scrnSize := fmt.Sprintf("%dx%d", c.ScreenSize.Width, c.ScreenSize.Height)
	_, err := imagick.ConvertImageCommand([]string{"convert", f, "-resize", scrnSize, "/tmp/tempout.png"})
	if err != nil {
		panic(err)
	}

	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()
	mw.ReadImage("/tmp/tempout.png")

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp, and sharp DOES NOT WORK FOR US!
	err = mw.ResizeImage(uint(c.ScreenSize.Width), uint(c.ScreenSize.Height), imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}
	pw.SetColor(c.DateStamp.BackgroundColor)
	dw.SetFillColor(pw)
	mw.SetBackgroundColor(pw)

	// Display the Date and Time if DateStamp.Display is true.
	if c.DateStamp.Display {
		// Current coordinates of text
		var dx, dy float64
		// Width of a space in current font/size
		var sx float64
		dx = 1
		dw.SetFontSize(72)
		dw.SetFont(c.DateStamp.Font)
		fm := mw.QueryFontMetrics(dw, "M")
		dy = fm.CharacterHeight + fm.Descender + 1000
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
	var untilHellFreezesOver bool = true // How long does this run for?
	cmd := string("/usr/bin/gsettings")
	parm1 := []string{"set", "org.gnome.desktop.background", "picture-uri", "file:///tmp/out.png"}
	parm2 := []string{"set", "org.gnome.desktop.background", "picture-options", "centered"}
	Exec_command(cmd, parm1)
	Exec_command(cmd, parm2)

	config := readConf()
	fmt.Printf("dimensions: %dX%d Delay: %d Display Date is %v\n",
		config.ScreenSize.Width,
		config.ScreenSize.Height,
		config.Delay,
		config.DateStamp.Display)

	fmt.Println("Getting image files from the following path(s):")

	for _, xpath := range config.FsPath {
		fmt.Printf("%s\n", xpath)
	}

	// Loading the images array, slice of images
	if len(config.FsPath) > 1 {
		for _, wildcard := range config.FsPath {
			flist := List_dir(wildcard)
			config.images = append(config.images, flist...)
		}
	} else {
		config.images = List_dir(config.FsPath[0])
	}

	// Continually loop thru the images array 4-EVER!
	for untilHellFreezesOver {
		for _, image := range config.images {
			process(image, config)
			time.Sleep(config.Delay * time.Second)
		}
	}
}
