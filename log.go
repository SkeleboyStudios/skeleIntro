package main

import (
	"image/color"
	"sync"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type CombatLogMessage struct {
	Msg  string
	Fnt  *common.Font
	Clip *common.Player
}

var CombatLogMessageType = "CombatLogMessage"

func (m CombatLogMessage) Type() string {
	return CombatLogMessageType
}

type CombatLogDoneMessage struct {
	Done bool
}

var CombatLogDoneMessageType = "Combat Log Done Message"

func (m *CombatLogDoneMessage) Type() string {
	return CombatLogDoneMessageType
}

type CombatLogPauseMessage struct {
	Pause bool
}

var CombatLogPauseMessageType = "Combat Log Pause Message"

func (m CombatLogPauseMessage) Type() string {
	return CombatLogPauseMessageType
}

type CombatLogClearMessage struct{}

var CombatLogClearMessageType = "Combat Log Clear Message"

func (m CombatLogClearMessage) Type() string { return CombatLogClearMessageType }

type CombatLogSystem struct {
	lock                                      sync.RWMutex
	log                                       []CombatLogMessage
	idx, charAt                               int
	done, moved                               bool
	elapsed                                   float32
	BackgroundURL, FontURL, DotURL            string
	font                                      *common.Font
	LineDelay, LetterDelay                    float32
	bg, dot1, dot2, dot3, line1, line2, line3 sprite
	paused, skipNextFrame                     bool
	dot1Shown, dot2Shown, dot3Shown           bool
}

func (s *CombatLogSystem) New(w *ecs.World) {
	//bg
	s.bg = sprite{BasicEntity: ecs.NewBasic()}
	s.bg.Drawable, _ = common.LoadedSprite(s.BackgroundURL)
	s.bg.SetShader(common.HUDShader)
	s.bg.SetZIndex(10002)
	s.bg.Width = s.bg.Drawable.Width()
	s.bg.Height = s.bg.Drawable.Height()
	s.bg.SetCenter(engo.Point{X: 320, Y: s.bg.Height / 2})
	w.AddEntity(&s.bg)

	dotTex, _ := common.LoadedSprite(s.DotURL)
	//dot1
	s.dot1 = sprite{BasicEntity: ecs.NewBasic()}
	s.dot1.Drawable = dotTex
	s.dot1.SetShader(common.HUDShader)
	s.dot1.SetZIndex(10003)
	s.dot1.SetCenter(engo.Point{X: 84, Y: 55})
	s.dot1.Hidden = true
	w.AddEntity(&s.dot1)
	//dot2
	s.dot2 = sprite{BasicEntity: ecs.NewBasic()}
	s.dot2.Drawable = dotTex
	s.dot2.SetShader(common.HUDShader)
	s.dot2.SetZIndex(10003)
	s.dot2.SetCenter(engo.Point{X: 84, Y: 35})
	s.dot2.Hidden = true
	w.AddEntity(&s.dot2)
	//dot3
	s.dot3 = sprite{BasicEntity: ecs.NewBasic()}
	s.dot3.Drawable = dotTex
	s.dot3.SetShader(common.HUDShader)
	s.dot3.SetZIndex(10003)
	s.dot3.SetCenter(engo.Point{X: 84, Y: 15})
	s.dot3.Hidden = true
	w.AddEntity(&s.dot3)

	s.font = &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  s.FontURL,
	}
	s.font.CreatePreloaded()
	//line1
	s.line1 = sprite{BasicEntity: ecs.NewBasic()}
	s.line1.Drawable = common.Text{
		Font: s.font,
		Text: "",
	}
	s.line1.SetShader(common.TextHUDShader)
	s.line1.Scale = engo.Point{X: 0.35, Y: 0.35}
	s.line1.SetZIndex(10003)
	s.line1.Position = engo.Point{X: 99, Y: 48}
	w.AddEntity(&s.line1)
	//line2
	s.line2 = sprite{BasicEntity: ecs.NewBasic()}
	s.line2.Drawable = common.Text{
		Font: s.font,
		Text: "",
	}
	s.line2.SetShader(common.TextHUDShader)
	s.line2.Scale = engo.Point{X: 0.35, Y: 0.35}
	s.line2.SetZIndex(10003)
	s.line2.Position = engo.Point{X: 99, Y: 28}
	w.AddEntity(&s.line2)
	//line3
	s.line3 = sprite{BasicEntity: ecs.NewBasic()}
	s.line3.Drawable = common.Text{
		Font: s.font,
		Text: "",
	}
	s.line3.SetShader(common.TextHUDShader)
	s.line3.Scale = engo.Point{X: 0.35, Y: 0.35}
	s.line3.SetZIndex(10003)
	s.line3.Position = engo.Point{X: 99, Y: 8}
	w.AddEntity(&s.line3)

	engo.Mailbox.Listen(CombatLogMessageType, func(message engo.Message) {
		msg, ok := message.(CombatLogMessage)
		if !ok {
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.log = append(s.log, msg)
	})

	engo.Mailbox.Listen(CombatLogDoneMessageType, func(message engo.Message) {
		msg, ok := message.(*CombatLogDoneMessage)
		if !ok {
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		msg.Done = s.done && s.idx >= len(s.log)-1
	})

	engo.Mailbox.Listen(CombatLogPauseMessageType, func(message engo.Message) {
		msg, ok := message.(CombatLogPauseMessage)
		if !ok {
			return
		}
		if msg.Pause {
			s.pause()
		} else {
			s.unpause()
		}
	})

	engo.Mailbox.Listen(CombatLogClearMessageType, func(message engo.Message) {
		_, ok := message.(CombatLogClearMessage)
		if !ok {
			return
		}
		s.clear()
	})
}

func (s *CombatLogSystem) Remove(basic ecs.BasicEntity) {}

func (s *CombatLogSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}
	s.elapsed += dt
	if s.done {
		if s.idx < len(s.log)-1 {
			s.idx++
			s.moved = false
			s.done = false
		}
	} else {
		if !s.moved && len(s.log) > 0 {
			if s.elapsed < s.LineDelay {
				return
			}
			s.elapsed = 0
			s.dot1.Hidden = false
			s.dot1Shown = true
			s.line3.Drawable = s.line2.Drawable
			s.line2.Drawable = s.line1.Drawable
			txt := s.line1.Drawable.(common.Text)
			txt.Font = s.log[s.idx].Fnt
			txt.Text = ""
			s.line1.Drawable = txt
			s.moved = true
			txt2 := s.line2.Drawable.(common.Text)
			txt3 := s.line3.Drawable.(common.Text)
			if txt2.Text != "" {
				s.dot2.Hidden = false
				s.dot2Shown = true
			}
			if txt3.Text != "" {
				s.dot3.Hidden = false
				s.dot3Shown = true
			}
		}
		if len(s.log) > 0 && s.elapsed > s.LetterDelay {
			s.charAt++
			txt := s.line1.Drawable.(common.Text)
			txt.Text = s.log[s.idx].Msg[:s.charAt]
			s.line1.Drawable = txt
			s.elapsed = 0
			if !s.log[s.idx].Clip.IsPlaying() {
				s.log[s.idx].Clip.Rewind()
				s.log[s.idx].Clip.Play()
			}
		}
		if engo.Input.Button("Y").JustPressed() || (engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonRight) {
			txt := s.line1.Drawable.(common.Text)
			txt.Text = s.log[s.idx].Msg
			s.line1.Drawable = txt
			s.charAt = 0
			s.elapsed = 0
			s.done = true
		}
		if len(s.log) > 0 && s.charAt >= len(s.log[s.idx].Msg) {
			s.charAt = 0
			s.elapsed = 0
			s.done = true
		}
	}
}

func (s *CombatLogSystem) clear() {
	s.line1.Drawable = common.Text{
		Font: s.font,
		Text: "",
	}
	s.dot1Shown = false
	s.line2.Drawable = common.Text{
		Font: s.font,
		Text: "",
	}
	s.dot2Shown = false
	s.line3.Drawable = common.Text{
		Font: s.font,
		Text: "",
	}
	s.dot3Shown = false
	s.idx = -1
	s.charAt = 0
	s.log = make([]CombatLogMessage, 0)
}

func (s *CombatLogSystem) pause() {
	s.bg.Hidden = true
	s.line1.Hidden = true
	s.line2.Hidden = true
	s.line3.Hidden = true
	s.dot1.Hidden = true
	s.dot2.Hidden = true
	s.dot3.Hidden = true
	s.paused = true
}

func (s *CombatLogSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
	s.bg.Hidden = false
	s.line1.Hidden = false
	s.line2.Hidden = false
	s.line3.Hidden = false
	if s.dot1Shown {
		s.dot1.Hidden = false
	}
	if s.dot2Shown {
		s.dot2.Hidden = false
	}
	if s.dot3Shown {
		s.dot3.Hidden = false
	}
}
