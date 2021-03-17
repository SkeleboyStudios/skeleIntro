package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type DoorSystemPauseMessage struct {
	Pause bool
}

var DoorSystemPauseMessageType = "Door System Pause Message"

func (DoorSystemPauseMessage) Type() string { return DoorSystemPauseMessageType }

type DoorComponent struct {
	IsOpen     bool
	DoorButton string
	TeleportTo engo.Point
	OpenFrame  int
}

func (c *DoorComponent) GetDoorComponent() *DoorComponent {
	return c
}

type DoorFace interface {
	GetDoorComponent() *DoorComponent
}

type Doorable interface {
	common.BasicFace
	common.AnimationFace
	common.CollisionFace
	DoorFace
}

type DoorEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
	*common.CollisionComponent
	*DoorComponent
}

type DoorSystem struct {
	entities []DoorEntity

	skipNextFrame, paused bool
}

func (s *DoorSystem) New(w *ecs.World) {
	engo.Mailbox.Listen(DoorSystemPauseMessageType, func(message engo.Message) {
		msg, ok := message.(DoorSystemPauseMessage)
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

func (s *DoorSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent, collision *common.CollisionComponent, door *DoorComponent) {
	s.entities = append(s.entities, DoorEntity{basic, anim, collision, door})
}

func (s *DoorSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(Doorable)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetAnimationComponent(), o.GetCollisionComponent(), o.GetDoorComponent())
}

func (s *DoorSystem) Remove(basic ecs.BasicEntity) {
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

func (s *DoorSystem) Update(dt float32) {
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.paused {
		return
	}
	for i, entity := range s.entities {
		if entity.Collides&CollisionGroupDoor != 0 {
			if entity.IsOpen {
				if entity.CurrentFrame >= entity.OpenFrame && engo.Input.Button(entity.DoorButton).Down() {
					engo.Mailbox.Dispatch(TeleportPlayerMessage{Pt: entity.TeleportTo})
				}
			} else {
				entity.SelectAnimationByName("open")
				s.entities[i].IsOpen = true
			}
		} else {
			if entity.IsOpen {
				entity.SelectAnimationByName("close")
				s.entities[i].IsOpen = false
			}
		}
	}
}

func (s *DoorSystem) pause() {
	s.paused = true
}

func (s *DoorSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
}
