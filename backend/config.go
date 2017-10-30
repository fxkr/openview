package backend

import (
	"github.com/fxkr/openview/backend/util/safe"
)

type Config struct {
	ResourceDir safe.Path
	CacheDir    safe.Path
	ImageDir    safe.Path

	ListenAddress string
}
