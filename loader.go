package claim

import (
	"github.com/go-gl/mathgl/mgl64"
)

// Loader is an interface to load claims with the given vec3.
type Loader interface {
	// LoadWithPos loads a claim with the given vec3.
	// it may be nil, but I highly suggest that you use some kind of Wilderness claim instead.
	LoadWithPos(mgl64.Vec3) (*Claim, error)
}
