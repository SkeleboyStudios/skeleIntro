package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var YesSelectedMessageType = "Yes Selected Message"

type YesSelectedMessage struct {
	Selected bool
}

func (*YesSelectedMessage) Type() string { return YesSelectedMessageType }

var AcceptSystemPauseMessageType = "Accept System Pause Message"

type AcceptSystemPauseMessage struct {
	Pause bool
}

func (AcceptSystemPauseMessage) Type() string { return AcceptSystemPauseMessageType }

type AcceptSystem struct {
	Fnt           *common.Font
	BackgroundURL string

	cursor *CursorSystem

	bg      sprite
	yes, no selection
}

func (s *AcceptSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.cursor = sys
		}
	}

	s.bg = sprite{BasicEntity: ecs.NewBasic()}
	s.bg.Drawable, _ = common.LoadedSprite(s.BackgroundURL)
	s.bg.SetShader(common.HUDShader)
	s.bg.SetZIndex(10003)
	s.bg.Position = engo.Point{X: 200, Y: 70}
	s.bg.Scale = engo.Point{X: 0.4, Y: 0.5}
	w.AddEntity(&s.bg)

	s.yes = selection{BasicEntity: ecs.NewBasic()}
	s.yes.Drawable = common.Text{
		Text: "yes",
		Font: s.Fnt,
	}
	s.yes.SetShader(common.HUDShader)
	s.yes.SetZIndex(10005)
	s.yes.Position = engo.Point{X: 225, Y: 80}
	s.yes.Scale = engo.Point{X: 0.35, Y: 0.35}
	w.AddEntity(&s.yes)

	s.no = selection{BasicEntity: ecs.NewBasic()}
	s.no.Drawable = common.Text{
		Text: "no",
		Font: s.Fnt,
	}
	s.no.SetShader(common.HUDShader)
	s.no.SetZIndex(10005)
	s.no.Position = engo.Point{X: 360, Y: 80}
	s.no.Scale = engo.Point{X: 0.35, Y: 0.35}
	w.AddEntity(&s.no)

	engo.Mailbox.Listen(AcceptSystemPauseMessageType, func(message engo.Message) {
		msg, ok := message.(AcceptSystemPauseMessage)
		if !ok {
			return
		}
		if msg.Pause {
			s.pause()
		} else {
			s.unpause()
		}
	})
	engo.Mailbox.Listen(YesSelectedMessageType, func(message engo.Message) {
		msg, ok := message.(*YesSelectedMessage)
		if !ok {
			return
		}
		msg.Selected = s.yes.Selected
	})
}

func (s *AcceptSystem) Remove(basic ecs.BasicEntity) {}

func (s *AcceptSystem) Update(dt float32) {}

func (s *AcceptSystem) pause() {
	s.bg.Hidden = true
	s.yes.Hidden = true
	s.cursor.Remove(s.yes.BasicEntity)
	s.no.Hidden = true
	s.cursor.Remove(s.no.BasicEntity)
}

func (s *AcceptSystem) unpause() {
	s.bg.Hidden = false
	s.yes.Hidden = false
	s.yes.Selected = true
	s.cursor.AddByInterface(&s.yes)
	s.no.Hidden = false
	s.cursor.AddByInterface(&s.no)
}
