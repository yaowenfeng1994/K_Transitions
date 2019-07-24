package transitions

import (
	"context"
	"fmt"
	"testing"
)

// 状态
const (
	Initial    = "Initial"
	WaitPay    = "WaitPay"
	PaySuccess = "PaySuccess"
	PayFail    = "PayFail"
)

// 动作
const (
	CreateEvent = "Create"
	PayEvent    = "Pay"
)

const (
	PayCondition = false
)

func CreateEventAction(ctx context.Context, from State, event Event, to []State) (state State, e error) {
	println(fmt.Sprintf("CreateEventAction: [%v] --%s--> %v", ctx.Value("data"), event, to))
	return to[0], nil
}

func PayEventAction(ctx context.Context, from State, event Event, to []State) (state State, e error) {
	println(fmt.Sprintf("PayEventAction: [%v] --%s--> %v", ctx.Value("data"), event, to))
	if PayCondition {
		return PaySuccess, nil
	}
	return PayFail, nil
}

var Transitions = []Transition{
	{Initial, CreateEvent, []State{WaitPay}, CreateEventAction, nil},
	{WaitPay, PayEvent, []State{PaySuccess, PayFail}, PayEventAction, nil},
}

var states = StatesDef{
	Initial:    "开始",
	WaitPay:    "待支付",
	PaySuccess: "支付成功",}

var events = EventsDef{
	CreateEvent: "创建订单",
	PayEvent:    "支付",}

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
	//orderStateMachine.Processor(&OrderEventProcessor{})
	order := context.WithValue(context.TODO(), "data", "order object data")
	//fmt.Println(orderStateMachine.sg.state)
	//_, err := orderStateMachine.Trigger(order, CreateEvent)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(orderStateMachine.sg.state)
	//_, err = orderStateMachine.Trigger(order, CreateEvent)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(orderStateMachine.sg.state)

	orderStateMachine.State(WaitPay)
	_, err := orderStateMachine.Trigger(order, PayEvent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderStateMachine.sg.state)
}
