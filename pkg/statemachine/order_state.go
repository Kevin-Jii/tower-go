package statemachine

import (
	"errors"
	"fmt"
)

// State 状态类型
type State int8

// Action 动作类型
type Action string

// 采购单状态
const (
	StatePending   State = 1 // 待确认
	StateConfirmed State = 2 // 已确认
	StateCompleted State = 3 // 已完成
	StateCancelled State = 4 // 已取消
)

// 采购单动作
const (
	ActionConfirm  Action = "confirm"
	ActionComplete Action = "complete"
	ActionCancel   Action = "cancel"
)

// Transition 状态转换
type Transition struct {
	From   State
	To     State
	Action Action
}

// StateMachine 状态机
type StateMachine struct {
	transitions []Transition
	hooks       map[Action][]func(from, to State) error
}

// NewOrderStateMachine 创建采购单状态机
func NewOrderStateMachine() *StateMachine {
	sm := &StateMachine{
		transitions: []Transition{
			{From: StatePending, To: StateConfirmed, Action: ActionConfirm},
			{From: StatePending, To: StateCancelled, Action: ActionCancel},
			{From: StateConfirmed, To: StateCompleted, Action: ActionComplete},
			{From: StateConfirmed, To: StateCancelled, Action: ActionCancel},
		},
		hooks: make(map[Action][]func(from, to State) error),
	}
	return sm
}

// CanTransition 检查是否可以转换
func (sm *StateMachine) CanTransition(from State, action Action) bool {
	for _, t := range sm.transitions {
		if t.From == from && t.Action == action {
			return true
		}
	}
	return false
}

// GetNextState 获取下一个状态
func (sm *StateMachine) GetNextState(from State, action Action) (State, error) {
	for _, t := range sm.transitions {
		if t.From == from && t.Action == action {
			return t.To, nil
		}
	}
	return 0, fmt.Errorf("invalid transition: from %d with action %s", from, action)
}

// OnAction 注册动作钩子
func (sm *StateMachine) OnAction(action Action, hook func(from, to State) error) {
	sm.hooks[action] = append(sm.hooks[action], hook)
}

// Execute 执行状态转换
func (sm *StateMachine) Execute(from State, action Action) (State, error) {
	to, err := sm.GetNextState(from, action)
	if err != nil {
		return 0, err
	}

	// 执行钩子
	for _, hook := range sm.hooks[action] {
		if err := hook(from, to); err != nil {
			return 0, err
		}
	}

	return to, nil
}

// ValidateTransition 验证状态转换（兼容旧代码）
func ValidateTransition(currentStatus, newStatus int8) bool {
	validTransitions := map[State][]State{
		StatePending:   {StateConfirmed, StateCancelled},
		StateConfirmed: {StateCompleted, StateCancelled},
		StateCompleted: {},
		StateCancelled: {},
	}

	allowedStatuses, ok := validTransitions[State(currentStatus)]
	if !ok {
		return false
	}

	for _, allowed := range allowedStatuses {
		if State(newStatus) == allowed {
			return true
		}
	}
	return false
}

// GetStateName 获取状态名称
func GetStateName(status State) string {
	names := map[State]string{
		StatePending:   "待确认",
		StateConfirmed: "已确认",
		StateCompleted: "已完成",
		StateCancelled: "已取消",
	}
	if name, ok := names[status]; ok {
		return name
	}
	return "未知"
}

// GetAvailableActions 获取当前状态可用的动作
func (sm *StateMachine) GetAvailableActions(from State) []Action {
	var actions []Action
	for _, t := range sm.transitions {
		if t.From == from {
			actions = append(actions, t.Action)
		}
	}
	return actions
}

var ErrInvalidTransition = errors.New("invalid state transition")
