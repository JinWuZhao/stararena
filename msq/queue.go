package msq

import "github.com/jinwuzhao/stararena/command"

type CommandQueue struct {
	ch chan command.Command
}

func NewCommandQueue(cap int) *CommandQueue {
	return &CommandQueue{
		ch: make(chan command.Command, cap),
	}
}

func (m *CommandQueue) Push(cmd command.Command) bool {
	select {
	case m.ch <- cmd:
		return true
	default:
		return false
	}
}

func (m *CommandQueue) Pop() (command.Command, bool) {
	select {
	case cmd := <-m.ch:
		return cmd, true
	default:
		return nil, false
	}
}

func (m *CommandQueue) PushChan() chan<- command.Command {
	return m.ch
}

func (m *CommandQueue) PopChan() <-chan command.Command {
	return m.ch
}
