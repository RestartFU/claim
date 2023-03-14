package main

import (
	"database/sql"
	"fmt"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/go-gl/mathgl/mgl64"
)

type sqlite3Claimer struct {
	db *sql.DB
}

func (s *sqlite3Claimer) Store(claim *claim.Claim) error {
	max, min := claim.Area().Max(), claim.Area().Min()
	_, err := s.db.Exec("INSERT OR REPLACE INTO claims (name, x1, z1, x2, z2) VALUES (?, ?, ?, ?, ?);",
		claim.Name(),
		max.X(),
		max.Y(),
		min.X(),
		min.Y(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlite3Claimer) LoadWithPos(v mgl64.Vec3) (*claim.Claim, error) {
	var name string
	var x1, z1, x2, z2 int

	x := v.X()
	z := v.Z()
	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM claims WHERE %v <= x1 AND %v >= x2 AND %v <= z1 AND %v >= z2;", x, x, z, z))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&name, &x1, &z1, &x2, &z2)
	}
	if name == "" {
		wild := claim.NewClaim("The Wilderness", area.Vec2{})
		wild.Handle(&ClaimHandler{c: wild})
		return wild, nil
	}
	area := area.NewVec2(mgl64.Vec2{float64(x1), float64(z1)}, mgl64.Vec2{float64(x2), float64(z2)})
	c := claim.NewClaim(name, area)
	c.Handle(&ClaimHandler{c: c})
	return c, nil
}
