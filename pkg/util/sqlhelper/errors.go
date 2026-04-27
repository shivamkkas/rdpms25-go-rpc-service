package sqlhelper

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorCode struct {
	Code int
	Msg  string
}

func ServerErrorWithMsg(msg string) *ErrorCode {
	return &ErrorCode{Code: http.StatusInternalServerError, Msg: msg}
}

func ServerError() *ErrorCode {
	return &ErrorCode{Code: http.StatusInternalServerError, Msg: "Server side error occurred!"}
}

func SetDbErrorGinContext(c *gin.Context, e error) {
	err := DbError(e)
	c.JSON(err.Code, err.Msg)
}

func DbError(e error) *ErrorCode {
	if e == nil {
		return nil
	}

	// if gin.Mode() != gin.ReleaseMode {
	// return &ErrorCode{Code: http.StatusInternalServerError, Msg: e.Error()}
	// }

	msg := e.Error()

	violatesUnique := strings.Contains(msg, "violates unique constraint")
	violatesForeign := strings.Contains(msg, "violates foreign key constraint")

	operationType := ""
	for _, t := range []string{"delete", "insert", "update"} {
		if strings.Contains(msg, fmt.Sprintf("unable to %s", t)) {
			operationType = t
		}
	}

	composedErrMsg := ""
	if violatesUnique && len(operationType) > 0 {
		composedErrMsg = fmt.Sprintf("Entity already exists! Unable to %s", operationType)
	} else if violatesForeign && len(operationType) > 0 {
		composedErrMsg = fmt.Sprintf("Related entity did not matched! Unable to %s", operationType)
	}

	if len(composedErrMsg) > 0 {
		return &ErrorCode{Code: http.StatusConflict, Msg: composedErrMsg}
	} else if strings.Contains(msg, "pq") {
		return &ErrorCode{Code: http.StatusInternalServerError, Msg: "Something went wrong in database"}
	} else {
		slog.Error("unhandled backend error", "error", msg)
		return &ErrorCode{Code: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}
}

func InvalidDataError(msg string) *ErrorCode {
	return &ErrorCode{Code: http.StatusBadRequest, Msg: msg}
}

func RequiredFieldError(key string) *ErrorCode {
	return &ErrorCode{Code: http.StatusBadRequest, Msg: key + " is required"}
}
