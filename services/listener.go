package services

import (
	"faceit/helpers"
	"fmt"
)

type testListener struct {
}

func (tl *testListener) Listen(eventType EventType, object EventObject) {
	switch eventType {
	case CreateEvent:
		helpers.NotifyChan <- fmt.Sprintf("%s object: %T -> %+v", eventType, object.object, object.object)
		break
	case UpdateEvent:
		helpers.NotifyChan <- fmt.Sprintf("%s object(%d): %T -> original values: %+v -> new values: %+v", eventType, *object.ID, object.object, object.object, object.object2)
		break
	case DeleteEvent:
		helpers.NotifyChan <- fmt.Sprintf("%s object(%d): %T -> %+v", eventType, *object.ID, object.object, object.object)
		break
	case ListEvent:
		helpers.NotifyChan <- fmt.Sprintf("%s objects: %T -> %+v", eventType, object.object, object.object)
		break
	case OtherEvent:
		helpers.NotifyChan <- fmt.Sprintf("%s event of object(s): %T -> %+v", eventType, object.object, object.object)
		break
	default:
		helpers.NotifyChan <- fmt.Sprintf("default case no prepared event for object(s): %T -> %+v", object.object, object.object)
		break
	}
}

func NewTestListener() Listener {
	return &testListener{}
}
