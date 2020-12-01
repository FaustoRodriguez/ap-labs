package entities

import (
	"math/rand"
	"snakes/util"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var imgFood *ebiten.Image

//Sets food Image
func init() {
	imgFood = mustLoadTexture("assets/textures/food.png")
	rand.Seed(time.Now().UnixNano())
}

type Food struct {
	position Point
	eaten    bool
}

func NewFoodItem(g *Game) *Food {
	f := Food{}
	f.eaten = false
	f.setPosition()
	return &f
}

func (f *Food) Update() error {
	return nil
}

func (f *Food) Render(screen *ebiten.Image) error {
	if f.eaten {
		return nil
	}
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(f.position.X, f.position.Y)
	screen.DrawImage(imgFood, &opt)
	return nil
}

func (c *Food) setPosition() {
	maxX := util.GridWidth
	maxY := util.GridHeight
	x := rand.Intn(maxX)
	y := rand.Intn(maxY)
	c.position = NewGridPoint(Point{
		X: float64(x),
		Y: float64(y),
	})
}
