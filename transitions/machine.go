package transitions

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

type State = string
type Event = string
type StatesDef map[State]string
type EventsDef map[Event]string
type Action func(ctx context.Context, source State, event Event, to []State) (State, error)

// 动作发生时的预处理
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
	name        string // 状态图名称
	states      StatesDef
	events      EventsDef
	state       State                           // 当前状态
	transitions map[State]map[Event]*Transition // 表示图上每个节点的状态，每个状态要执行的动作，每个动作对应的触发器
}

type stateMachine struct {
	//processor EventProcessor  // 暂时没想到这个字段有什么作用，先注释掉
	sg *stateGraph
}

/*
默认实现
*/
type DefaultProcessor struct{}

func (*DefaultProcessor) OnExit(ctx context.Context, state State, event Event) error {
	//log.Printf("exit [%s]", state)
	return nil
}

func (*DefaultProcessor) OnActionFailure(ctx context.Context, from State, event Event, to []State, err error) error {
	//log.Printf("failure %s -(%s)-> [%s]: (%s)", from, event, strings.Join(to, "|"), err.Error())
	return nil
}

func (*DefaultProcessor) OnEnter(ctx context.Context, state State) error {
	//log.Printf("enter [%s]", state)
	return nil
}

var NoopProcessor = &DefaultProcessor{}

func New(name string) *stateMachine {
	return (&stateMachine{sg: &stateGraph{transitions: map[State]map[Event]*Transition{
	}}}).Name(name)
}

func (sm *stateMachine) Name(s string) *stateMachine {
	sm.sg.name = s
	return sm
}

func (sm *stateMachine) States(states StatesDef) *stateMachine {
	sm.sg.states = states
	return sm
}

/*
从持久化数据中拿出状态机数据后重新给状态机赋值
*/
func (sm *stateMachine) State(state State) *stateMachine {
	sm.sg.state = state
	return sm
}

func (sm *stateMachine) Events(events EventsDef) *stateMachine {
	sm.sg.events = events
	return sm
}

//func (sm *stateMachine) Processor(processor EventProcessor) *stateMachine {
//	sm.processor = processor
//	return sm
//}

// 初始化状态机
func (sm *stateMachine) Transitions(transitions ...Transition) *stateMachine {
	for index := range transitions {
		// 为状态机初始化上第一个状态
		if index == 0 {
			sm.sg.state = transitions[index].Source
		}
		newTransfer := &transitions[index]
		events, ok := sm.sg.transitions[newTransfer.Source]
		// 先附上所有节点的状态，有重复节点状态表示该状态有多个触发动作
		if !ok {
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
			events[newTransfer.Event] = newTransfer
		}
	}
	return sm
}

func removeDuplicatesAndEmpty(stateList []State) (ret []State) {
	aLen := len(stateList)
	for i := 0; i < aLen; i++ {
		if (i > 0 && stateList[i-1] == stateList[i]) || len(stateList[i]) == 0 {
			continue
		}
		ret = append(ret, stateList[i])
	}
	return
}

// 触发状态转换
func (sm *stateMachine) Trigger(ctx context.Context, event Event) (State, error) {
	if _, ok := sm.sg.states[sm.sg.state]; !ok {
		return "", errors.New("状态机不包含状态" + sm.sg.state)
	}
	if _, ok := sm.sg.events[event]; !ok {
		return "", errors.New("状态机不包含事件 " + event)
	}
	var stateExist bool
	stateExist = false
	// 保存进来时状态机当前状态
	from := sm.sg.state
	if transfer, ok := sm.sg.transitions[from][event]; ok {
		//processor := sm.processor
		var processor EventProcessor
		// 离开状态处理，转换之前
		if transfer.Processor != nil {
			processor = transfer.Processor
		}
		if processor == nil {
			processor = NoopProcessor
		}
		_ = processor.OnExit(ctx, from, event)
		to, err := transfer.Action(ctx, from, event, transfer.To)
		if err != nil {
			// 发生错误时 转换执行错误处理 将状态机当前状态置为来源状态
			_ = processor.OnActionFailure(ctx, from, event, transfer.To, err)
			sm.sg.state = from
			return to, err
		}
		// 校验经过动作后的目标状态是否与定义的一致
		for _, toState := range transfer.To {
			if toState == to {
				stateExist = true
				break
			}
		}
		if !stateExist {
			_ = processor.OnActionFailure(ctx, from, event, transfer.To, err)
			sm.sg.state = from
			return to, errors.New(fmt.Sprintf("返回的状态 %v 不存在目标状态 %v 中", to, transfer.To))
		}
		sm.sg.state = to
		// 进入状态处理，转换之后
		_ = processor.OnEnter(ctx, to)
		return to, err
	}
	return "", errors.New(fmt.Sprintf("该状态下没有定义状态转换事件 %v", event))
}
