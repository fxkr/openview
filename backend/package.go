package backend

import (
	"github.com/fxkr/openview/backend/image"
)

func Initialize() {
	image.Initialize()
}

func Terminate() {
	image.Terminate()
}
