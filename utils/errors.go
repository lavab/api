package utils

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
)

type Error struct {
	Success  bool        `json:"success"`
	Code     int         `json:"code,omitempty"`
	Location string      `json:"location,omitempty"`
	Error    interface{} `json:"error"`
	Severe   bool        `json:"-"`
}

func NewError(code int, input interface{}, severe bool) *Error {
	_, file, line, _ := runtime.Caller(1)

	return &Error{
		Code:     code,
		Error:    input,
		Location: filepath.Base(file) + ":" + strconv.Itoa(line),
		Severe:   severe,
	}
}

func (e *Error) String() string {
	buf := &bytes.Buffer{}

	if e.Code != 0 {
		buf.WriteString("[")
		buf.WriteString(strconv.Itoa(e.Code))
		buf.WriteString("] ")
	}

	if e.Location != "" {
		buf.WriteString(e.Location)
		buf.WriteString(": ")
	}

	switch v := e.Error.(type) {
	case error:
		buf.WriteString(v.Error())
	case string:
		buf.WriteString(v)
	default:
		buf.WriteString(fmt.Sprintf("%+v", v))
	}

	return buf.String()
}

const (
	AccountsListUnknown = 10000 + iota
)

const (
	AccountsCreateUnknown = 10100 + iota
	AccountsCreateInvalidInput
	AccountsCreateUnknownStep
	AccountsCreateInvalidLength
	AccountsCreateUsernameTaken
	AccountsCreateEmailUsed
	AccountsCreateUnableToInsertAccount
	AccountsCreateUserNotFound
	AccountsCreateInvalidToken
	AccountsCreateInvalidTokenOwner
	AccountsCreateInvalidTokenType
	AccountsCreateExpiredToken
	AccountsCreateAlreadyConfigured
	AccountsCreateWeakPassword
	AccountsCreateUnableToHash
	AccountsCreateUnableToPrepareLabels
	AccountsCreateUnableToCreateAddress
	AccountsCreateUnableToUpdateAccount
)

const (
	AccountsGetUnknown = 10200 + iota
	AccountsGetOnlyMe
	AccountsGetUnableToGet
)

const (
	AccountsUpdateUnknown = 10300 + iota
	AccountsUpdateInvalidInput
	AccountsUpdateOnlyMe
	AccountsUpdateUnableToGet
	AccountsUpdateInvalidCurrentPassword
	AccountsUpdateWeakPassword
	AccountsUpdateUnableToHash
	AccountsUpdateInvalidPublicKey
	AccountsUpdateInvalidPublicKeyOwner
	AccountsUpdateUnableToUpdate
)

const (
	AccountsDeleteUnknown = 10400 + iota
	AccountsDeleteOnlyMe
	AccountsDeleteUnableToGet
	AccountsDeleteUnableToDelete
)

const (
	AccountsWipeDataUnknown = 10500 + iota
	AccountsWipeDataOnlyMe
	AccountsWipeDataUnableToGet
	AccountsWipeDataUnableToDelete
)

const (
	AccountsStartOnboardingUnknown = 10600 + iota
	AccountsStartOnboardingOnlyMe
	AccountsStartOnboardingUnableToGet
	AccountsStartOnboardingMisconfigured
	AccountsStartOnboardingUnableToInit
)
