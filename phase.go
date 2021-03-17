package main

import (
	"sync"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

type Phase uint8

const (
	BeginingPhase Phase = iota
	ListenPhase
	WalkPhase
	LogClearPhase
	AcceptPhase
)

var PhaseSetMessageType = "Phase Set Message"

type PhaseSetMessage struct {
	Phase
}

func (PhaseSetMessage) Type() string { return PhaseSetMessageType }

var AcceptSetMessageType = "Accept Set Message"

type AcceptSetMessage struct {
	AcceptFunc func()
}

func (AcceptSetMessage) Type() string { return AcceptSetMessageType }

var PhaseDequeuMessageType = "Accept Dequeue Message"

type PhaseDequeuMessage struct{}

func (PhaseDequeuMessage) Type() string { return PhaseDequeuMessageType }

type PhaseSystem struct {
	entities []ecs.Identifier

	move   *MoveSystem
	cursor *CursorSystem

	currentPhase, setPhase Phase
	queue                  []Phase
	lock                   sync.Mutex

	acceptFunc    func()
	acceptLogWait bool
}

func (s *PhaseSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *MoveSystem:
			s.move = sys
		case *CursorSystem:
			s.cursor = sys
		}
	}

	engo.Mailbox.Listen(PhaseSetMessageType, func(message engo.Message) {
		msg, ok := message.(PhaseSetMessage)
		if !ok {
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		if s.currentPhase == BeginingPhase && s.setPhase == BeginingPhase {
			s.setPhase = msg.Phase
		} else {
			s.queue = append(s.queue, msg.Phase)
		}
	})

	engo.Mailbox.Listen(AcceptSetMessageType, func(message engo.Message) {
		msg, ok := message.(AcceptSetMessage)
		if !ok {
			return
		}
		s.acceptFunc = msg.AcceptFunc
	})

	engo.Mailbox.Listen(PhaseDequeuMessageType, func(message engo.Message) {
		_, ok := message.(PhaseDequeuMessage)
		if !ok {
			return
		}
		s.dequeue()
	})
}

func (s *PhaseSystem) Remove(basic ecs.BasicEntity) {}

func (s *PhaseSystem) AddByInterface(i ecs.Identifier) {
	s.entities = append(s.entities, i)
}

func (s *PhaseSystem) Update(dt float32) {
	if s.currentPhase != s.setPhase {
		for _, entity := range s.entities {
			if mover, ok := entity.(Moveable); ok {
				s.move.Remove(*mover.GetBasicEntity())
			}
			if cur, ok := entity.(CursorAble); ok {
				s.cursor.Remove(*cur.GetBasicEntity())
			}
		}
		engo.Mailbox.Dispatch(CombatLogPauseMessage{
			Pause: true,
		})
		engo.Mailbox.Dispatch(DoorSystemPauseMessage{
			Pause: true,
		})
		engo.Mailbox.Dispatch(InterestSystemPauseMessage{
			Pause: true,
		})
		engo.Mailbox.Dispatch(AcceptSystemPauseMessage{
			Pause: true,
		})
		switch s.setPhase {
		case ListenPhase:
			engo.Mailbox.Dispatch(CombatLogPauseMessage{
				Pause: false,
			})
		case WalkPhase:
			for _, entity := range s.entities {
				if mover, ok := entity.(Moveable); ok {
					s.move.AddByInterface(mover)
				}
			}
			engo.Mailbox.Dispatch(DoorSystemPauseMessage{
				Pause: false,
			})
			engo.Mailbox.Dispatch(InterestSystemPauseMessage{
				Pause: false,
			})
		case LogClearPhase:
			engo.Mailbox.Dispatch(CombatLogClearMessage{})
		case AcceptPhase:
			engo.Mailbox.Dispatch(CombatLogPauseMessage{
				Pause: false,
			})
			s.acceptLogWait = true
		}
		s.currentPhase = s.setPhase
	}

	switch s.currentPhase {
	case ListenPhase:
		if !s.logDone() {
			return
		}
		if engo.Input.Button("A").JustPressed() || engo.Input.Button("B").JustPressed() ||
			engo.Input.Button("X").JustPressed() || (engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft) {
			s.dequeue()
		}
	case WalkPhase:
		//idk if this does anything honestly
	case LogClearPhase:
		s.dequeue()
	case AcceptPhase:
		if !s.logDone() {
			return
		} else {
			if s.acceptLogWait {
				engo.Mailbox.Dispatch(AcceptSystemPauseMessage{
					Pause: false,
				})
				s.acceptLogWait = false
			}
		}
		if engo.Input.Button("A").JustPressed() || (engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft) {
			accept := &YesSelectedMessage{}
			engo.Mailbox.Dispatch(accept)
			if accept.Selected && s.acceptFunc != nil {
				s.acceptFunc()
				s.dequeue()
			} else {
				s.dequeue()
			}
		}
		if engo.Input.Button("B").JustPressed() {
			s.dequeue()
		}
	}
}

func (s *PhaseSystem) dequeue() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.queue) > 0 {
		s.setPhase = s.queue[0]
		s.queue = s.queue[1:]
	} else {
		s.setPhase = BeginingPhase
	}
}

func (s *PhaseSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
