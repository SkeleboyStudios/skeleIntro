package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type Character struct {
	Name      string
	Sprite    common.Drawable
	HP, MP    float32
	Str, Def  float32
	Dex, Int  float32
	Abilities []*Ability
	Font      *common.Font
}

type Ability struct {
	Title       string
	Shorthand   string
	Description string
	ActionCost  int
	MPCost      float32
	TargetType  Target
	EffectFunc  func(You *Player, TargetPlayers []*Player, TargetBaddies []*Baddie)
}

type Target uint

const (
	TargetTypeNone Target = iota
	TargetTypeSingleEnemey
	TargetTypeAllEnemy
	TargetTypeSingleFriend
	TargetTypeAllFriend
	TargetTypeSingleAny
	TargetTypeAll
)

type Player struct {
	ecs.BasicEntity

	HPComponent
	MPComponent
	StaminaBarComponent
	StatusComponent
	TargetComponent
	AttackComponent
}

func (p *Player) RemoveAbility(title string) {
	delete := -1
	for i, a := range p.Abilities {
		if a.Title == title {
			idx = i
			break
		}
	}
	if delete >= 0 {
		p.Abilities = append(p.Abilities[:delete], p.Abilities[delete+1:]...)
	}
}
