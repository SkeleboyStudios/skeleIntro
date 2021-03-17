package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type sprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type selection struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	CursorComponent
}

type animation struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	common.AnimationComponent
}

type audio struct {
	ecs.BasicEntity
	common.AudioComponent
}

type wall struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.CollisionComponent
}

type wallPtr struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.CollisionComponent
}

type playa struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	common.AnimationComponent
	common.CollisionComponent

	MoveComponent
}

type door struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
	common.AnimationComponent

	DoorComponent
}

type doorPtr struct {
	*ecs.BasicEntity
	*common.CollisionComponent
	*common.RenderComponent
	*common.SpaceComponent
	*common.AnimationComponent

	*DoorComponent
}

type interest struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent

	InterestComponent
}

type interestPtr struct {
	*ecs.BasicEntity
	*common.CollisionComponent
	*common.RenderComponent
	*common.SpaceComponent

	*InterestComponent
}

type room struct {
	bg        sprite
	walls     []wallPtr
	doors     []doorPtr
	interests []interestPtr
}

type doorInfo struct {
	URL                                              string
	Position, TeleportTo                             engo.Point
	CellWidth, CellHeight, BorderWidth, BorderHeight int
	Shapes                                           []common.Shape
	OpenFrames, CloseFrames                          []int
	Button                                           string
}

type wallInfo struct {
	Position      engo.Point
	Width, Height float32
	Shapes        []common.Shape
}

type interestInfo struct {
	URL      string
	Position engo.Point
	Shapes   []common.Shape
	Func     func()
}

func newRoom(w *ecs.World, start engo.Point, bgPath string, wallInfos []wallInfo, doorInfos []doorInfo, interestInfos []interestInfo) room {
	ret := room{bg: sprite{BasicEntity: ecs.NewBasic()}}
	ret.bg.Drawable, _ = common.LoadedSprite(bgPath)
	ret.bg.Position = start
	w.AddEntity(&ret.bg)
	for _, wallInfo := range wallInfos {
		wa := wall{BasicEntity: ecs.NewBasic()}
		wa.Position = engo.Point{
			X: start.X + wallInfo.Position.X,
			Y: start.Y + wallInfo.Position.Y,
		}
		wa.Width = wallInfo.Width
		wa.Height = wallInfo.Height
		for _, shape := range wallInfo.Shapes {
			wa.AddShape(shape)
		}
		wa.CollisionComponent = common.CollisionComponent{Group: CollisionGroupPlayaWall}
		ret.walls = append(ret.walls, wallPtr{wa.GetBasicEntity(), wa.GetSpaceComponent(), wa.GetCollisionComponent()})
		w.AddEntity(&wa)
	}
	for _, doorInfo := range doorInfos {
		d := door{BasicEntity: ecs.NewBasic()}
		dSS := common.NewSpritesheetWithBorderFromFile(doorInfo.URL, doorInfo.CellWidth, doorInfo.CellHeight, doorInfo.BorderWidth, doorInfo.BorderHeight)
		d.Drawable = dSS.Drawable(0)
		d.SetZIndex(1)
		d.Position = engo.Point{
			X: start.X + doorInfo.Position.X,
			Y: start.Y + doorInfo.Position.Y,
		}
		d.Scale = engo.Point{X: 2, Y: 2}
		for _, shape := range doorInfo.Shapes {
			d.AddShape(shape)
		}
		d.AnimationComponent = common.NewAnimationComponent(dSS.Drawables(), 0.1)
		d.AddAnimation(&common.Animation{Name: "open", Frames: doorInfo.OpenFrames})
		d.AddAnimation(&common.Animation{Name: "close", Frames: doorInfo.CloseFrames})
		d.CollisionComponent = common.CollisionComponent{Main: CollisionGroupDoor}
		d.DoorButton = doorInfo.Button
		d.TeleportTo = doorInfo.TeleportTo
		d.OpenFrame = doorInfo.OpenFrames[len(doorInfo.OpenFrames)-1]
		ret.doors = append(ret.doors, doorPtr{d.GetBasicEntity(), d.GetCollisionComponent(), d.GetRenderComponent(), d.GetSpaceComponent(), d.GetAnimationComponent(), d.GetDoorComponent()})
		w.AddEntity(&d)
	}
	for _, interestInfo := range interestInfos {
		i := interest{BasicEntity: ecs.NewBasic()}
		i.Drawable, _ = common.LoadedSprite(interestInfo.URL)
		i.SetZIndex(interestInfo.Position.Y + start.Y + (i.Drawable.Height() / 2) - 82) //y position (of middle!!!) minus 82 for height of player
		i.Position = engo.Point{
			X: start.X + interestInfo.Position.X,
			Y: start.Y + interestInfo.Position.Y,
		}
		for _, shape := range interestInfo.Shapes {
			i.AddShape(shape)
		}
		i.CollisionComponent = common.CollisionComponent{Main: CollisionGroupInterest}
		i.InterestFunc = interestInfo.Func
		ret.interests = append(ret.interests, interestPtr{i.GetBasicEntity(), i.GetCollisionComponent(), i.GetRenderComponent(), i.GetSpaceComponent(), i.GetInterestComponent()})
		w.AddEntity(&i)
	}
	return ret
}
