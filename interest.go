package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type InterestSystemPauseMessage struct {
	Pause bool
}

var InterestSystemPauseMessageType = "Interest System Pause Message"

func (InterestSystemPauseMessage) Type() string { return InterestSystemPauseMessageType }

type InterestComponent struct {
	InterestFunc func()
}

func (c *InterestComponent) GetInterestComponent() *InterestComponent { return c }

type InterestFace interface {
	GetInterestComponent() *InterestComponent
}

type InterestAble interface {
	common.BasicFace
	common.CollisionFace
	InterestFace
}

type interestEntity struct {
	*ecs.BasicEntity
	*common.CollisionComponent
	*InterestComponent
}

type InterestSystem struct {
	entities []interestEntity

	skipNextFrame, paused bool
}

func (s *InterestSystem) New(w *ecs.World) {
	engo.Mailbox.Listen(InterestSystemPauseMessageType, func(message engo.Message) {
		msg, ok := message.(InterestSystemPauseMessage)
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

func (s *InterestSystem) Add(basic *ecs.BasicEntity, collision *common.CollisionComponent, interest *InterestComponent) {
	s.entities = append(s.entities, interestEntity{basic, collision, interest})
}

func (s *InterestSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(InterestAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetCollisionComponent(), o.GetInterestComponent())
}

func (s *InterestSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *InterestSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}
	for _, entity := range s.entities {
		if entity.Collides&CollisionGroupInterest != 0 {
			if engo.Input.Button("A").JustPressed() || (engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft) {
				if entity.InterestFunc != nil {
					entity.InterestFunc()
				}
			}
		}
	}
}

func (s *InterestSystem) pause() {
	s.paused = true
}

func (s *InterestSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
}
