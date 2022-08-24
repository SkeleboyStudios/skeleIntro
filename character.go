package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type CharacterInfo struct {
	Name          string
	CardSprite    common.Drawable
	BoxSprite     common.Drawable
	CardTextScale engo.Point
	HP, MP        float32
	MaxHP, MaxMP  float32
	Str, Def      float32
	Dex, Int      float32
	Font          *common.Font
	Clip          *common.Player
}

type Item struct {
	Title       string
	Shorthand   string
	Description string
	Quantity    int
	TargetType  Target
	EffectFunc  func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie)
}

type Target uint

const (
	TargetTypeNone Target = iota
	TargetTypeSingleEnemy
	TargetTypeAllEnemy
	TargetTypeSingleFriend
	TargetTypeAllFriend
	TargetTypeSingleAny
	TargetTypeAll
)

type StatsComponent struct {
	HP, MP       float32
	MaxHP, MaxMP float32
	Str, Def     float32
	Dex, Int     float32
}

type StatusComponent struct {
	Name           string
	IsCardSelected bool
	TextScale      engo.Point
}

type TargetComponent struct {
	TargetPlayers []*Character
	TargetBaddies []*Baddie
}

type AbilityComponent struct {
	Abilities       []Ability
	SelectedAbility Ability
	IsAbilitySelected bool
}

type InventoryComponent struct {
	Inventory    []Item
	SelectedItem Item
	IsItemSelected bool
}

type ChatComponent struct {
	Font *common.Font
	Clip *common.Player
	Dot  *common.Drawable
}

type Characterface interface {
	GetCharacter() *Character
}

type Characterable interface {
	Characterface
}

type Character struct {
	card     *sprite
	cardText *sprite
	box      *sprite
	hpBar    *sprite
	mpBar    *sprite
	castBar  *sprite

	ecs.BasicEntity

	StatsComponent
	CastBarComponent
	StatusComponent
	TargetComponent
	AbilityComponent
	InventoryComponent
	ChatComponent
}

func (c *Character) GetCharacter() *Character {
	return c
}

func (c *Character) AddAbility(a Ability) {
	for _, a2 := range c.Abilities {
		if a.Title == a2.Title {
			return
		}
	}
	c.Abilities = append(c.Abilities, a)
}

func (c *Character) RemoveAbility(title string) {
	idx := -1
	for i, a := range c.Abilities {
		if a.Title == title {
			idx = i
			break
		}
	}
	if idx >= 0 {
		c.Abilities = append(c.Abilities[:idx], c.Abilities[idx+1:]...)
	}
}

func (c *Character) AddItem(i Item) {
	for _, i2 := range c.Inventory {
		if i.Title == i2.Title {
			i2.Quantity += i.Quantity
			return
		}
	}
	c.Inventory = append(c.Inventory, i)
}

func (c *Character) RemoveItem(title string, quantity int) {
	for i, item := range c.Inventory {
		if item.Title == title {
			item.Quantity -= quantity
			if item.Quantity <= 0 {
				c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			}
			return
		}
	}
}

func (c *Character) MoveCard(p engo.Point) {
	c.card.Position = p
	c.cardText.Position = engo.Point{X: p.X + 10, Y: p.Y + 3}
	c.hpBar.Position = engo.Point{X: p.X + 8, Y: p.Y + 32}
	c.mpBar.Position = engo.Point{X: p.X + 8, Y: p.Y + 54}
	c.castBar.Position = engo.Point{X: p.X + 8, Y: p.Y + 76}
}

func AddCharacter(info CharacterInfo, w *ecs.World) *Character {
	chara := &Character{BasicEntity: ecs.NewBasic()}
	chara.Name = info.Name
	chara.TextScale = info.CardTextScale
	chara.totalCastTime = 1
	chara.StatsComponent = StatsComponent{
		HP:    info.HP,
		MP:    info.MP,
		MaxHP: info.MaxHP,
		MaxMP: info.MaxMP,
		Str:   info.Str,
		Def:   info.Def,
		Int:   info.Int,
		Dex:   info.Dex,
	}
	chara.card = &sprite{BasicEntity: ecs.NewBasic()}
	chara.card.Drawable = info.CardSprite
	chara.card.Height = chara.card.Drawable.Height()
	chara.card.Width = chara.card.Drawable.Width()
	w.AddEntity(chara.card)
	chara.cardText = &sprite{BasicEntity: ecs.NewBasic()}
	chara.cardText.Drawable = common.Text{
		Text: info.Name,
		Font: info.Font,
	}
	chara.cardText.Scale = info.CardTextScale
	chara.cardText.SetZIndex(1)
	w.AddEntity(chara.cardText)
	chara.hpBar = &sprite{BasicEntity: ecs.NewBasic()}
	chara.hpBar.Drawable = common.Rectangle{}
	chara.hpBar.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	chara.hpBar.Width = 83
	chara.hpBar.Height = 13
	chara.hpBar.SetZIndex(1)
	w.AddEntity(chara.hpBar)
	chara.mpBar = &sprite{BasicEntity: ecs.NewBasic()}
	chara.mpBar.Drawable = common.Rectangle{}
	chara.mpBar.Color = color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	chara.mpBar.Width = 83
	chara.mpBar.Height = 13
	chara.mpBar.SetZIndex(1)
	w.AddEntity(chara.mpBar)
	chara.castBar = &sprite{BasicEntity: ecs.NewBasic()}
	chara.castBar.Drawable = common.Rectangle{}
	chara.castBar.Color = color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}
	chara.castBar.Width = 83
	chara.castBar.Height = 13
	chara.castBar.SetZIndex(1)
	w.AddEntity(chara.castBar)
	chara.box = &sprite{BasicEntity: ecs.NewBasic()}
	chara.box.Drawable = info.BoxSprite
	chara.box.Height = chara.box.Drawable.Height()
	chara.box.Width = chara.box.Drawable.Width()
	chara.box.Position = engo.Point{X: 15, Y: 360 - chara.box.Height - 10}
	chara.box.SetZIndex(2)
	chara.box.Hidden = true
	w.AddEntity(chara.box)
	chara.Font = info.Font
	chara.Clip = info.Clip
	w.AddEntity(chara)
	return chara
}
