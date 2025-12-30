package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Internal errors
var (
	ErrCacheGetNoRows error = redis.Nil
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
)

// Internal error codes
const (
	CodeBCryptHashingFailed   errCode = "ERR_BCRYPT_HASHING_FAILED"
	CodeCacheQueryExec        errCode = "ERR_CACHE_QUERY_EXECUTION"
	CodeCacheScriptExec       errCode = "ERR_CACHE_SCRIPT_EXECUTION"
	CodeDBQueryExec           errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTX                  errCode = "ERR_DB_TRANSACTION"
	CodeEmailNotAvailable     errCode = "ERR_EMAIL_NOT_AVAILABLE"
	CodeEventPublishingFailed errCode = "ERR_EVENT_PUBLISHING_FAILED"
	CodeInvalidEmail          errCode = "ERR_INVALID_EMAIL"
	CodeInvalidPassword       errCode = "ERR_INVALID_PASSWORD"
	CodeInvalidRequestMeta    errCode = "ERR_INVALID_REQUEST_META"
	CodeJWTGenerationFailed   errCode = "ERR_JWT_GENERATION_FAILED"
	CodeUnknown               errCode = "ERR_UNKNOWN"
)

// External error messages
const (
	MsgInternalServer string = "Internal server error"
)
