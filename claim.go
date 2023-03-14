package claim

import (
	"sync"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dragonfly-on-steroids/area"
)

// this map is used to keep track of the claim that a player may be in.
// this is needed because we need to know in which claim the player was before trying to enter a new one.
var players sync.Map

// NewClaim returns a new *Claim
func NewClaim(name string, area area.Vec2) *Claim {
	return &Claim{
		area: area,
		h:    NopHandler{},
		name: name,
	}
}

// Claim is a struct in which you can find all the information you need to define a Claim.
type Claim struct {
	name   string
	world  *world.World
	area   area.Vec2
	hMutex sync.RWMutex
	h      Handler
}

// Compare compares two claims, by checking if the two claim names are the same or not.
func (c *Claim) Compare(claim2 interface{}) bool {
	if claim, ok := claim2.(*Claim); ok && c.name == claim.name {
		return true
	}
	return false
}

// Name returns the name of the claim.
func (c *Claim) Name() string { return c.name }

// Area returns the area of the claim.
func (c *Claim) Area() area.Vec2 { return c.area }

// handler returns the claim handler.
func (c *Claim) handler() Handler { return c.h }

// Handle handles the claim using the given handler.
func (c *Claim) Handle(h Handler) {
	c.hMutex.Lock()
	defer c.hMutex.Unlock()
	if h == nil {
		h = NopHandler{}
	}
	c.h = h
}

// Enter enters the claim if not already in it.
// It also leaves the previous claim, if any.
func (c *Claim) Enter(ctx *event.Context, p *player.Player) {
	if claim, _ := players.Load(p); !c.Compare(claim) {
		c.h.HandleEnter(ctx, p, c)
		ctx.Continue(func() {
			if claim, ok := claim.(*Claim); ok {
				claim.Leave(ctx, p)
			}
			players.Store(p, c)
		})
	}
}

// Leave leaves the claim if the given player is in it.
func (c *Claim) Leave(ctx *event.Context, p *player.Player) {
	if claim, _ := players.Load(p); c.Compare(claim) {
		c.h.HandleLeave(ctx, p, c)
		ctx.Continue(func() {
			players.Delete(p)
		})
	}
}
