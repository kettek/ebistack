package ebistack

import (
	"bytes"
	"encoding/binary"
	"image"
	_ "image/png" // This is justified, go eat a kumquat.

	"github.com/hajimehoshi/ebiten/v2"
)

// Sheet is a collection of Stacks.
type Sheet struct {
	Stacks      []Stack
	FrameWidth  int
	FrameHeight int
	image       *ebiten.Image
}

func (s *Sheet) populateSubImages() error {
	if s.image == nil {
		panic("sheet image is nil")
	}
	for i, stack := range s.Stacks {
		for j, anim := range stack.Animations {
			for k, frame := range anim.Frames {
				for l, slice := range frame.Slices {
					img := s.image.SubImage(image.Rect(slice.x, slice.y, slice.x+s.FrameWidth, slice.y+s.FrameHeight)).(*ebiten.Image)
					s.Stacks[i].Animations[j].Frames[k].Slices[l].Image = img
				}
			}
		}
	}
	return nil
}

// AddStack adds a stack.
func (s *Sheet) AddStack(st Stack) {
	s.Stacks = append(s.Stacks, st)
}

// Stack is a collection of Animations.
type Stack struct {
	Name       string
	Animations []Animation
	sliceCount int
}

// AddAnimation adds an animation.
func (s *Stack) AddAnimation(a Animation) {
	s.Animations = append(s.Animations, a)
}

// Animation is a collection of Frames.
type Animation struct {
	Name      string
	Frametime int
	Frames    []Frame
}

// AddFrame adds a frame.
func (a *Animation) AddFrame(f Frame) {
	a.Frames = append(a.Frames, f)
}

// Frame is a collection of Slices.
type Frame struct {
	Slices []Slice
}

// Slice is a position and shading value.
type Slice struct {
	x       int
	y       int
	Shading float64
	Image   *ebiten.Image // Sub-image of the sheet if read from Staxie or whatever the image is from AddFrame.
}

// NewSheetFromStaxie creates a new Sheet from a staxie PNG.
func NewSheetFromStaxie(data []byte) (Sheet, error) {
	sheet := Sheet{}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return sheet, err
	}

	sheet.image = ebiten.NewImageFromImage(img)

	offset := 0

	readUint32 := func() uint32 {
		if offset+4 > len(data) {
			panic("out of bounds (readUint32)")
		}
		v := binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4
		return v
	}

	readUint16 := func() uint16 {
		if offset+2 > len(data) {
			panic("out of bounds (readUint16)")
		}
		v := binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2
		return v
	}

	readUint8 := func() int {
		if offset+1 > len(data) {
			panic("out of bounds (readUint8)")
		}
		v := data[offset]
		offset++
		return int(v)
	}

	readString := func() string {
		if offset+1 > len(data) {
			panic("out of bounds (readString)")
		}
		strlen := int(data[offset])
		offset++
		if offset+strlen > len(data) {
			panic("out of bounds")
		}
		str := string(data[offset : offset+strlen])
		offset += strlen
		return str
	}

	readSection := func() string {
		if offset+4 > len(data) {
			panic("out of bounds (readSection)")
		}
		str := string(data[offset : offset+4])
		offset += 4
		return str
	}

	offset += 8 // Skip PNG header

	for offset < len(data) {
		chunkSize := readUint32()
		chunkType := readSection()
		switch chunkType {
		case "stAx":
			version := readUint8()
			if version != 0 {
				panic("unsupported version")
			}
			frameWidth := readUint16()
			frameHeight := readUint16()
			stackCount := readUint16()

			sheet.FrameWidth = int(frameWidth)
			sheet.FrameHeight = int(frameHeight)

			y := 0
			for range stackCount {
				stack := Stack{}
				name := readString()
				sliceCount := readUint16()
				animationCount := readUint16()

				stack.sliceCount = int(sliceCount)
				stack.Name = name

				for range animationCount {
					animation := Animation{}
					name := readString()
					frameTime := readUint32()
					frameCount := readUint16()

					animation.Name = name
					animation.Frametime = int(frameTime)

					for range frameCount {
						frame := Frame{}

						for l := range int(sliceCount) {
							slice := Slice{
								x:       l * int(frameWidth),
								y:       y,
								Shading: float64(readUint8()) / 255.0,
							}
							frame.Slices = append(frame.Slices, slice)
						}
						if sliceCount > 0 {
							y += int(frameHeight)
						}
						animation.Frames = append(animation.Frames, frame)
					}
					stack.Animations = append(stack.Animations, animation)
				}
				sheet.Stacks = append(sheet.Stacks, stack)
			}
		default:
			offset += int(chunkSize)
		}
		readUint32() // Skip CRC
	}

	err = sheet.populateSubImages()

	return sheet, err
}
