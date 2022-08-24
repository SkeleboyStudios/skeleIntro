package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/SkeleboyStudios/skeleIntro/assets"

	"github.com/Noofbiz/pixelshader"
)

type GhostFightScene struct {
	PlayerLocation engo.Point
	files          []string
}

func (*GhostFightScene) Type() string { return "Ghost Fight!!!" }

func (s *GhostFightScene) Preload() {
	s.files = []string{
		"title/move.wav",
		"title/cursor.png",
		"title/log.wav",
		"title/log.ttf",
		"fight/log.png",
		"fight/dots.png",
		"fight/log.ttf",
		"fight/bg.ogg",
		"fight/log.wav",
		"fight/bg.png",
		"fight/cards.png",
		"fight/cash.wav",
		"fight/mimic.png",
		"fight/you.ttf",
		"fight/boxes.png",
		"fight/me.ttf",
		"fight/len.ttf",
		"fight/me.wav",
		"fight/len.wav",
	}

	for _, file := range s.files {
		data, err := assets.Asset(file)
		if err != nil {
			log.Fatalf("Unable to locate asset with URL: %v\n", file)
		}
		err = engo.Files.LoadReaderData(file, bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Unable to load asset with URL: %v\n At %v", file, s.Type())
		}
	}

	engo.Input.RegisterButton("up", engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton("down", engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton("left", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("right", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("A", engo.KeyJ, engo.KeyZ)
	engo.Input.RegisterButton("B", engo.KeyK, engo.KeyX)
	engo.Input.RegisterButton("X", engo.KeyL, engo.KeyC)
	engo.Input.RegisterButton("Y", engo.KeySemicolon, engo.KeyV)
	engo.Input.RegisterButton("FullScreen", engo.KeyFour, engo.KeyF4)
	engo.Input.RegisterButton("Exit", engo.KeyEscape)
}

func (s *GhostFightScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	rand.Seed(time.Now().UnixNano())

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var animatable *common.Animationable
	var notanimatable *common.NotAnimationable
	var animSys = &common.AnimationSystem{}
	w.AddSystemInterface(animSys, animatable, notanimatable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	var audioSys = &common.AudioSystem{}
	w.AddSystemInterface(audioSys, audioable, notaudioable)

	var cursorable *CursorAble
	var notcursorable *NotCursorAble
	var curSys CursorSystem
	curSys.ClickSoundURL = "title/move.wav"
	curSys.CursorURL = "title/cursor.png"
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	w.AddSystem(&CombatLogSystem{
		BackgroundURL: "fight/log.png",
		DotURL:        "fight/dots.png",
		FontURL:       "fight/log.ttf",
		LineDelay:     0.3,
		LetterDelay:   0.1,
	})

	w.AddSystem(&FullScreenSystem{})
	w.AddSystem(&ExitSystem{})

	var characterable *Characterable
	w.AddSystemInterface(&BarSystem{}, characterable, nil)
	w.AddSystemInterface(&CardSelectSystem{}, characterable, nil)

	var phaseable *common.BasicFace
	w.AddSystemInterface(&PhaseSystem{}, phaseable, nil)

	selFont := &common.Font{
		Size: 64,
		FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
		URL:  "fight/log.ttf",
	}
	selFont.CreatePreloaded()

	w.AddSystemInterface(&AbilitySelectSystem{fnt: selFont}, characterable, nil)
	w.AddSystemInterface(&ItemSelectSystem{fnt: selFont}, characterable, nil)

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer("fight/bg.ogg")
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	logSnd := audio{BasicEntity: ecs.NewBasic()}
	logPlayer, _ := common.LoadedPlayer("fight/log.wav")
	logSnd.AudioComponent = common.AudioComponent{Player: logPlayer}
	logSnd.AudioComponent.Player.SetVolume(0.15)
	w.AddEntity(&logSnd)

	youSnd := audio{BasicEntity: ecs.NewBasic()}
	youPlayer, _ := common.LoadedPlayer("title/log.wav")
	youSnd.AudioComponent = common.AudioComponent{Player: youPlayer}
	youSnd.AudioComponent.Player.SetVolume(0.15)
	w.AddEntity(&youSnd)

	meSnd := audio{BasicEntity: ecs.NewBasic()}
	mePlayer, _ := common.LoadedPlayer("fight/me.wav")
	meSnd.AudioComponent = common.AudioComponent{Player: mePlayer}
	meSnd.AudioComponent.Player.SetVolume(0.15)
	w.AddEntity(&meSnd)

	lenSnd := audio{BasicEntity: ecs.NewBasic()}
	lenPlayer, _ := common.LoadedPlayer("fight/len.wav")
	lenSnd.AudioComponent = common.AudioComponent{Player: lenPlayer}
	lenSnd.AudioComponent.Player.SetVolume(0.15)
	w.AddEntity(&lenSnd)

	bg := sprite{BasicEntity: ecs.NewBasic()}
	tex0, _ := common.LoadedSprite("fight/bg.png")
	bg.Drawable = pixelshader.PixelRegion{
		Tex0: tex0,
	}
	bg.SetShader(fightShader)
	bg.SetZIndex(0)
	w.AddEntity(&bg)

	cards := common.NewSpritesheetWithBorderFromFile("fight/cards.png", 102, 105, 1, 1)
	boxes := common.NewSpritesheetWithBorderFromFile("fight/boxes.png", 600, 144, 1, 1)
	youFnt := &common.Font{
		Size: 64,
		FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
		URL:  "fight/you.ttf",
	}
	youFnt.CreatePreloaded()
	you := AddCharacter(CharacterInfo{
		Name:          "You",
		CardSprite:    cards.Drawable(0),
		BoxSprite:     boxes.Drawable(0),
		HP:            100,
		MaxHP:         100,
		MP:            100,
		MaxMP:         100,
		Str:           25,
		Def:           25,
		Dex:           35,
		Int:           40,
		Font:          youFnt,
		Clip:          youPlayer,
		CardTextScale: engo.Point{X: 0.35, Y: 0.35},
	}, w)
	you.MoveCard(engo.Point{X: 320 - you.card.Width/2, Y: 360 - you.card.Height - 10})
	you.AddAbility(LookAroundAbility)
	you.AddAbility(DefendAbility)
	you.AddAbility(AskPinAbility)
	you.AddAbility(GuessPinAbility)
	you.AddAbility(RegularAttackAbility)

	CurrentSave.RecruitedLen = true
	CurrentSave.RecruitedMe = true
	me := &Character{}
	if CurrentSave.RecruitedMe {
		meFnt := &common.Font{
			Size: 64,
			FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
			URL:  "fight/me.ttf",
		}
		meFnt.CreatePreloaded()
		me = AddCharacter(CharacterInfo{
			Name:          "Me",
			CardSprite:    cards.Drawable(1),
			BoxSprite:     boxes.Drawable(1),
			HP:            150,
			MaxHP:         150,
			MP:            75,
			MaxMP:         75,
			Str:           22,
			Def:           30,
			Dex:           30,
			Int:           25,
			Font:          meFnt,
			Clip:          mePlayer,
			CardTextScale: engo.Point{X: 0.25, Y: 0.25},
		}, w)
		you.MoveCard(engo.Point{X: 315 - you.card.Width, Y: 360 - you.card.Height - 10})
		me.MoveCard(engo.Point{X: 340, Y: 360 - me.card.Height - 10})
	}

	len := &Character{}
	if CurrentSave.RecruitedLen {
		lenFnt := &common.Font{
			Size: 64,
			FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
			URL:  "fight/len.ttf",
		}
		lenFnt.CreatePreloaded()
		len = AddCharacter(CharacterInfo{
			Name:          "Len",
			CardSprite:    cards.Drawable(2),
			BoxSprite:     boxes.Drawable(2),
			HP:            85,
			MaxHP:         85,
			MP:            120,
			MaxMP:         120,
			Str:           15,
			Def:           45,
			Dex:           20,
			Int:           48,
			Font:          lenFnt,
			Clip:          lenPlayer,
			CardTextScale: engo.Point{X: 0.35, Y: 0.35},
		}, w)
		if CurrentSave.RecruitedLen && CurrentSave.RecruitedMe {
			you.MoveCard(engo.Point{X: 280 - (3 * you.card.Width / 2), Y: 360 - you.card.Height - 10})
			me.MoveCard(engo.Point{X: 300 - (me.card.Width / 2), Y: 360 - me.card.Height - 10})
			len.MoveCard(engo.Point{X: 320 + (len.card.Width / 2), Y: 360 - len.card.Height - 10})
		} else {
			you.MoveCard(engo.Point{X: 315 - you.card.Width, Y: 360 - you.card.Height - 10})
			len.MoveCard(engo.Point{X: 340, Y: 360 - len.card.Height - 10})
		}
	}

	msgs := []string{
		"A Blood Mouthed Ghost   Appearerated!",
	}
	for _, msg := range msgs {
		engo.Mailbox.Dispatch(CombatLogMessage{
			Msg:  msg,
			Fnt:  selFont,
			Clip: logPlayer,
		})
	}
	engo.Mailbox.Dispatch(PhaseSetMessage{
		Phase: CardSelectPhase,
	})
}
