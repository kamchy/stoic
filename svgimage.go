package stoic

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	svg "github.com/ajstarks/svgo"
)

func genColorStr(attr string, h int, satperc int, lightperc int, alpha float32) string {
	return fmt.Sprintf("%s:hsla(%d, %d%%, %d%%, %f", attr, h, satperc, lightperc, alpha)
}

// GenerateSvgImage writes svg image to stdout
func GenerateSvgImage(w io.Writer) {
	s := svg.New(w)
	rand.Seed(time.Now().UnixMicro())

	wi, hi := 500, 100
	s.Startview(wi, hi, 0, 0, wi, hi)
	// draw_fala(s, wi, hi)
	draw_circles(s, wi, hi)
	s.End()
}

func draw_circles(s *svg.SVG, wi int, hi int) {
	miny := -hi / 2
	base := rand.Intn(360)
	for y := miny; y < hi; y += rand.Intn(hi / 3) {
		for x := miny; x < wi; x += rand.Intn(wi / 3) {
			col := genColorStr("fill", (base+rand.Intn(200))%360, 50, 70, 0.8)
			r := 10.0 + rand.Intn(30)
			s.Circle(x, y, r, col)
			// s.Use(x, y, "#cir", col)
		}
	}

}

func draw_fala(s *svg.SVG, wi int, hi int) {

	s.Def()
	s.Gid("fala")

	s.Path(fmt.Sprintf("M %d,%d Q %d %d %d %d Q %d %d %d %d L %d,%d %d,%d Z", 0, hi/3, wi/3, -hi/4, wi/2, hi/3, 2*wi/3, hi/2, wi, 0, wi, hi*3, 0, hi*3))
	s.Gend()
	s.DefEnd()

	miny := -hi / 2
	base := rand.Intn(360)
	for y := miny; y < hi; y += rand.Intn(hi / 3) {
		col := genColorStr("fill", (base+rand.Intn(200))%360, 50, 70, 0.8)
		s.Use(0, y, "#fala", col)
	}

}
