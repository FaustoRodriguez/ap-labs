package entities

import (
	"math/rand"
	"snakes/util"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type Enemy struct {
	game        *Game
	parts       []*SnakeBlock //Same Blocks as the Main Snake
	lastDir     Direction
	updateCount uint
	dead        bool
}

//Builds a new enemy
func NewEnemy(g *Game) *Enemy {
	//Random origin location is set
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(util.GridWidth)
	y := rand.Intn(util.GridHeight)
	center := NewGridPoint(Point{
		X: float64(x),
		Y: float64(y),
	})
	//Builds enemy
	e := Enemy{
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
	// Random directions set on a thread
	ticker := time.NewTicker(1200 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				switch rand.Intn(4) {
				case 0:
					if e.lastDir != Down {
						e.lastDir = Up
					}
				case 1:
					if e.lastDir != Up {
						e.lastDir = Down
					}
				case 2:
					if e.lastDir != Right {
						e.lastDir = Left
					}
				case 3:
					if e.lastDir != Left {
						e.lastDir = Right
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return &e
}

//Update function for Ebiten to work with it
func (e *Enemy) Update() error {
	if e.dead {
		return nil
	}
	if e.lastDir != 0 {
		if e.updateCount == util.SnakeSpeed {
			e.updateCount = 0
			newHead := &SnakeBlock{
				partType: Head,
			}
			newHead.SetPos(e.parts[len(e.parts)-1].position)
			newHead.Move(e.lastDir)

			//Check if head hits with something
			if newHead.movingTo.GridX() > util.GridWidth-1 || newHead.movingTo.GridX() < 0 || newHead.movingTo.GridY() > util.GridHeight-1 || newHead.movingTo.GridY() < 0 {
				e.dead = true
				//With the borders of the window
			} else {
				//With itself
				for _, part := range e.parts {
					if newHead.movingTo == part.position {
						e.dead = true
					}
				}
				//With other Enemy snakes
				for _, enemy := range e.game.enemies {
					if !enemy.dead && enemy.parts[0].position != e.parts[0].position {
						for _, part := range enemy.parts {
							if newHead.movingTo == part.position {
								e.dead = true
							}
						}
					}
				}
				// With Main Snake
				for _, part := range e.game.snake.parts {
					if newHead.movingTo == part.position {
						e.dead = true
					}
				}
			}

			//Check if it's touching the food
			for i := 0; i < len(e.game.food); i++ {
				if e.game.food[i].position == newHead.movingTo {
					//Eating
					newHead.isEating = true
					e.game.food[i].eaten = true
					if e.game.AllFoodIsEaten() {
						e.game.End()
					}
				}
			}

			//Snake stops moving when game is finished
			if !e.game.IsRunning() {
				return nil
			}
			e.parts = append(e.parts, newHead)
			e.parts[len(e.parts)-2].partType = Body
			if !e.parts[0].isEating {
				if !e.parts[0].IsMoving() && e.parts[0].position == e.parts[1].position {
					e.parts = append(e.parts[:0], e.parts[0+1:]...)
				}
				e.parts[0].MoveTo(e.parts[1].position)

			} else {
				e.parts[0].isEating = false
			}
			e.parts[0].partType = Tail
		}
		e.updateCount++
	}
	for _, p := range e.parts {
		err := p.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

//Render Function for Ebiten to work with this
func (e Enemy) Render(screen *ebiten.Image) error {
	if e.dead {
		return nil
	}
	for _, o := range e.parts {
		err := o.Render(screen)
		if err != nil {
			return err
		}
	}
	e.parts[0].Render(screen)
	return nil
}
