package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"time"
)

//go:embed Go_sprites/*
var EmbeddedAssets embed.FS
var textWidget *widget.Text
var rootContainer *widget.Container
var game Game

const (
	GameWidth    = 700
	GameHeight   = 700
	PlayerSpeed  = 10
	SnakeHead    = 0
	SnakeBody    = 1
	StartingSize = 5
)

type Game struct {
	snake     []Sprite
	food      []Sprite
	shortener []Sprite
	score     int
	drawOps   ebiten.DrawImageOptions
	AppUI     *ebitenui.UI
}

func (g *Game) Update() error {
	g.AppUI.Update()
	processPlayerInput(g, 0)
	for i := 1; i < len(g.snake); i++ {
		moveSnakeBody(g, i)
	}
	displayScore(g)
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	sNum := len(g.snake)
	for ind, _ := range g.snake {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(g.snake[(sNum-1)-ind].xloc), float64(g.snake[(sNum-1)-ind].yloc))
		screen.DrawImage(g.snake[(sNum-1)-ind].pict, &g.drawOps)
	}
	for _, foodSprite := range g.food {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(foodSprite.xloc), float64(foodSprite.yloc))
		screen.DrawImage(foodSprite.pict, &g.drawOps)
	}
	for _, shortSprite := range g.shortener {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(shortSprite.xloc), float64(shortSprite.yloc))
		screen.DrawImage(shortSprite.pict, &g.drawOps)
	}
	g.AppUI.Draw(screen)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func main() {
	rand.NewSource(time.Now().UnixNano())
	println(rand.Intn(600))
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("Final project")
	game = Game{score: 0}
	for i := 0; i < StartingSize; i++ {
		if i < 1 {
			game.snake = append(game.snake, Sprite{
				pict:      loadPNGImageFromEmbedded("snakeHead.png"),
				xloc:      GameWidth / 2,
				yloc:      GameHeight / 2,
				dX:        0,
				dY:        0,
				SnakeType: SnakeHead,
			})
		} else {
			game.snake = append(game.snake, Sprite{
				pict:      loadPNGImageFromEmbedded("snakeTail.png"),
				xloc:      GameWidth / 2,
				yloc:      GameHeight / 2,
				dX:        0,
				dY:        0,
				SnakeType: SnakeBody,
			})
		}
	}
	game.AppUI = MakeUIWindow()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	rootContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
	)
	textInfo := widget.TextOptions{}.Text("Score: 0", basicfont.Face7x13, color.White)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)
	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("Go_sprites")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("Go_sprites/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func moveSnakeBody(theGame *Game, i int) {
	if theGame.snake[i].yloc < theGame.snake[0].yloc {
		theGame.snake[i].yloc += PlayerSpeed - i
	} else if theGame.snake[i].yloc > theGame.snake[0].yloc {
		theGame.snake[i].yloc -= PlayerSpeed - i
	}
	if theGame.snake[i].xloc < theGame.snake[0].xloc {
		theGame.snake[i].xloc += PlayerSpeed - i
	} else if theGame.snake[i].xloc != theGame.snake[0].xloc {
		theGame.snake[i].xloc -= PlayerSpeed - i
	}
}
func processPlayerInput(theGame *Game, snakeNum int) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		theGame.snake[snakeNum].dY = -PlayerSpeed + snakeNum
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		theGame.snake[snakeNum].dY = PlayerSpeed - snakeNum
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		theGame.snake[snakeNum].dY = 0
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		theGame.snake[snakeNum].dX = -PlayerSpeed + snakeNum
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		theGame.snake[snakeNum].dX = PlayerSpeed - snakeNum
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		theGame.snake[snakeNum].dX = 0
	}
	theGame.snake[snakeNum].yloc += theGame.snake[snakeNum].dY
	if theGame.snake[snakeNum].yloc <= 0 {
		theGame.snake[snakeNum].dY = 0
		theGame.snake[snakeNum].yloc = 0
	} else if theGame.snake[snakeNum].yloc+theGame.snake[snakeNum].pict.Bounds().Size().Y > GameHeight {
		theGame.snake[snakeNum].dY = 0
		theGame.snake[snakeNum].yloc = GameHeight - theGame.snake[snakeNum].pict.Bounds().Size().Y
	}
	theGame.snake[snakeNum].xloc += theGame.snake[snakeNum].dX
	if theGame.snake[snakeNum].xloc <= 0 {
		theGame.snake[snakeNum].dX = 0
		theGame.snake[snakeNum].xloc = 0
	} else if theGame.snake[snakeNum].xloc+theGame.snake[snakeNum].pict.Bounds().Size().X > GameWidth {
		theGame.snake[snakeNum].dX = 0
		theGame.snake[snakeNum].xloc = GameWidth - theGame.snake[snakeNum].pict.Bounds().Size().X
	}
}

func displayScore(theGame *Game) {
	message := fmt.Sprintf("Score: %d", theGame.score)
	textWidget.Label = message
}

func removeComSprite(sl *[]Sprite, s Sprite) []Sprite {
	for idx, v := range *sl {
		if v == s {
			return append((*sl)[0:idx], (*sl)[idx+1:]...)
		}
	}
	return *sl
	// code above taken from https://stackoverflow.com/questions/32798061/how-can-i-remove-a-struct-from-a-slice-of-structs-in-go
}
