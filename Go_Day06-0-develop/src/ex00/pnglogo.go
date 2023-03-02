package main

import (
    "image"
    "image/color"
    "image/png"
    "os"
	"golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
)

func main() {
    img := image.NewRGBA(image.Rect(0, 0, 300, 300))
    addLabel(img, 15, 25, "Worthy Miyu")

    f, err := os.Create("amazing_logo.png")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if err := png.Encode(f, img); err != nil {
        panic(err)
    }
}

func addLabel(img *image.RGBA, x, y int, label string) {
    col := color.RGBA{200, 100, 0, 255}
    point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    d.DrawString(label)
}