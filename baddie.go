package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type BaddieState struct {
	PhaseStartFunc func(bad *Baddie)
	PhaseEndFunc   func(bad *Baddie)
}

type BaddieInfo struct {
	Name         string
	Spritesheet  string
	HP, MP       float32
	MaxHP, MaxMP float32
	Str, Def     float32
	Dex, Int     float32
	Font         *common.Font
	Clip         *common.Player
	Phases       map[string]BaddieState
	Attacks      []Attack
	StartPhase   string
}

type BaddieComponent struct {
	Spritesheet *common.Spritesheet
	Phases      map[string]BaddieState
}

type BaddieFace interface {
	GetBaddie() *Baddie
}

type Baddieable interface {
	BaddieFace
}

type Baddie struct {
	spr     *sprite
	hpBar   *sprite
	castBar *sprite

	ecs.BasicEntity

	StatsComponent
	CastBarComponent
	TargetComponent
	BaddieComponent
	ChatComponent
	AIComponent
}

func AddBaddie(bad BaddieInfo) *Baddie {
	return nil
}
