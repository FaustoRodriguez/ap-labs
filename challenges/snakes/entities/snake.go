package entities

import (
	"fmt"
	"snakes/util"

	"github.com/hajimehoshi/ebiten"
)

type Snake struct {
	game        *Game
	parts       []*SnakeBlock
	lastDir     Direction
	updateCount uint
}

func NewSnake(g *Game) *Snake {
	center := NewGridPoint(Point{
		X: util.GridWidth / 2,
		Y: util.GridHeight / 2,
	})
	s := Snake{
		game: g,
		parts: []*SnakeBlock{
			{
				position: Point{
					X: center.X - util.GridSize,
					Y: center.Y,
				},
				movingTo: Point{
					X: center.X - util.GridSize,
					Y: center.Y,
				},
				partType: Tail,
			},
			{
				position: center,
				movingTo: center,
				partType: Body,
			},
			{
				position: Point{
					X: center.X + util.GridSize,
					Y: center.Y,
				},
				movingTo: Point{
					X: center.X + util.GridSize,
					Y: center.Y,
				},
				partType: Head,
			},
		}}

	return &s
}

func (s *Snake) Update() error {
	//Key press changes direction
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		if s.lastDir != Down {
			s.lastDir = Up
		}
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		if s.lastDir != Up {
			s.lastDir = Down
		}
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		if s.lastDir != Right {
			s.lastDir = Left
		}
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		if s.lastDir != Left {
			s.lastDir = Right
		}
	}

	if s.lastDir != 0 {
		if s.updateCount == util.SnakeSpeed {
			s.updateCount = 0

			// create new head
			newHead := &SnakeBlock{
				partType: Head,
			}

			newHead.SetPos(s.parts[len(s.parts)-1].position)
			newHead.Move(s.lastDir)

			//Check if head collides with something
			//With window borders
			if newHead.movingTo.GridX() > util.GridWidth-1 || newHead.movingTo.GridX() < 0 ||
				newHead.movingTo.GridY() > util.GridHeight-1 || newHead.movingTo.GridY() < 0 {
				s.game.End()
			} else {
				//With itself
				for _, part := range s.parts {
					if newHead.movingTo == part.position {
						s.game.End()
						fmt.Printf("self")
						break
					}
				}
				//With Enemies
				for _, enemy := range s.game.enemies {
					if !enemy.dead {
						for _, part := range enemy.parts {
							if newHead.movingTo == part.position {
								s.game.End()
								fmt.Printf("Enemy")
								break
							}
						}
					}
				}
			}

			//Check if it's touching the food
			for i := 0; i < len(s.game.food); i++ {
				if s.game.food[i].position == newHead.movingTo {
					s.game.incScore()
					//Eating
					newHead.isEating = true
					s.game.food[i].eaten = true
					if s.game.AllFoodIsEaten() {
						s.game.End()
					}
				}
			}
			//Snake stops moving when game is finished
			if !s.game.IsRunning() {
				return nil
			}

			s.parts = append(s.parts, newHead)
			s.parts[len(s.parts)-2].partType = Body
			if !s.parts[0].isEating {
				if !s.parts[0].IsMoving() && s.parts[0].position == s.parts[1].position {
					s.parts = append(s.parts[:0], s.parts[0+1:]...)
				}
				s.parts[0].MoveTo(s.parts[1].position)

			} else {
				s.parts[0].isEating = false
			}

			s.parts[0].partType = Tail
		}

		s.updateCount++
	}
	for _, p := range s.parts {
		err := p.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Snake) Render(screen *ebiten.Image) error {
	for _, o := range s.parts {
		err := o.Render(screen)
		if err != nil {
			return err
		}
	}
	s.parts[0].Render(screen)
	return nil
}
