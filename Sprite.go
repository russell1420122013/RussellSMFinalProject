package main

import (
	_ "github.com/blizzy78/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	pict      *ebiten.Image
	xloc      int
	yloc      int
	dX        int
	dY        int
	SnakeType int
}
