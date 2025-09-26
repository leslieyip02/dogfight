package message

import "encoding/json"

type Builder struct {
	eventType EventType
	data      interface{}
}

func NewEvent(eventType EventType) *Builder {
	return &Builder{eventType: eventType}
}

func (b *Builder) WithData(data interface{}) *Builder {
	b.data = data
	return b
}

func (b *Builder) Build() ([]byte, error) {
	dataBytes, err := json.Marshal(b.data)
	if err != nil {
		return nil, err
	}

	event := Event{
		Type: b.eventType,
		Data: dataBytes,
	}
	return json.Marshal(event)
}
