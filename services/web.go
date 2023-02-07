package services

import (
	"encoding/json"
	"faceit/db"
	"faceit/db/models"
	"faceit/helpers"
	"faceit/requests"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type webService struct {
	dbm      *db.Manager
	e        *echo.Echo
	listener Listener
}

// create new web service interface
func NewWebService(e *echo.Echo, dbm *db.Manager, listener Listener) (Service, error) {
	return &webService{
		e:        e,
		dbm:      dbm,
		listener: listener,
	}, nil
}

// in case of multiple listeners on event types it would broadcast the events to each listeners
func (ws webService) Broadcast(eventType EventType, object EventObject) {
	if ws.listener == nil {
		return
	}

	ws.listener.Listen(eventType, object)
}

// HealthCheck function
func (ws webService) HealthCheck(c echo.Context) error {
	defer func() {
		if f := recover(); f != nil {
			helpers.ErrChan <- fmt.Errorf("panic occured: %+v", f)
		}
	}()

	ws.Broadcast(OtherEvent, EventObject{object: "healthcheck", ID: nil})

	return c.String(http.StatusOK, "it's alive!")
}

// CreateUser path of user creation
func (ws webService) CreateUser(c echo.Context) error {
	defer func() {
		if f := recover(); f != nil {
			helpers.ErrChan <- fmt.Errorf("panic occured: %+v", f)
		}
	}()

	var createUserReq requests.CreateUserRequest
	err := c.Bind(&createUserReq)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, requests.ErrBindRequest(err))
	}

	if err = createUserReq.Validate(); err != nil {
		return ws.ValidationError(c, err)
	}

	_, err = ws.dbm.GetUserToEmail(createUserReq.Email, nil)
	if err == nil {
		return echo.NewHTTPError(http.StatusConflict, requests.ErrAlreadyExist)
	}

	model, err := models.CreateUserFromReq(createUserReq)
	if err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	ws.Broadcast(CreateEvent, EventObject{object: model})

	if err = ws.dbm.Save(&model, nil); err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	return c.JSON(http.StatusCreated, model)
}

// UpdateUser path of user update
func (ws webService) UpdateUser(c echo.Context) error {
	defer func() {
		if f := recover(); f != nil {
			helpers.ErrChan <- fmt.Errorf("panic occured: %+v", f)
		}
	}()

	var req requests.UpdateUserRequest
	err := c.Bind(&req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, requests.ErrBindRequest(err))
	}

	if err = req.Validate(); err != nil {
		fmt.Println("validate", err)
		return ws.ValidationError(c, err)
	}

	userUUID, err := uuid.Parse(c.Param("UUID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user UUID")
	}

	user, err := ws.dbm.GetUserToUUID(userUUID.String(), nil)
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, requests.ErrNotFound("user", err))
	}

	if err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	if req.Email != nil {
		existingEmailUser, err := ws.dbm.GetUserToEmail(*req.Email, nil)
		if err == nil && existingEmailUser.ID != user.ID {
			return echo.NewHTTPError(http.StatusConflict, requests.ErrAlreadyExist)
		}
	}

	model, err := models.CreateUpdateUserFromReq(user, req)
	if err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	ws.Broadcast(UpdateEvent, EventObject{object: user, object2: model, ID: &user.ID})

	if err = ws.dbm.Save(&model, nil); err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	return c.JSON(http.StatusOK, model)
}

// DeleteUser path of user deletion
func (ws webService) DeleteUser(c echo.Context) error {
	defer func() {
		if f := recover(); f != nil {
			helpers.ErrChan <- fmt.Errorf("panic occured: %+v", f)
		}
	}()

	userUUID, err := uuid.Parse(c.Param("UUID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user UUID")
	}

	user, err := ws.dbm.GetUserToUUID(userUUID.String(), nil)
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, requests.ErrNotFound("user", err))
	}

	if err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	ws.Broadcast(DeleteEvent, EventObject{object: user, ID: &user.ID})

	if err = ws.dbm.Delete(&user, nil); err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	return c.NoContent(http.StatusOK)
}

// ListUsers listing users based on query parameters
func (ws webService) ListUsers(c echo.Context) error {
	defer func() {
		if f := recover(); f != nil {
			helpers.ErrChan <- fmt.Errorf("panic occured: %+v", f)
		}
	}()

	var listUserReq requests.ListUsersRequest
	err := c.Bind(&listUserReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, requests.ErrBindRequest(err))
	}

	if err = listUserReq.Validate(); err != nil {
		return ws.ValidationError(c, err)
	}

	users, err := ws.dbm.ListUsersTo(listUserReq, nil)
	if err != nil {
		helpers.ErrChan <- err
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrInternal)
	}

	ws.Broadcast(ListEvent, EventObject{object: listUserReq})

	if len(users) == 0 {
		return echo.NewHTTPError(http.StatusNoContent, requests.ErrNotFoundSimple)
	}

	return c.JSON(http.StatusOK, users)
}

// ValidationError prepare validation error of structs to response
func (ws webService) ValidationError(c echo.Context, err error) error {
	var validationErrors = requests.ValidationErrorMap{}

	errBytes, err := json.Marshal(err)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrMarshall(err))
	}
	if err = json.Unmarshal(errBytes, &validationErrors); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, requests.ErrUnmarshall("validation object", err))
	}

	return c.JSON(http.StatusBadRequest, validationErrors)
}

// set upweb service routes
func (ws webService) Setup() error {
	ws.e.GET("/test", ws.HealthCheck)

	userGroup := ws.e.Group("user")
	userGroup.POST("", ws.CreateUser, helpers.JsonMiddleware)
	userGroup.PUT("/:UUID", ws.UpdateUser, helpers.JsonMiddleware)
	userGroup.DELETE("/:UUID", ws.DeleteUser)
	userGroup.GET("/list", ws.ListUsers)

	return nil
}
