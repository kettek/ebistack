package ebistack

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite provides an interface for managing and rendering a sprite stack.
type Sprite struct {
	sheet      *Sheet
	stack      *Stack
	animation  *Animation
	frame      *Frame
	frameindex int
	frametime  int
	FrameRate  int // FrameRate expressed of units per update.
	Rotation   float64
	Scale      float64
}

// MakeSprite makes a sprite stack from the given sheet using the provided stack and its animation names.
func MakeSprite(sheet *Sheet, stackName string, animName string) Sprite {
	tps := ebiten.TPS() // 60 per default.
	framerate := 1000 / tps
	s := Sprite{sheet: sheet, Scale: 1, FrameRate: framerate}
	s.SetStack(stackName)
	s.SetAnimation(animName)
	return s
}

// Update increments frames if they exist.
func (s *Sprite) Update() {
	s.frametime += s.FrameRate
	if s.frametime >= s.animation.Frametime {
		s.frametime = 0
		s.frameindex++
		if s.frameindex >= len(s.animation.Frames) {
			s.frameindex = 0
		}
		s.frame = &s.animation.Frames[s.frameindex]
	}
}

// Draw draws the sprite stack to the given position.
func (s *Sprite) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	sliceOpts := &ebiten.DrawImageOptions{}
	sliceOpts.GeoM.Translate(float64(s.sheet.FrameWidth)/-2, float64(s.sheet.FrameHeight)/-2)
	sliceOpts.GeoM.Scale(s.Scale, s.Scale)
	sliceOpts.GeoM.Rotate(s.Rotation)
	sliceOpts.GeoM.Concat(opts.GeoM)
	for _, slice := range s.frame.Slices {
		screen.DrawImage(slice.Image, sliceOpts)
		sliceOpts.GeoM.Translate(0, -1*s.Scale)
	}
}

// SetStack sets the sprite stack to the given value.
func (s *Sprite) SetStack(name string) {
	var animName string
	if s.animation != nil {
		animName = s.animation.Name
	}
	for _, stack := range s.sheet.Stacks {
		if stack.Name == name {
			s.stack = &stack
			s.SetAnimation(animName)
			return
		}
	}
}

// SetAnimation sets the animation used by the sprite stack.
func (s *Sprite) SetAnimation(name string) {
	for _, anim := range s.stack.Animations {
		if anim.Name == name {
			s.animation = &anim
			s.frametime = anim.Frametime
			s.SetFrame(s.frameindex)
			return
		}
	}
}

// SetFrame sets the frame of the current animation.
func (s *Sprite) SetFrame(index int) {
	if index >= len(s.animation.Frames) {
		index = 0
	}
	s.frameindex = index
	s.frame = &s.animation.Frames[s.frameindex]
}

// Slices returns the current frame's slices.
func (s *Sprite) Slices() []Slice {
	if s.frame == nil {
		return nil
	}
	return s.frame.Slices
}
