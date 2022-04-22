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
		"fight/log.png",
		"fight/dots.png",
		"fight/log.ttf",
		"fight/bg.ogg",
		"fight/log.wav",
		"fight/bg.png",
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
	// w.AddSystem(&systems.ExitSystem{})

	var phaseable *common.BasicFace
	w.AddSystemInterface(&PhaseSystem{}, phaseable, nil)

	selFont := &common.Font{
		Size: 72,
		FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
		URL:  "fight/log.ttf",
	}
	selFont.CreatePreloaded()

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

	bg := sprite{BasicEntity: ecs.NewBasic()}
	tex0, _ := common.LoadedSprite("fight/bg.png")
	bg.Drawable = pixelshader.PixelRegion{
		Tex0: tex0,
	}
	bg.SetShader(fightShader)
	bg.SetZIndex(0)
	w.AddEntity(&bg)

	cards := common.NewSpriteSheetWithBorderFromFile("fight/cards.png", 0, 0, 1, 1)
	cardNum := 1
	if CurrentSave.RecruitedMe {
		cardNum++
	}
	if CurrentSave.RecruitedLen {
		cardNum++
	}
	youFnt := &common.Font{
		Size: 64,
		FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
		URL:  "fight/you.ttf",
	}
	youFnt.CreatePreloaded()
	you := AddPlayer(&Character{
		Name:   "You",
		Sprite: cards.Drawable(0),
		HP:     100,
		MP:     100,
		Str:    25,
		Def:    25,
		Dex:    35,
		Int:    40,
		Abilities: []Ability{
			Ability{},
		},
		Font: youFnt,
	}, cardNum)
	me := &Player{}
	if CurrentSave.RecrutedMe {
		meFnt := &common.Font{
			Size: 64,
			FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
			URL:  "fight/me.ttf",
		}
		meFnt.CreatePreloaded()
		// me = AddPlayer(&Character{
		//
		// })
	}
	len := &Player{}
	if CurrentSave.RecrutedLen {
		lenFnt := &common.Font{
			Size: 64,
			FG:   color.RGBA{R: 0xdc, G: 0xd2, B: 0xd2, A: 0xff},
			URL:  "fight/len.ttf",
		}
		lenFnt.CreatePreloaded()
		// len = AddPlayer(&Character{
		//
		// })
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
		Phase: ListenPhase,
	})
}
