package main

import (
	"database/sql"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/go-gl/mathgl/mgl64"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	c := server.DefaultConfig()
	c.Players.SaveData = false
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel
	s := server.New(&c, log)
	s.Start()
	//
	db, _ := sql.Open("sqlite3", "./test.db")
	_, _ = db.Exec("CREATE TABLE IF NOT EXISTS claims(name TEXT PRIMARY KEY, x1 INT, z1 INT, x2 INT, z2 INT);")
	testClaim := claim.NewClaim("test", area.NewVec2(mgl64.Vec2{0, 0}, mgl64.Vec2{10, 10}))
	st := &sqlite3Claimer{db: db}
	st.Store(testClaim)
	//
	for {
		p, err := s.Accept()
		if err != nil {
			return
		}
		p.Handle(claim.NewClaimHandler(p, st))
	}
}

type ClaimHandler struct {
	claim.NopHandler
	c *claim.Claim
}

func (c *ClaimHandler) HandleEnter(ctx *event.Context, p *player.Player) {
	p.Message("entered", c.c.Name())
}
func (c *ClaimHandler) HandleLeave(ctx *event.Context, p *player.Player) {
	p.Message("left", c.c.Name())
}
func (c *ClaimHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	ctx.Cancel()
}
