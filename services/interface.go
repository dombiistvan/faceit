package services

import "github.com/labstack/echo/v4"

const (
	CreateEvent EventType = "create"
	UpdateEvent EventType = "update"
	DeleteEvent EventType = "delete"
	ListEvent   EventType = "list"
	OtherEvent  EventType = "other"
)

type EventType string
type EventObject struct {
	object  any
	object2 any
	ID      *uint
}

type Listener interface {
	Listen(eventType EventType, object EventObject)
}

type Service interface {
	HealthCheck(c echo.Context) error
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	ListUsers(c echo.Context) error
	Setup() error
	Broadcast(eventType EventType, object EventObject)
}
