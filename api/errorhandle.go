package api

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/logger"
	"net/http"
	"runtime"
)

const (
	// data error
	PermissionDeny     = 1
	JsonUnmarshalError = 2
	TimeoutError       = 3

	// http error
	ParseTokenEorror = 11
	GetBodyError     = 12
	GetTokenError    = 13
	GetHeaderError   = 14
	GetIdError       = 15

	// login error
	AuthenticationError = 31
	LdapConnectError    = 32
	LdapSearchError     = 33
	LdapUserAuthError   = 34
	ParseHeaderError    = 35
	GoogleCodeError     = 36

	// connect other components error
	ConnectAuthenticationError = 41
	ConnectTaskWorkerError     = 42

	// DB connnect error
	SelectDbError = 51
	AddDbError    = 52
	UpdateDbError = 53
	DeleteDbError = 54

	// Redis connect error
	ConnectRedisError = 61
	SelectRedisError  = 62

	// reflect error
	ReflectError           = 71
	UnmarshalResponseError = 72
	SpecialReturnError     = 73

	// task error
	CreateTaskError     = 81
	GetTaskBodyError    = 82
	CannotFindTaskError = 83
	GetTaskStatusError  = 84
)

func ReturnErr(code int, ctx iris.Context, err error, response *Response) {
	var (
		functionName = "unknown_function"
		functionFile = ""
		functionLine = 0
	)
	pc, pcFile, pcLine, ok := runtime.Caller(1)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionFile = pcFile
		functionLine = pcLine
	}
	errMsg := fmt.Sprintf("%s", err.Error())
	logger.Log.WithField("traceId", response.TraceId).WithField("function", functionName).
		WithField("callerFile", functionFile).WithField("callerLine", functionLine).Warnf(errMsg)
	response.Code = code

	if code == PermissionDeny {
		ctx.StatusCode(http.StatusForbidden)
	} else {
		ctx.StatusCode(http.StatusNotImplemented)
	}
	ResponseBody(ctx, response, errMsg)
}
