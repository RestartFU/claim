package claim

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"log"
)

type Handler interface {
	// HandleEnter handles when a player enters a claim.
	// it also calls HandleLeave for the previous claim of the player.
	// if ctx is cancelled, it will act as if you cancelled a HandleMove context.
	HandleEnter(ctx *event.Context, p *player.Player, c *Claim)

	// HandleLeave handles when a player enters a claim.
	// if ctx is cancelled, it will act as if you cancelled a HandleMove context.
	HandleLeave(ctx *event.Context, p *player.Player, c *Claim)

	// HandleBlockBreak handles when a block is broken in a claim.
	HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, c *Claim)

	// HandleAttackEntity handles when an entity is hit while being in a claim.
	// warning: it may be called even if the source is not in the claim.
	HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, c *Claim)
}

// NopHandler ...
type NopHandler struct{}

// HandleEnter ...
func (NopHandler) HandleEnter(ctx *event.Context, p *player.Player, c *Claim) {}

// HandleLeave ...
func (NopHandler) HandleLeave(ctx *event.Context, p *player.Player, c *Claim) {}

// HandleBlockBreak ...
func (NopHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, c *Claim) {}

// HandleAttackEntity ...
func (NopHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, c *Claim) {
}

// PlayerHandler is the handler which is used to handle:
// When a block is broken in a claim.
// When a player enters or leaves a claim.
// And when a player is hurt in a claim.
type PlayerHandler struct {
	player.NopHandler
	p *player.Player
	Loader
}

// Name returns the name of the handler.
// this may be needed if you're using libraries in which a name is needed for a handler.
func (*PlayerHandler) Name() string { return "ClaimHandler" }

// NewPlayerHandler returns a new *ClaimHandler.
func NewPlayerHandler(p *player.Player, loader Loader) *PlayerHandler {
	return &PlayerHandler{
		p:      p,
		Loader: loader,
	}
}

// HandleBlockBreak handles when a block is broken,
func (c *PlayerHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	claim, err := c.LoadWithPos(pos.Vec3())
	if err != nil {
		log.Println(err)
		return
	}
	claim.h.HandleBlockBreak(ctx, pos, drops, claim)
}

// HandleAttackEntity handles when an entity is hit while being in a claim.
// warning: it may be called even if the source is not in the claim.
func (c *PlayerHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64) {
	claim, err := c.LoadWithPos(e.Position())
	if err != nil {
		log.Println(err)
		return
	}
	claim.h.HandleAttackEntity(ctx, e, force, height, claim)
}

// This makes sure that the two positions are not the same.
// This is to see if only the yaw or pitch values were changed.
func actuallyMovedXZ(old, new mgl64.Vec3) bool {
	return old.X() != new.X() || old.Z() != new.Z()
}

// HandleMove handles when a player moves.
// It calls Enter if it finds out that the player is in a claim.
func (c *PlayerHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if actuallyMovedXZ(c.p.Position(), newPos) {
		claim, err := c.LoadWithPos(newPos)
		if err != nil {
			log.Println(err)
			return
		}
		claim.Enter(ctx, c.p)
	}
}
