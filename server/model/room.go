package model

import (
	"errors"
	"fmt"
	"regexp"

	json "github.com/json-iterator/go"
	"github.com/synctv-org/synctv/internal/conf"

	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyRoomId          = errors.New("empty room id")
	ErrRoomIdTooLong        = errors.New("room id too long")
	ErrRoomIdHasInvalidChar = errors.New("room id has invalid char")

	ErrPasswordTooLong        = errors.New("password too long")
	ErrPasswordHasInvalidChar = errors.New("password has invalid char")

	ErrEmptyUsername          = errors.New("empty username")
	ErrUsernameTooLong        = errors.New("username too long")
	ErrUsernameHasInvalidChar = errors.New("username has invalid char")
)

var (
	alphaNumReg        = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	alphaNumChineseReg = regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_\-]+$`)
)

type FormatEmptyPasswordError string

func (f FormatEmptyPasswordError) Error() string {
	return fmt.Sprintf("%s password empty", string(f))
}

type CreateRoomReq struct {
	RoomId       string `json:"roomId"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	UserPassword string `json:"userPassword"`
	Hidden       bool   `json:"hidden"`
}

func (c *CreateRoomReq) Decode(ctx *gin.Context) error {
	return json.NewDecoder(ctx.Request.Body).Decode(c)
}

func (c *CreateRoomReq) Validate() error {
	if c.RoomId == "" {
		return ErrEmptyRoomId
	} else if len(c.RoomId) > 32 {
		return ErrRoomIdTooLong
	} else if !alphaNumChineseReg.MatchString(c.RoomId) {
		return ErrRoomIdHasInvalidChar
	}

	if c.Password != "" {
		if len(c.Password) > 32 {
			return ErrPasswordTooLong
		} else if !alphaNumReg.MatchString(c.Password) {
			return ErrPasswordHasInvalidChar
		}
	} else if conf.Conf.Room.MustPassword {
		return FormatEmptyPasswordError("room")
	}

	if c.Username == "" {
		return ErrEmptyUsername
	} else if len(c.Username) > 32 {
		return ErrUsernameTooLong
	} else if !alphaNumChineseReg.MatchString(c.Username) {
		return ErrUsernameHasInvalidChar
	}

	if c.UserPassword == "" {
		return FormatEmptyPasswordError("user")
	} else if len(c.UserPassword) > 32 {
		return ErrPasswordTooLong
	} else if !alphaNumReg.MatchString(c.UserPassword) {
		return ErrPasswordHasInvalidChar
	}

	return nil
}

type RoomListResp struct {
	RoomId       string `json:"roomId"`
	PeopleNum    int64  `json:"peopleNum"`
	NeedPassword bool   `json:"needPassword"`
	Creator      string `json:"creator"`
	CreatedAt    int64  `json:"createdAt"`
}

type LoginRoomReq struct {
	RoomId       string `json:"roomId"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	UserPassword string `json:"userPassword"`
}

func (l *LoginRoomReq) Decode(ctx *gin.Context) error {
	return json.NewDecoder(ctx.Request.Body).Decode(l)
}

func (l *LoginRoomReq) Validate() error {
	if l.RoomId == "" {
		return ErrEmptyRoomId
	}

	if l.Username == "" {
		return ErrEmptyUsername
	}

	if l.UserPassword == "" {
		return FormatEmptyPasswordError("user")
	}

	return nil
}

type SetRoomPasswordReq struct {
	Password string `json:"password"`
}

func (s *SetRoomPasswordReq) Decode(ctx *gin.Context) error {
	return json.NewDecoder(ctx.Request.Body).Decode(s)
}

func (s *SetRoomPasswordReq) Validate() error {
	if len(s.Password) > 32 {
		return ErrPasswordTooLong
	} else if !alphaNumReg.MatchString(s.Password) {
		return ErrPasswordHasInvalidChar
	}
	return nil
}

type UsernameReq struct {
	Username string `json:"username"`
}

func (u *UsernameReq) Decode(ctx *gin.Context) error {
	return json.NewDecoder(ctx.Request.Body).Decode(u)
}

func (u *UsernameReq) Validate() error {
	if u.Username == "" {
		return ErrEmptyUsername
	} else if len(u.Username) > 32 {
		return ErrUsernameTooLong
	} else if !alphaNumChineseReg.MatchString(u.Username) {
		return ErrUsernameHasInvalidChar
	}
	return nil
}
