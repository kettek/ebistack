package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebistack"
)

type Game struct {
	spriteStacks []*ebistack.Sprite
	rotate       float64
}

func (g *Game) Update() error {
	g.rotate += 0.01
	for _, sprite := range g.spriteStacks {
		sprite.Rotation = g.rotate
		sprite.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(200, 200)
	for _, sprite := range g.spriteStacks {
		sprite.Draw(screen, opts)
	}
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ow, oh
}

func main() {
	g := &Game{}

	f, err := os.ReadFile("autoCannon.png")
	if err != nil {
		panic(err)
	}
	sheet, err := ebistack.NewSheetFromStaxie(f)
	if err != nil {
		panic(err)
	}
	sprite := ebistack.MakeSprite(&sheet, "top", "attack")
	sprite.Scale = 4

	g.spriteStacks = append(g.spriteStacks, &sprite)

	ebiten.SetWindowSize(400, 400)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
