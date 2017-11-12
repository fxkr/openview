package image

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

func Initialize() {
	imagick.Initialize()
}

func Terminate() {
	imagick.Terminate()
}
