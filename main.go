package main

import (
	"./transitions"
	"context"
	"fmt"
)

// 状态
const (
	Initial  = "Initial"
	WaitPay  = "WaitPay"
	Paying   = "Paying"
	WaitSend = "WaitSend"
)

// 动作
const (
	CreateEvent     = "Create"
	PayEvent        = "Pay"
	PaySuccessEvent = "PaySuccess"
)


func doAction(ctx context.Context, from transitions.State, event transitions.Event, to []transitions.State) (state transitions.State, e error) {
	println(fmt.Sprintf("doAction: [%v] --%s--> %v", ctx.Value("data"), event, to))
	return to[0], nil
}

func main() {
	var Transitions = []transitions.Transition{
		{Initial, CreateEvent, []transitions.State{WaitPay}, doAction, nil},
		{WaitPay, PayEvent, []transitions.State{Paying}, doAction, nil},
		{Paying, PaySuccessEvent, []transitions.State{WaitSend}, doAction, nil},
	}

	orderStateMachine := transitions.New("myFirstStateMachine")
	orderStateMachine.Transitions(Transitions...)
}
