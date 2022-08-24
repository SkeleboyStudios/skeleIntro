package main

import (
	"github.com/EngoEngine/ecs"
)

type CastBarComponent struct {
	barHP, barMP                   float32
	totalCastTime, currentCastTime float32
	isCasting                      bool
}

type barEntity struct {
	chara  *Character
	baddie *Baddie
}

type BarSystem struct {
	entities []barEntity
}

func (s *BarSystem) Add(chara *Character, bad *Baddie) {
	s.entities = append(s.entities, barEntity{chara, bad})
}

func (s *BarSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(Characterable)
	if ok {
		s.Add(o.GetCharacter(), nil)
	}
	o2, ok := i.(Baddieable)
	if ok {
		s.Add(nil, o2.GetBaddie())
	}
}

func (s *BarSystem) Remove(b ecs.BasicEntity) {
	d := -1
	for i, e := range s.entities {
		if e.chara != nil {
			if e.chara.BasicEntity.ID() == b.ID() {
				d = i
				break
			}
		} else if e.baddie != nil {
			if e.baddie.BasicEntity.ID() == b.ID() {
				d = i
				break
			}
		}
	}
	if d >= 0 {
		s.entities = append(s.entities[:d], s.entities[d+1:]...)
	}
}

func (s *BarSystem) Update(dt float32) {
	for _, e := range s.entities {
		if e.chara != nil {
			if e.chara.barHP != e.chara.HP {
				e.chara.barHP = e.chara.HP
				e.chara.hpBar.Width = 83 * (e.chara.barHP / e.chara.MaxHP)
			}
			if e.chara.barMP != e.chara.MP {
				e.chara.barMP = e.chara.MP
				e.chara.mpBar.Width = 83 * (e.chara.barMP / e.chara.MaxMP)
			}
			if e.chara.isCasting {
				e.chara.currentCastTime += dt
				if e.chara.currentCastTime >= e.chara.totalCastTime {
					e.chara.currentCastTime = 0
					e.chara.totalCastTime = 1
					if e.chara.SelectedAbility.EffectFunc != nil {
						e.chara.SelectedAbility.EffectFunc(e.chara, e.chara.TargetPlayers, e.chara.TargetBaddies)
						e.chara.SelectedAbility = Ability{}
						e.chara.TargetPlayers = make([]*Character, 0)
						e.chara.TargetBaddies = make([]*Baddie, 0)
					}
				}
				e.chara.castBar.Width = 83 * (e.chara.currentCastTime / e.chara.totalCastTime)
			} else {
				e.chara.castBar.Width = 0
			}
		} else if e.baddie != nil {

		}
	}
}
