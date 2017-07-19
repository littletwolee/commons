package commons

import (
	"net/http"
)

const (
	ERROR_Path_EMPTY           = "Your path is empty?"
	ERROR_PARSE                = "Failed parse log level"
	ERROR_PARAMETER_SQL_ATTACK = "Value [%s] is a sensitive word"
	ERROR_PARAMETER_EMPTY      = "Parameter is empty"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

const (
	LEFT  = "left"
	RIGHT = "right"
)

//es
const (
	ES_INDEX_NOT_EXISTS = "index is not exists"
	ES_GTE              = "gte"
	ES_LTE              = "lte"
	ES_ERROR_DATA_EMPTY = "data is empty"
)

const (
	HtpGet      = http.MethodGet
	HttpPost    = http.MethodPost
	HttpPut     = http.MethodPut
	HttpDelete  = http.MethodDelete
	HttpOptions = http.MethodOptions
)

const (
	YEAR = iota
	MONTH
	DAY
	HOUR
	MINUTE
	SECOND
)

const (
	StatusContinue           = 100
	StatusSwitchingProtocols = 101

	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206

	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusTemporaryRedirect = 307

	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418

	StatusInternalServerError     = 500
	StatusNotImplemented          = 501
	StatusBadGateway              = 502
	StatusServiceUnavailable      = 503
	StatusGatewayTimeout          = 504
	StatusHTTPVersionNotSupported = 505

	statusPreconditionRequired          = 428
	statusTooManyRequests               = 429
	statusRequestHeaderFieldsTooLarge   = 431
	statusNetworkAuthenticationRequired = 511
)
