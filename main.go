package main

import (
	"bytes"
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"github.com/gopherjs/gopherjs/js"

	"github.com/SkeleboyStudios/skeleIntro/assets"
)

var (
	centerAnimation *common.Animation
	width           = 580
	height          = 370
	dpi             float32
)

type anim struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
}

type OpeningScene struct{}

func (*OpeningScene) Type() string { return "Opening Scene" }

func (*OpeningScene) Preload() {
	b, err := assets.Asset("welcome.png")
	if err != nil {
		panic("no welcome.png found")
	}
	engo.Files.LoadReaderData("welcome.png", bytes.NewReader(b))

	centerAnimation = &common.Animation{Name: "center", Frames: []int{0, 1, 2}, Loop: true}
}

func (*OpeningScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.AnimationSystem{})

	welcomeSheet := common.NewSpritesheetFromFile("welcome.png", width, height)

	centerEntity := &anim{BasicEntity: ecs.NewBasic()}
	centerEntity.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: engo.ResizeXOffset / (2 * dpi),
			Y: engo.ResizeYOffset / (2 * dpi),
		},
		Width:  float32(width),
		Height: float32(height),
	}
	centerEntity.RenderComponent = common.RenderComponent{
		Drawable: welcomeSheet.Cell(0),
		Scale:    engo.Point{1, 1},
	}
	centerEntity.AnimationComponent = common.NewAnimationComponent(welcomeSheet.Drawables(), 0.1)
	centerEntity.AnimationComponent.AddAnimation(centerAnimation)
	centerEntity.AnimationComponent.AddDefaultAnimation(centerAnimation)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&centerEntity.BasicEntity, &centerEntity.RenderComponent, &centerEntity.SpaceComponent)
		case *common.AnimationSystem:
			sys.Add(&centerEntity.BasicEntity, &centerEntity.AnimationComponent, &centerEntity.RenderComponent)
		}
	}
}

func main() {
	x := js.Global.Get("document").Get("body").Get("clientWidth").Int()
	println(x)
	y := js.Global.Get("document").Get("body").Get("clientHeight").Int()
	println(y)
	dpi = float32(js.Global.Get("devicePixelRatio").Float())
	opts := engo.RunOptions{
		Title:  "Animation Demo",
		Width:  x,
		Height: y,
		GlobalScale: engo.Point{
			X: dpi * float32(x) / float32(width),
			Y: dpi * float32(y) / float32(width),
		},
	}
	engo.Run(opts, &OpeningScene{})
}
