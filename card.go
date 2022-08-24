package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

var CardSelectSystemPauseMessageType = "CardSelectSystemPauseMessage"

type CardSelectSystemPauseMessage struct {
	Pause bool
}

func (CardSelectSystemPauseMessage) Type() string { return CardSelectSystemPauseMessageType }

type cardSelectEntity struct {
	*Character
}

type CardSelectSystem struct {
	entities              []cardSelectEntity
	paused, skipNextFrame bool
	setIdx, curIdx        int
}

func (s *CardSelectSystem) New(w *ecs.World) {
	s.setIdx = -1
	engo.Mailbox.Listen(CardSelectSystemPauseMessageType, func(message engo.Message) {
		msg, ok := message.(CardSelectSystemPauseMessage)
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

func (s *CardSelectSystem) Add(chara *Character) {
	s.entities = append(s.entities, cardSelectEntity{chara})
}

func (s *CardSelectSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(Characterable)
	if ok {
		s.Add(o.GetCharacter())
	}
}

func (s *CardSelectSystem) Remove(b ecs.BasicEntity) {
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

func (s *CardSelectSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}
	if s.setIdx != s.curIdx {
		if s.setIdx >= 0 {
			s.entities[s.setIdx].IsCardSelected = false
			s.entities[s.curIdx].IsCardSelected = false
			s.entities[s.setIdx].MoveCard(engo.Point{X: s.entities[s.setIdx].card.Position.X, Y: s.entities[s.setIdx].card.Position.Y + 10})
		}
		s.entities[s.curIdx].IsCardSelected = true
		s.entities[s.curIdx].MoveCard(engo.Point{X: s.entities[s.curIdx].card.Position.X, Y: s.entities[s.curIdx].card.Position.Y - 10})
		s.setIdx = s.curIdx
	}
	if engo.Input.Button("left").JustPressed() {
		s.curIdx--
		if s.curIdx < 0 {
			s.curIdx = len(s.entities) - 1
		}
	} else if engo.Input.Button("right").JustPressed() {
		s.curIdx++
		if s.curIdx > len(s.entities)-1 {
			s.curIdx = 0
		}
	}
	if engo.Input.Button("A").JustPressed() {
		engo.Mailbox.Dispatch(PhaseSetMessage{
			Phase: AbilitySelectPhase,
		})
		engo.Mailbox.Dispatch(PhaseDequeuMessage{})
	} else if engo.Input.Button("B").JustPressed() {
		engo.Mailbox.Dispatch(PhaseSetMessage{
			Phase: ItemSelectPhase,
		})
		engo.Mailbox.Dispatch(PhaseDequeuMessage{})
	} else if engo.Input.Button("X").JustPressed() {
		s.entities[s.setIdx].SelectedAbility = RegularAttackAbility
		engo.Mailbox.Dispatch(PhaseSetMessage{
			Phase: TargetPhase,
		})
		engo.Mailbox.Dispatch(PhaseDequeuMessage{})
	} else if engo.Input.Button("Y").JustPressed() {
		s.entities[s.setIdx].SelectedAbility = DefendAbility
		engo.Mailbox.Dispatch(PhaseSetMessage{
			Phase: CardSelectPhase,
		})
		engo.Mailbox.Dispatch(PhaseDequeuMessage{})
	}
}

func (s *CardSelectSystem) pause() {
	s.paused = true
	for _, e := range s.entities {
		e.card.Hidden = true
		e.cardText.Hidden = true
		e.hpBar.Hidden = true
		e.mpBar.Hidden = true
		e.castBar.Hidden = true
		e.totalCastTime = 1
		e.currentCastTime = 0
	}
	if s.setIdx >= 0 {
		s.entities[s.setIdx].MoveCard(engo.Point{X: s.entities[s.setIdx].card.Position.X, Y: s.entities[s.setIdx].card.Position.Y + 10})
	}
	s.setIdx = -1
	s.curIdx = 0
}

func (s *CardSelectSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
	for _, e := range s.entities {
		e.card.Hidden = false
		e.cardText.Hidden = false
		e.hpBar.Hidden = false
		e.mpBar.Hidden = false
		e.castBar.Hidden = false
		e.IsCardSelected = false
	}
	s.entities[0].IsCardSelected = true
}
