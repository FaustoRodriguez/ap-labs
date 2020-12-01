package entities

import (
	"snakes/util"

	"github.com/hajimehoshi/ebiten"
)

type PartType int

const (
	Head = iota
	Body
	Tail
)

var imgHead, imgBody, imgTail *ebiten.Image

func init() {
	imgHead = mustLoadTexture("assets/textures/snake_head.png")
	imgBody = mustLoadTexture("assets/textures/snake.png")
}

func resolveSnakeImage(t PartType) *ebiten.Image {
	switch t {
	case Head:
		return imgHead
	case Body:
		return imgBody
	case Tail:
		return imgBody
	default:
		return nil
	}
}

type SnakeBlock struct {
	position       Point
	movingTo       Point
	movingProgress float64
	partType       PartType
	isEating       bool
}

func (s *SnakeBlock) SetPos(pos Point) {
	s.position = pos
	s.movingTo = pos
	s.movingProgress = 0
}

func (s *SnakeBlock) Move(dir Direction) {
	s.MoveTo(dir.translate(s.position))
	s.movingTo = dir.translate(s.position)
}

func (s *SnakeBlock) MoveTo(pos Point) {
	s.movingTo = pos
	s.movingProgress = 0
}

func (s *SnakeBlock) IsMoving() bool {
	return s.position != s.movingTo
}

func (s *SnakeBlock) Update() error {
	if !s.IsMoving() {
		return nil
	}
	//Calc grid step
	step := float64(util.GridSize) / float64(util.SnakeSpeed-1)
	s.movingProgress = s.movingProgress + step
	if s.movingProgress >= util.GridSize {
		s.position = s.movingTo
	}
	return nil
}

func (s *SnakeBlock) Render(screen *ebiten.Image) error {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(s.position.X, s.position.Y)

	if s.IsMoving() {
		opt.GeoM.Translate(
			float64(s.movingTo.GridX()-s.position.GridX())*s.movingProgress,
			float64(s.movingTo.GridY()-s.position.GridY())*s.movingProgress)
	}

	screen.DrawImage(resolveSnakeImage(s.partType), &opt)
	return nil
}
