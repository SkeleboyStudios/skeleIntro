package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

var TargetSystemPauseMessageType = "TargetSystemPauseMessage"

type TargetSystemPauseMessage struct {
	Pause bool
}

func (TargetSystemPauseMessage) Type() string { return TargetSystemPauseMessageType }

type targetEntity struct {
	chara  *Character
	baddie *Baddie
}

type TargetSystem struct {
	entities              []targetEntity
	paused, skipNextFrame bool
}

func (s *TargetSystem) New(w *ecs.World) {
	engo.Mailbox.Listen(TargetSystemPauseMessageType, func(message engo.Message) {
		msg, ok := message.(TargetSystemPauseMessage)
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

func (s *TargetSystem) Add(chara *Character, baddie *Baddie) {
	s.entities = append(s.entities, targetEntity{chara, baddie})
}

func (s *TargetSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(Characterable)
	if ok {
		s.Add(o.GetCharacter(), nil)
	}
	o2, ok := i.(Baddieable)
	if ok {
		s.Add(nil, o2.GetBaddie())
	}
}

func (s *TargetSystem) Remove(b ecs.BasicEntity) {
	d := -1
	for i, e := range s.entities {
		if e.chara != nil {
			if e.chara.ID() == b.ID() {
				d = 1
				break
			}
		}
		if e.baddie != nil {
			if e.baddie.ID() == b.ID() {
				d = i
				break
			}
		}
	}
	if d >= 0 {
		s.entities = append(s.entities[:d], s.entities[d+1:]...)
	}
}

func (s *TargetSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}

	var chara *Character
	for _, e := range s.entities {
		if e.chara != nil {
			if e.chara.IsCardSelected {
				chara = e.chara
				break
			}
		}
	}
	if chara == nil {
		return
	}

}

func (s *TargetSystem) pause() {
	s.paused = true
}

func (s *TargetSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
}
