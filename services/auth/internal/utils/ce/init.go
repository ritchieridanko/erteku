package ce

import (
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errCode string

type Error struct {
	code    errCode
	message string
	err     error
	fields  []logger.Field
}

func NewError(ec errCode, message string, err error, fields ...logger.Field) *Error {
	return &Error{
		code:    ec,
		message: message,
		err:     err,
		fields:  fields,
	}
}

func (e *Error) Code() errCode {
	return e.code
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}
	return e.message
}

func (e *Error) Fields() []logger.Field {
	return e.fields
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) ToGRPCStatus() error {
	switch e.code {
	case CodeInvalidEmail, CodeInvalidPassword, CodeInvalidRequestMeta:
		return status.Error(codes.InvalidArgument, e.message)
	case CodeAuthNotFound, CodeWrongPassword, CodeWrongSignInMethod:
		return status.Error(codes.Unauthenticated, e.message)
	case CodeSessionNotFound:
		return status.Error(codes.NotFound, e.message)
	case CodeEmailNotAvailable:
		return status.Error(codes.AlreadyExists, e.message)
	case CodeBCryptHashingFailed, CodeCacheQueryExec, CodeCacheScriptExec,
		CodeDBQueryExec, CodeDBTX, CodeEventPublishingFailed, CodeJWTGenerationFailed,
		CodeUnknown:
		return status.Error(codes.Internal, e.message)
	default:
		return status.Error(codes.Internal, e.message)
	}
}
