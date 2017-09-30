package safe

import (
	"encoding/json"
	"path/filepath"

	"github.com/pkg/errors"
)

// Path represents a posix file system path.
//
// Note: Paths may be empty!
type Path struct {

	// Internal representation.
	//
	// Note: Unlike the result of String(), this may be empty!
	raw string
}

// RelativePath represents a Path statically known not to be absolute.
//
// Note: RelativePaths may be empty!
type RelativePath struct {
	Path
}

// For json.Marshaler
func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// For json.Unmarshaler
func (p *RelativePath) UnmarshalJSON(data []byte) error {
	var rawPath string
	err := json.Unmarshal(data, &rawPath)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isSafeRelativePath(rawPath) {
		return errors.Errorf("Unsafe path: %v\n", rawPath)
	}
	p.raw = rawPath
	return nil
}

// UnsafeNewPath converts a string to a path.
//
// safe must be trustworthy.
// safe may be absolute.
func UnsafeNewPath(safe string) Path {
	return Path{safe}
}

// UnsafeNewPath converts a string to a path.
//
// safeRelativePath must be trustworthy.
// safeRelativePath must be relative.
// safeRelativePath must not be empty.
func UnsafeNewRelativePath(safeRelativePath string) RelativePath {
	return RelativePath{Path{safeRelativePath}}
}

// SafeNewRelativePath converts a string to a path.
//
// unsafe may be user/attacker-controlled.
// unsafe should not be absolute, otherwise an error is returned.
// unsafe should be normalized, otherwise an error is returned.
// unsafe should not be empty, otherwise an error is returned.
func SafeNewRelativePath(unsafe string) (RelativePath, error) {
	if !isSafeRelativePath(unsafe) {
		return RelativePath{}, errors.Errorf("Unsafe path: %v", unsafe)
	}
	return RelativePath{Path{unsafe}}, nil
}

// String returns the Path as a string.
//
// An empty path will be returned as ".".
func (p Path) String() string {
	if p.raw == "" {
		return "."
	}
	return p.raw
}

func (p Path) IsEmpty() bool {
	return p.raw == ""
}

// Base returns the last component of the path as a string.
func (p Path) Base() string {
	s := p.raw
	if s == "" {
		return ""
	}
	return filepath.Base(s)
}

// Join concatenates two paths.
func (p Path) Join(extensionPath RelativePath) Path {
	if p.raw == "" {
		return Path{extensionPath.raw}
	}
	return Path{filepath.Join(p.raw, extensionPath.raw)}
}

// Join concatenates two paths. trustedString must be a safe, relative path.
func (p Path) JoinUnsafe(trustedString string) Path {
	if p.raw == "" {
		return UnsafeNewPath(trustedString)
	}
	return Path{filepath.Join(p.raw, trustedString)}
}

// Join concatenates two relative paths.
func (p RelativePath) Join(extensionPath RelativePath) RelativePath {
	return RelativePath{Path{filepath.Join(p.raw, extensionPath.raw)}}
}

// isSafeRelativePath returns true if path is a possibly empty, normalized, relative path without null bytes.
func isSafeRelativePath(path string) bool {
	type State int

	const (
		state_begin  State = iota // Initial state at start of string. Expect anything except slash.
		state_slash               // Previous character was a slash between two components.
		state_dot                 // Previous character was a dot at the beginning of a component
		state_dotdot              // Previous character was a second dot at the beginning of a component
		state_safe                // We're safely within a component (we know it's not empty, "." or "..").
	)

	state := state_begin

	for _, char := range path {
		switch char {

		case '\x00':
			return false // Null bytes are forbidden.

		case '/':
			switch state {
			case state_begin:
				return false // Path starts with slash -> not relative.
			case state_slash:
				return false // Path contains two consecutive slashes -> not normalized.
			case state_dot:
				return false // Paths contains "." component -> not normalized.
			case state_dotdot:
				return false // Paths contains ".." component -> not normalized or goes above the base.
			case state_safe:
				state = state_slash // This slash is safely preceded by a directory name.
			}

		case '.':
			switch state {
			case state_begin, state_slash:
				state = state_dot // This dot is at the start of a new component.
			case state_dot:
				state = state_dotdot // This and the previous dot are at the start of a new component.
			case state_dotdot, state_safe:
				state = state_safe // Three or more leading dots in a component are acceptable.
			}

		default:
			state = state_safe // Any characters other than null, slash or dot are considered safe.
		}
	}

	// Note: empty paths are allowed
	return state == state_begin || state == state_safe
}
