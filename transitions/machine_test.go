package transitions

import (
	"context"
	"fmt"
	"testing"
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

func doAction(ctx context.Context, from State, event Event, to []State) (state State, e error) {
	println(fmt.Sprintf("doAction: [%v] --%s--> %v", ctx.Value("data"), event, to))
	return to[0], nil
}


var Transitions = []Transition{
	{Initial, CreateEvent, []State{WaitPay}, doAction, nil},
	{WaitPay, PayEvent, []State{Paying}, doAction, nil},
	{Paying, PaySuccessEvent, []State{WaitSend}, doAction, nil},
}


func TestStateMachine_Example_Order(t *testing.T) {
	orderStateMachine := New("myFirstStateMachine")
	orderStateMachine.Transitions(Transitions...)

}
