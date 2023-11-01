package ws

type MessageKind int

const (
	CHAT_MESSAGE MessageKind = iota
	CHAT_CREATED
	USER_OFFLINE
)

type MessageActionType int

const (
	Send MessageActionType = iota
	Delete
	Edit
)

type Message struct {
	ActionType MessageActionType `json:"actionType"`
	Kind       MessageKind       `json:"kind"`
	Payload    interface{}       `json:"payload"`
	strategy   BroadcastStrategy
}

func (m *Message) Strategy() BroadcastStrategy {
	return m.strategy
}

func (m *Message) SetStrategy(strategy BroadcastStrategy) {
	m.strategy = strategy
}
