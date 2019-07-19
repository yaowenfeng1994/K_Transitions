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
	PaySuccess   = "PaySuccess"
)

// 动作
const (
	CreateEvent     = "Create"
	PayEvent        = "Pay"
)

func doAction(ctx context.Context, from State, event Event, to []State) (state State, e error) {
	println(fmt.Sprintf("doAction: [%v] --%s--> %v", ctx.Value("data"), event, to))
	return to[0], nil
}

var Transitions = []Transition{
	{Initial, CreateEvent, []State{WaitPay}, doAction, nil},
	{WaitPay, PayEvent, []State{PaySuccess}, doAction, nil},
}

var states = StatesDef{
	Initial: "开始",
	WaitPay: "待支付",
	PaySuccess:  "支付成功",}

var events = EventsDef{
	CreateEvent:     "创建订单",
	PayEvent:        "支付",}

type OrderEventProcessor struct{}

func (*OrderEventProcessor) OnExit(ctx context.Context, state State, event Event) error {
	println(fmt.Sprintf("OnExit: [%v] Exit [%v] on event [%v]", ctx.Value("data"), state, event))
	return nil
}

func (*OrderEventProcessor) OnActionFailure(ctx context.Context, from State, event Event, to []State, err error) error {
	println(fmt.Sprintf("OnActionFailure: [%v] do action error %v --%v--> %v", ctx.Value("data"), from, event, to))
	return nil
}

func (*OrderEventProcessor) OnEnter(ctx context.Context, state State) error {
	println(fmt.Sprintf("OnEnter: [%v] Enter [%v]", ctx.Value("data"), state))
	return nil
}

func TestStateMachine_Example_Order(t *testing.T) {
	orderStateMachine := New("myFirstStateMachine")
	orderStateMachine.Transitions(Transitions...)
	orderStateMachine.States(states)
	orderStateMachine.Events(events)
	orderStateMachine.Processor(&OrderEventProcessor{})
	//order := context.WithValue(context.TODO(), "data", "order object data")

	//state, err := orderStateMachine.Trigger(order, Initial, CreateEvent)
	//println(fmt.Sprintf("====: %v : %v", state, err))

	//state, err = orderStateMachine.Trigger(order, Initial, PayEvent)
	//println(fmt.Sprintf("====: %v : %v", state, err))

	//state, err = orderStateMachine.Trigger(order, WaitPay, PayEvent)
	//println(fmt.Sprintf("====: %v : %v", state, err))
}
