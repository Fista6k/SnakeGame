package main

import (
	"errors"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	backGround color.RGBA

	player *Player

	isFruit bool

	fruit *Fruit
}

type Player struct {
	x, y float64

	velocity float64

	width, height float64

	playerColor color.RGBA
}

type Fruit struct {
	x, y float64

	width, height float64

	fruitColor color.RGBA
}

func Init() *Game {
	return &Game{
		backGround: color.RGBA{0, 0, 255, 255},
		player: &Player{
			x:           0,
			y:           0,
			width:       20,
			height:      20,
			playerColor: color.RGBA{255, 0, 0, 255},
			velocity:    5,
		},
		fruit: &Fruit{
			x:          320,
			y:          240,
			width:      20,
			height:     20,
			fruitColor: color.RGBA{100, 0, 0, 100},
		},
		isFruit: false,
	}
}
func (g *Game) Update() error {
	if g.Inside() {
		g.player.Update()
	} else {
		log.Fatal(errors.New("Defeat"))
	}
	if !g.isFruit {
		g.SpawnFruit()
		g.isFruit = true
	}

	if g.CheckEatingFruit() {
		g.isFruit = false
	}
	return nil
}

func (g *Game) Inside() bool {
	w := g.player.width
	h := g.player.height
	x, y := g.player.x, g.player.y
	return x+w <= screenWidth && x >= 0 && y >= 0 && y-h <= screenHeight
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.y -= p.velocity
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.y += p.velocity
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.x += p.velocity
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.x -= p.velocity
	}
}

func (g *Game) CheckEatingFruit() bool {
	return CheckCollisions(g.player.x, g.player.y, g.player.width, g.player.height,
		g.fruit.x, g.fruit.y, g.fruit.width, g.fruit.height)
}

func CheckCollisions(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backGround)
	ebitenutil.DrawRect(screen, g.player.x, g.player.y, g.player.width, g.player.height, g.player.playerColor)
	ebitenutil.DrawRect(screen, g.fruit.x, g.fruit.y, g.fruit.width, g.fruit.height, g.fruit.fruitColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) SpawnFruit() {
	g.fruit.x = float64(rand.Intn(screenWidth - 50))
	g.fruit.y = float64(rand.Intn(screenHeight - 50))
}

func main() {
	game := Init()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
