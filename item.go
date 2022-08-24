package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var ItemSelectSystemPauseMessageType = "ItemSelectSystemPauseMessage"

type ItemSelectSystemPauseMessage struct {
	Pause bool
}

func (ItemSelectSystemPauseMessage) Type() string { return ItemSelectSystemPauseMessageType }

type itemSelectEntity struct {
	*Character
}

type ItemSelectSystem struct {
	entities               []itemSelectEntity
	paused, skipNextFrame  bool
	cursor                 sprite
	short1, short2         sprite
	short3, short4         sprite
	name, desc             sprite
	fnt                    *common.Font
	scale                  engo.Point
	setIdx, curIdx, topIdx int
}

func (s *ItemSelectSystem) New(w *ecs.World) {
	s.scale = engo.Point{X: 0.5, Y: 0.5}
	curTex, _ := common.LoadedSprite("title/cursor.png")
	s.cursor = sprite{BasicEntity: ecs.NewBasic()}
	s.cursor.Drawable = curTex
	s.cursor.Width = curTex.Width()
	s.cursor.Height = curTex.Height()
	s.cursor.SetZIndex(3)
	s.cursor.Hidden = true
	w.AddEntity(&s.cursor)

	s.short1 = sprite{BasicEntity: ecs.NewBasic()}
	s.short1.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short1.Position = engo.Point{X: 50, Y: 220}
	s.short1.SetZIndex(3)
	s.short1.Scale = s.scale
	s.short1.Hidden = true
	w.AddEntity(&s.short1)

	s.short2 = sprite{BasicEntity: ecs.NewBasic()}
	s.short2.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short2.Position = engo.Point{X: 50, Y: 250}
	s.short2.SetZIndex(3)
	s.short2.Scale = s.scale
	s.short2.Hidden = true
	w.AddEntity(&s.short2)

	s.short3 = sprite{BasicEntity: ecs.NewBasic()}
	s.short3.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short3.Position = engo.Point{X: 50, Y: 280}
	s.short3.SetZIndex(3)
	s.short3.Scale = s.scale
	s.short3.Hidden = true
	w.AddEntity(&s.short3)

	s.short4 = sprite{BasicEntity: ecs.NewBasic()}
	s.short4.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short4.Position = engo.Point{X: 50, Y: 310}
	s.short4.SetZIndex(3)
	s.short4.Scale = s.scale
	s.short4.Hidden = true
	w.AddEntity(&s.short4)

	s.name = sprite{BasicEntity: ecs.NewBasic()}
	s.name.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.name.Position = engo.Point{X: 205, Y: 220}
	s.name.SetZIndex(3)
	s.name.Scale = engo.Point{X: s.scale.X * 1.05, Y: s.scale.Y * 1.05}
	s.name.Hidden = true
	w.AddEntity(&s.name)

	s.desc = sprite{BasicEntity: ecs.NewBasic()}
	s.desc.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.desc.Position = engo.Point{X: 205, Y: 250}
	s.desc.SetZIndex(3)
	s.desc.Scale = engo.Point{X: s.scale.X * 0.95, Y: s.scale.Y * 0.95}
	s.desc.Hidden = true
	w.AddEntity(&s.desc)

	engo.Mailbox.Listen(ItemSelectSystemPauseMessageType, func(message engo.Message) {
		msg, ok := message.(ItemSelectSystemPauseMessage)
		if !ok {
			return
		}
		if msg.Pause {
			s.pause()
		} else {
			s.unpause()
		}
	})
}

func (s *ItemSelectSystem) Add(chara *Character) {
	s.entities = append(s.entities, itemSelectEntity{chara})
}

func (s *ItemSelectSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(Characterable)
	if ok {
		s.Add(o.GetCharacter())
	}
}

func (s *ItemSelectSystem) Remove(b ecs.BasicEntity) {
	d := -1
	for i, e := range s.entities {
		if e.ID() == b.ID() {
			d = i
			break
		}
	}
	if d >= 0 {
		s.entities = append(s.entities[:d], s.entities[d+1:]...)
	}
}

func (s *ItemSelectSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}
	var chara *Character
	for _, e := range s.entities {
		if e.IsCardSelected {
			chara = e.GetCharacter()
			break
		}
	}
	if chara == nil {
		return
	}
	if engo.Input.Button("up").JustPressed() {
		s.curIdx--
		if s.curIdx < 0 {
			s.curIdx = 0
		}
	} else if engo.Input.Button("down").JustPressed() {
		s.curIdx++
		if s.curIdx >= len(chara.Inventory) {
			s.curIdx = len(chara.Inventory) - 1
			if s.curIdx < 0 {
				s.curIdx = 0
			}
		}
	} else if engo.Input.Button("A").JustPressed() {
		if s.curIdx < len(chara.Inventory) {
			chara.SelectedItem = chara.Inventory[s.curIdx]
		} else {
			engo.Mailbox.Dispatch(PhaseSetMessage{
				Phase: CardSelectPhase,
			})
			engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			return
		}
		engo.Mailbox.Dispatch(PhaseSetMessage{
			Phase: TargetPhase,
		})
		engo.Mailbox.Dispatch(PhaseDequeuMessage{})
	} else if engo.Input.Button("B").JustPressed() {
		engo.Mailbox.Dispatch(PhaseSetMessage{
			Phase: CardSelectPhase,
		})
		engo.Mailbox.Dispatch(PhaseDequeuMessage{})
	}

	if s.curIdx != s.setIdx {
		if s.curIdx < s.topIdx {
			s.short1.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.short2.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.short3.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.short4.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.name.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.desc.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			switch {
			case len(chara.Inventory) > 4:
				s.short4.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx+3].Shorthand,
				}
				s.short4.Scale = chara.TextScale
				fallthrough
			case len(chara.Inventory) > 3:
				s.short3.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx+2].Shorthand,
				}
				s.short3.Scale = chara.TextScale
				fallthrough
			case len(chara.Inventory) > 2:
				s.short2.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx+1].Shorthand,
				}
				s.short2.Scale = chara.TextScale
				fallthrough
			case len(chara.Inventory) > 1:
				s.short1.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx].Shorthand,
				}
				s.short1.Scale = chara.TextScale
			}
			s.topIdx = s.curIdx
		} else if s.curIdx > s.topIdx+3 {
			s.short1.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.short2.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.short3.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.short4.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.name.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			s.desc.Drawable = common.Text{
				Font: s.fnt,
				Text: "---",
			}
			switch {
			case len(chara.Inventory) > 4:
				s.short4.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx].Shorthand,
				}
				s.short4.Scale = chara.TextScale
				fallthrough
			case len(chara.Inventory) > 3:
				s.short3.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx-1].Shorthand,
				}
				s.short3.Scale = chara.TextScale
				fallthrough
			case len(chara.Inventory) > 2:
				s.short2.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx-2].Shorthand,
				}
				s.short2.Scale = chara.TextScale
				fallthrough
			case len(chara.Inventory) > 1:
				s.short1.Drawable = common.Text{
					Font: chara.Font,
					Text: chara.Inventory[s.curIdx-3].Shorthand,
				}
				s.short1.Scale = chara.TextScale
			}
			s.topIdx = s.curIdx - 3
		} else {
			switch s.curIdx - s.topIdx {
			case 0:
				s.cursor.Position = engo.Point{X: s.short1.Position.X - s.cursor.Width - 2, Y: s.short1.Position.Y + 5}
			case 1:
				s.cursor.Position = engo.Point{X: s.short2.Position.X - s.cursor.Width - 2, Y: s.short2.Position.Y + 5}
			case 2:
				s.cursor.Position = engo.Point{X: s.short3.Position.X - s.cursor.Width - 2, Y: s.short3.Position.Y + 5}
			case 3:
				s.cursor.Position = engo.Point{X: s.short4.Position.X - s.cursor.Width - 2, Y: s.short4.Position.Y + 5}
			}
		}
		s.name.Drawable = common.Text{
			Font: chara.Font,
			Text: chara.Inventory[s.curIdx].Title,
		}
		s.name.Scale = engo.Point{X: chara.TextScale.X * 1.05, Y: chara.TextScale.Y * 1.05}
		s.desc.Drawable = common.Text{
			Font:        chara.Font,
			Text:        chara.Inventory[s.curIdx].Description,
			LineSpacing: 0.8,
		}
		s.desc.Scale = engo.Point{X: chara.TextScale.X * 0.95, Y: chara.TextScale.Y * 0.95}
		s.setIdx = s.curIdx
	}
}

func (s *ItemSelectSystem) pause() {
	s.paused = true
	s.cursor.Hidden = true
	s.short1.Hidden = true
	s.short1.Scale = s.scale
	s.short2.Hidden = true
	s.short2.Scale = s.scale
	s.short3.Hidden = true
	s.short3.Scale = s.scale
	s.short4.Hidden = true
	s.short4.Scale = s.scale
	s.name.Hidden = true
	s.name.Scale = engo.Point{X: s.scale.X * 1.05, Y: s.scale.Y * 1.05}
	s.desc.Hidden = true
	s.desc.Scale = engo.Point{X: s.scale.X * 0.95, Y: s.scale.Y * 0.95}
	for _, e := range s.entities {
		e.box.Hidden = true
	}
}

func (s *ItemSelectSystem) unpause() {
	s.paused = false
	s.cursor.Hidden = false
	s.cursor.Position = engo.Point{X: s.short1.Position.X - s.cursor.Width - 2, Y: s.short1.Position.Y + 5}
	s.skipNextFrame = true
	s.short1.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short1.Hidden = false
	s.short2.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short2.Hidden = false
	s.short3.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short3.Hidden = false
	s.short4.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.short4.Hidden = false
	s.name.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.name.Hidden = false
	s.desc.Drawable = common.Text{
		Font: s.fnt,
		Text: "---",
	}
	s.desc.Hidden = false
	for _, e := range s.entities {
		if e.IsCardSelected {
			e.box.Hidden = false
			switch {
			case len(e.Inventory) > 1:
				s.short1.Drawable = common.Text{
					Font: e.Font,
					Text: e.Inventory[0].Shorthand,
				}
				s.short1.Scale = e.TextScale
				s.name.Drawable = common.Text{
					Font: e.Font,
					Text: e.Inventory[0].Title,
				}
				s.name.Scale = e.TextScale
				s.desc.Drawable = common.Text{
					Font: e.Font,
					Text: e.Inventory[0].Description,
				}
				s.desc.Scale = e.TextScale
				fallthrough
			case len(e.Inventory) > 2:
				s.short2.Drawable = common.Text{
					Font: e.Font,
					Text: e.Inventory[1].Shorthand,
				}
				s.short2.Scale = e.TextScale
				fallthrough
			case len(e.Inventory) > 3:
				s.short3.Drawable = common.Text{
					Font: e.Font,
					Text: e.Inventory[2].Shorthand,
				}
				s.short3.Scale = e.TextScale
				fallthrough
			case len(e.Inventory) > 4:
				s.short4.Drawable = common.Text{
					Font: e.Font,
					Text: e.Inventory[3].Shorthand,
				}
				s.short4.Scale = e.TextScale
			}
		}
	}
}
