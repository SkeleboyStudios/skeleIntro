package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// ExitSystem exits the game when you press esc for 3 seconds.
type ExitSystem struct {
	f      *common.Font
	entity sprite
	time   float32
}

func (e *ExitSystem) New(w *ecs.World) {
	e.f = &common.Font{
		Size: 32,
		FG:   color.White,
		URL:  "title/log.ttf",
	}
	e.f.CreatePreloaded()

	e.entity = sprite{BasicEntity: ecs.NewBasic()}
	e.entity.SpaceComponent = common.SpaceComponent{
		Width:  15,
		Height: 15,
	}
	e.entity.RenderComponent = common.RenderComponent{
		Drawable: common.Text{
			Font: e.f,
			Text: "exiting",
		},
		Hidden:      true,
		StartZIndex: 25000,
		Scale:       engo.Point{X: 0.25, Y: 0.25},
	}

	w.AddEntity(&e.entity)
}

func (e *ExitSystem) Remove(basic ecs.BasicEntity) {}

func (e *ExitSystem) Update(dt float32) {
	if engo.Input.Button("Exit").Down() {
		e.entity.Hidden = false
		e.time += dt
		if e.time < 0.3 {
			e.entity.Drawable = common.Text{
				Text: "exiting .",
				Font: e.f,
			}
		} else if e.time < 0.7 {
			e.entity.Drawable = common.Text{
				Text: "exiting . .",
				Font: e.f,
			}
		} else if e.time < 1 {
			e.entity.Drawable = common.Text{
				Text: "exiting . . .",
				Font: e.f,
			}
		} else {
			engo.Exit()
		}
	} else {
		e.entity.Drawable = common.Text{
			Text: "exiting",
			Font: e.f,
		}
		e.entity.Hidden = true
		e.time = 0
	}
}
