package entities

import (
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	snake   *Snake
	hud     *Hud
	running bool
	points  int
	enemies []*Enemy
	food    []*Food
	won     bool
}

func NewGame(enemies int, food int) Game {
	g := Game{
		running: true,
		points:  0,
	}
	g.snake = NewSnake(&g)
	g.food = make([]*Food, food)
	for i := 0; i < food; i++ {
		g.food[i] = NewFoodItem(&g)
	}
	g.hud = NewHud(&g)
	g.enemies = make([]*Enemy, enemies)
	for i := 0; i < enemies; i++ {
		g.enemies[i] = NewEnemy(&g)
	}
	return g
}

// Check if this is the last food
func (g *Game) AllFoodIsEaten() bool {
	c := 0
	for i := 0; i < len(g.food); i++ {
		if g.food[i].eaten {
			c += 1
		}
	}
	return len(g.food)-c == 0
}

//Last instructions before the game ends
func (g *Game) End() {
	if g.AllFoodIsEaten() {
		g.won = true
		for i := 0; i < len(g.enemies); i++ {
			if len(g.enemies[i].parts) > len(g.snake.parts) && !g.enemies[i].dead {
				g.won = false
			}
		}
	} else {
		g.won = false
	}
	g.running = false
}

func (g Game) IsRunning() bool {
	return g.running
}

func (g *Game) incScore() {
	g.points++
}

//Sends each element to execute
func (g *Game) Update() error {
	if g.IsRunning() {
		if err := executeUpdates(g.snake.Update, g.hud.Update); err != nil {
			return err
		}
		for i := 0; i < len(g.enemies); i++ {
			executeUpdates(g.enemies[i].Update)
		}
		for i := 0; i < len(g.food); i++ {
			executeUpdates(g.food[i].Update)
		}
	}
	return nil
}

//This is where the renders execute
func (g *Game) Render(screen *ebiten.Image) error {
	if err := executeRenderers(screen, g.snake.Render, g.hud.Render); err != nil {
		return err
	}
	for i := 0; i < len(g.enemies); i++ {
		executeRenderers(screen, g.enemies[i].Render)
	}
	for i := 0; i < len(g.food); i++ {
		executeRenderers(screen, g.food[i].Render)
	}
	return nil
}

//The function just runs the updater
func executeUpdates(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

//The function just sends the screen to eache renderer
func executeRenderers(screen *ebiten.Image, fns ...func(screen *ebiten.Image) error) error {
	for _, fn := range fns {
		if err := fn(screen); err != nil {
			return err
		}
	}

	return nil
}
