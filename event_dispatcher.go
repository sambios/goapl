package goapl


type EventHandlerFunc func()

// 事件调度器中存放的单元
type EventSaver struct {
	evtid uint32
	Listeners map[*EventHandlerFunc]bool
}

// 事件调度器基类
type EventDispatcher struct {
	events map[uint32]EventSaver
}

// 事件调度接口
type IEventDispatcher interface {

	AddEventListener(eventType uint32, fn EventHandlerFunc)

	RemoveEventListener(eventType uint32, fn EventHandlerFunc) bool

	HasEventListener(eventType uint32) bool

	EventTrigger(eventType uint32) bool
}

// 创建事件派发器
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{events:make(map[uint32]EventSaver)}
}


// 事件调度器添加事件
func (this *EventDispatcher) AddEventListener(eventType uint32, handlerFunc EventHandlerFunc) {
	if evt, ok := this.events[eventType]; ok {
		evt.Listeners[&handlerFunc] = true
		return
	}

	evt:= EventSaver{evtid:eventType, Listeners:make(map[*EventHandlerFunc]bool)}
	this.events[eventType] = evt
	evt.Listeners[&handlerFunc] = true

}

// 事件调度器移除某个监听
func (this *EventDispatcher) RemoveEventListener(eventType uint32, listener EventHandlerFunc) bool {
	if evt, ok := this.events[eventType]; ok {
		delete(evt.Listeners, &listener)
		return true
	}

	return false
}

// 事件调度器是否包含某个类型的监听
func (this *EventDispatcher) HasEventListener(eventType uint32) bool {
	if _, ok := this.events[eventType]; ok {
		return true
	}
	return false
}

// 事件调度器派发事件
func (this *EventDispatcher) EventTrigger(eventType uint32) bool {
	for _, evt := range this.events {
		if evt.evtid == eventType {
			for listener := range evt.Listeners {
				(*listener)()
			}
			return true
		}
	}
	return false
}
