package transitions

import (
	"context"
	"fmt"
	"sort"
)

type State = string
type Event = string

type StatesDef map[State]string

//type EventsDef map[Event]string
type Action func(ctx context.Context, source State, event Event, to []State) (State, error)
type EventProcessor interface {
	OnExit(ctx context.Context, state State, event Event) error
	OnActionFailure(ctx context.Context, from State, event Event, to []State, err error) error
	OnEnter(ctx context.Context, state State) error
}

type Transition struct {
	Source    State   // 来源状态
	Event     Event   // 发生状态迁移的动作的名字
	To        []State // 目标状态
	Action    Action  // 需要自己实现该动作
	Processor EventProcessor
}

/*
状态机执行表述图
*/

type stateGraph struct {
	name string // 状态图名称
	//start  []State
	//end    []State
	//states StatesDef
	//events EventsDef
	state       State                           // 当前状态
	transitions map[State]map[Event]*Transition // 表示图上每个节点的状态，每个状态要执行的动作，每个动作对应的触发器
}

type stateMachine struct {
	processor EventProcessor
	sg        *stateGraph
}

func New(name string) *stateMachine {
	return (&stateMachine{sg: &stateGraph{transitions: map[State]map[Event]*Transition{
	}}}).Name(name)
}

func (sm *stateMachine) Name(s string) *stateMachine {
	sm.sg.name = s
	return sm
}

//func (sm *stateMachine) States(states StatesDef) *stateMachine {
//	sm.sg.states = states
//	return sm
//}

func (sm *stateMachine) Transitions(transitions ...Transition) *stateMachine {
	for index := range transitions {
		newTransfer := &transitions[index]
		events, ok := sm.sg.transitions[newTransfer.Source]
		// 先附上所有节点的状态
		if !ok {
			fmt.Println(events)
			events = map[Event]*Transition{}
			sm.sg.transitions[newTransfer.Source] = events
		}
		// 再附上每个状态要执行的动作
		if transfer, ok := events[newTransfer.Event]; ok {
			transfer.To = append(transfer.To, newTransfer.To...)
			// 去掉重复的状态
			sort.Strings(transfer.To)
			transfer.To = removeDuplicatesAndEmpty(transfer.To)
			events[newTransfer.Event] = transfer
		} else {
			fmt.Println(newTransfer)
			events[newTransfer.Event] = newTransfer
		}
	}
	return sm
}

func removeDuplicatesAndEmpty(a []State) (ret []State) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}