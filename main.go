package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 600
	screenHeight = 600
)

type Game struct {
	backGround    color.RGBA
	fruit         Fruit
	snake         []Segment
	direction     Direction
	nextDirection Direction
	gameOver      bool
	gridSize      int
	screenSize    int
	moveTimer     int
	moveDelay     int
	score         int
}

type Segment struct {
	x, y int
}

type Fruit struct {
	x, y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func Init() *Game {
	g := &Game{
		backGround: color.RGBA{0, 0, 0, 255},
		gridSize:   20,
		screenSize: 600,
		moveDelay:  10,
	}
	g.Reset()
	return g
}

func (g *Game) Reset() {
	g.snake = []Segment{
		{x: 10, y: 10},
		{x: 9, y: 10},
		{x: 8, y: 10},
	}
	g.gameOver = false
	g.direction = Right
	g.nextDirection = Right
	g.SpawnFruit()
	g.score = 0
	g.moveTimer = 0
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.Reset()
		}
		return nil
	}

	g.Handle()

	g.moveTimer++
	if g.moveTimer < g.moveDelay {
		return nil
	}

	g.moveTimer = 0

	g.direction = g.nextDirection

	newHead := g.snake[0]
	switch g.direction {
	case Up:
		newHead.y--
	case Down:
		newHead.y++
	case Left:
		newHead.x--
	case Right:
		newHead.x++
	}

	if !g.Inside(newHead) {
		g.gameOver = true
		return nil
	}

	for _, seg := range g.snake {
		if seg.x == newHead.x && seg.y == newHead.y {
			g.gameOver = true
			return nil
		}
	}

	g.snake = append([]Segment{newHead}, g.snake...)

	if newHead.x == g.fruit.x && newHead.y == g.fruit.y {
		g.score++
		g.SpawnFruit()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

func (g *Game) Inside(head Segment) bool {
	return head.x >= 0 && head.x < g.screenSize/g.gridSize &&
		head.y >= 0 && head.y < g.screenSize/g.gridSize
}

func (g *Game) Handle() {
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.direction != Down {
		g.nextDirection = Up
	} else if ebiten.IsKeyPressed(ebiten.KeyS) && g.direction != Up {
		g.nextDirection = Down
	} else if ebiten.IsKeyPressed(ebiten.KeyD) && g.direction != Left {
		g.nextDirection = Right
	} else if ebiten.IsKeyPressed(ebiten.KeyA) && g.direction != Right {
		g.nextDirection = Left
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backGround)

	for i, seg := range g.snake {
		snakeColor := color.RGBA{0, 255, 0, 255}
		if i == 0 {
			snakeColor = color.RGBA{0, 200, 0, 255}
		}
		g.DrawGrid(screen, seg.x, seg.y, snakeColor)
	}

	g.DrawGrid(screen, g.fruit.x, g.fruit.y, color.RGBA{255, 0, 0, 255})

	ebitenutil.DebugPrint(screen, "Score: "+strconv.Itoa(g.score))

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "Game over! Press R to restart", 100, 100)
	}
}

func (g *Game) DrawGrid(screen *ebiten.Image, x, y int, color color.RGBA) {
	cellSize := float32(g.gridSize)
	cellImage := ebiten.NewImage(int(cellSize)-1, int(cellSize)-1)
	cellImage.Fill(color)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*g.gridSize), float64(y*g.gridSize))
	screen.DrawImage(cellImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) SpawnFruit() {
	notFree := make(map[[2]int]bool)
	for _, seg := range g.snake {
		notFree[[2]int{seg.x, seg.y}] = true
	}

	var free [][2]int
	for x := 0; x < g.screenSize/g.gridSize; x++ {
		for y := 0; y < g.screenSize/g.gridSize; y++ {
			if !notFree[[2]int{x, y}] {
				free = append(free, [2]int{x, y})
			}
		}
	}

	if len(free) > 0 {
		pos := free[rand.Intn(len(free))]
		g.fruit = Fruit{x: pos[0], y: pos[1]}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := Init()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
