package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type TeleportPlayerMessage struct {
	Pt engo.Point
}

var TeleportPlayerMessageType = "Teleport Player Message"

func (TeleportPlayerMessage) Type() string { return TeleportPlayerMessageType }

type MoveComponent struct {
	PlayerCharacter bool
	Speed           float32

	velocity      engo.Point
	currentZIndex float32
}

func (c *MoveComponent) GetMoveComponent() *MoveComponent {
	return c
}

type MoveFace interface {
	GetMoveComponent() *MoveComponent
}

type Moveable interface {
	common.BasicFace
	common.AnimationFace
	common.SpaceFace
	common.RenderFace
	MoveFace
}

type moveEntity struct {
	*ecs.BasicEntity
	*common.AnimationComponent
	*common.SpaceComponent
	*common.RenderComponent
	*MoveComponent
}

type MoveSystem struct {
	entities []moveEntity
}

func (s *MoveSystem) New(w *ecs.World) {
	engo.Mailbox.Listen(TeleportPlayerMessageType, func(message engo.Message) {
		msg, ok := message.(TeleportPlayerMessage)
		if !ok {
			return
		}
		for i := 0; i < len(s.entities); i++ {
			if s.entities[i].PlayerCharacter {
				s.entities[i].Position = msg.Pt
			}
		}
	})
}

func (s *MoveSystem) Add(basic *ecs.BasicEntity, anim *common.AnimationComponent, space *common.SpaceComponent, render *common.RenderComponent, move *MoveComponent) {
	move.velocity = engo.Point{}
	s.entities = append(s.entities, moveEntity{basic, anim, space, render, move})
}

func (s *MoveSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(Moveable)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetAnimationComponent(), o.GetSpaceComponent(), o.GetRenderComponent(), o.GetMoveComponent())
}

func (s *MoveSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities[delete].SelectAnimationByName("downstop")
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *MoveSystem) Update(dt float32) {
	for i, entity := range s.entities {
		if entity.PlayerCharacter {
			if v, changed := s.getSpeed(); changed {
				v, _ = v.Normalize()
				v.MultiplyScalar(dt * entity.Speed)
				s.entities[i].velocity = v
			}
			s.setAnimation(entity)
			entity.Position.Add(entity.velocity)
		} else {
			//something else
		}
		if entity.currentZIndex != entity.Position.Y {
			entity.SetZIndex(entity.Position.Y)
			entity.currentZIndex = entity.Position.Y
		}
	}
}

func (s *MoveSystem) setAnimation(entity moveEntity) {
	if engo.Input.Button("left").JustPressed() {
		entity.AnimationComponent.SelectAnimationByName("left")
	} else if engo.Input.Button("right").JustPressed() {
		entity.AnimationComponent.SelectAnimationByName("right")
	} else if engo.Input.Button("up").JustPressed() {
		entity.AnimationComponent.SelectAnimationByName("up")
	} else if engo.Input.Button("down").JustPressed() {
		entity.AnimationComponent.SelectAnimationByName("down")
	}

	if engo.Input.Button("up").JustReleased() {
		entity.AnimationComponent.SelectAnimationByName("upstop")
		if engo.Input.Button("left").Down() {
			entity.AnimationComponent.SelectAnimationByName("left")
		} else if engo.Input.Button("right").Down() {
			entity.AnimationComponent.SelectAnimationByName("right")
		} else if engo.Input.Button("up").Down() {
			entity.AnimationComponent.SelectAnimationByName("up")
		} else if engo.Input.Button("down").Down() {
			entity.AnimationComponent.SelectAnimationByName("down")
		}
	} else if engo.Input.Button("down").JustReleased() {
		entity.AnimationComponent.SelectAnimationByName("downstop")
		if engo.Input.Button("left").Down() {
			entity.AnimationComponent.SelectAnimationByName("left")
		} else if engo.Input.Button("right").Down() {
			entity.AnimationComponent.SelectAnimationByName("right")
		} else if engo.Input.Button("up").Down() {
			entity.AnimationComponent.SelectAnimationByName("up")
		} else if engo.Input.Button("down").Down() {
			entity.AnimationComponent.SelectAnimationByName("down")
		}
	} else if engo.Input.Button("left").JustReleased() {
		entity.AnimationComponent.SelectAnimationByName("leftstop")
		if engo.Input.Button("left").Down() {
			entity.AnimationComponent.SelectAnimationByName("left")
		} else if engo.Input.Button("right").Down() {
			entity.AnimationComponent.SelectAnimationByName("right")
		} else if engo.Input.Button("up").Down() {
			entity.AnimationComponent.SelectAnimationByName("up")
		} else if engo.Input.Button("down").Down() {
			entity.AnimationComponent.SelectAnimationByName("down")
		}
	} else if engo.Input.Button("right").JustReleased() {
		entity.AnimationComponent.SelectAnimationByName("rightstop")
		if engo.Input.Button("left").Down() {
			entity.AnimationComponent.SelectAnimationByName("left")
		} else if engo.Input.Button("right").Down() {
			entity.AnimationComponent.SelectAnimationByName("right")
		} else if engo.Input.Button("up").Down() {
			entity.AnimationComponent.SelectAnimationByName("up")
		} else if engo.Input.Button("down").Down() {
			entity.AnimationComponent.SelectAnimationByName("down")
		}
	}
}

func (s *MoveSystem) getSpeed() (p engo.Point, changed bool) {
	p.X = engo.Input.Axis("horizontal").Value()
	p.Y = engo.Input.Axis("vertical").Value()
	origX, origY := p.X, p.Y

	if engo.Input.Button("up").JustPressed() {
		p.Y = -1
	} else if engo.Input.Button("down").JustPressed() {
		p.Y = 1
	}
	if engo.Input.Button("left").JustPressed() {
		p.X = -1
	} else if engo.Input.Button("right").JustPressed() {
		p.X = 1
	}

	if engo.Input.Button("up").JustReleased() || engo.Input.Button("down").JustReleased() {
		p.Y = 0
		changed = true
		if engo.Input.Button("up").Down() {
			p.Y = -1
		} else if engo.Input.Button("down").Down() {
			p.Y = 1
		} else if engo.Input.Button("left").Down() {
			p.X = -1
		} else if engo.Input.Button("right").Down() {
			p.X = 1
		}
	}
	if engo.Input.Button("left").JustReleased() || engo.Input.Button("right").JustReleased() {
		p.X = 0
		changed = true
		if engo.Input.Button("left").Down() {
			p.X = -1
		} else if engo.Input.Button("right").Down() {
			p.X = 1
		} else if engo.Input.Button("up").Down() {
			p.Y = -1
		} else if engo.Input.Button("down").Down() {
			p.Y = 1
		}
	}
	changed = changed || p.X != origX || p.Y != origY
	return
}
