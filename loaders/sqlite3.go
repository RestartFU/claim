package loaders

import (
	"database/sql"
	"fmt"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/go-gl/mathgl/mgl64"
)

type SQL struct {
	db *sql.DB
	h  claim.Handler
}

func NewSQL(db *sql.DB, h claim.Handler) (*SQL, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS claims(name TEXT PRIMARY KEY, x1 INT, z1 INT, x2 INT, z2 INT);")
	if err != nil {
		return nil, err
	}
	return &SQL{db: db, h: h}, nil
}

func (s *SQL) Delete(claim *claim.Claim) error {
	_, err := s.db.Exec(fmt.Sprintf("DELETE FROM claims WHERE name='%s';", claim.Name()))
	return err
}

func (s *SQL) Store(claim *claim.Claim) error {
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

func (s *SQL) LoadWithPos(v mgl64.Vec3) (*claim.Claim, error) {
	var (
		name           string
		x1, z1, x2, z2 int
	)
	x := v.X()
	z := v.Z()
	query := fmt.Sprintf("SELECT * FROM claims WHERE %v <= x1 AND %v >= x2 AND %v <= z1 AND %v >= z2;", x, x, z, z)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&name, &x1, &z1, &x2, &z2)
		if err != nil {
			return nil, err
		}
	}
	if name == "" {
		wild := claim.NewClaim("The Wilderness", area.Vec2{})
		wild.Handle(s.h)
		return wild, nil
	}
	xf := float64(x1)
	xl := float64(x2)
	zf := float64(z1)
	zl := float64(z2)

	max, min := mgl64.Vec2{xf, zf}, mgl64.Vec2{xl, zl}
	a := area.NewVec2(max, min)
	newClaim := claim.NewClaim(name, a)
	newClaim.Handle(s.h)
	return newClaim, nil
}
