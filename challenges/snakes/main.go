package main

import (
	"log"
	"math/rand"
	"os"
	"snakes/entities"
	"snakes/util"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var game entities.Game
var enemies int
var food int
var err error

func init() {
	enemies, err = strconv.Atoi(os.Args[1])
	food, err = strconv.Atoi(os.Args[2])
	if err != nil {

	}
	game = entities.NewGame(enemies, food)
}

func update(screen *ebiten.Image) error {
	if err := game.Update(); err != nil {
		return err
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	return game.Render(screen)
}

func main() {
	if len(os.Args) == 3 {
		if err != nil {
			log.Fatal("Wrong input inserted")
		} else {
			rand.Seed(time.Now().UnixNano())
			if err = ebiten.Run(update, util.Width, util.Height, 1, "Snakes"); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("Usage: snakes <Enemies> <Food dots>")
	}
}
