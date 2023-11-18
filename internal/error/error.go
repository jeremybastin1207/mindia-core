package mindiaerr

type ErrCode int

const (
	ErrCodeUnknown       ErrCode = 0
	ErrCodeMediaNotFound ErrCode = iota + 1
	ErrCodeApiKeyNotFound
	ErrBadRequest
	ErrCodeInternal
	ErrCodeMimeTypeNotSupported
	ErrCodeNamedTransformationNotFound
	ErrCodeTransformationNotFound
	ErrCodeUnauthorizedRequest
	ErrCodeServiceUnavailable
)

func (e ErrCode) Code() string {
	switch e {
	case ErrCodeMediaNotFound:
		return "err_media_not_found"
	case ErrCodeApiKeyNotFound:
		return "err_apikey_not_found"
	case ErrCodeMimeTypeNotSupported:
		return "err_mimetype_not_supported"
	case ErrCodeNamedTransformationNotFound:
		return "err_named_transformation_not_found"
	case ErrCodeTransformationNotFound:
		return "err_transformation_not_found"
	case ErrCodeUnauthorizedRequest:
		return "err_unauthorized_request"
	case ErrCodeServiceUnavailable:
		return "err_service_unavailable"
	case ErrBadRequest:
		return "err_bad_request"
	case ErrCodeInternal:
		return "err_internal"
	default:
		return "unknown error"
	}
}

func (e ErrCode) RawMessage() string {
	switch e {
	case ErrCodeMediaNotFound:
		return "unable to find the media"
	case ErrCodeApiKeyNotFound:
		return "unable to find the api key"
	case ErrCodeMimeTypeNotSupported:
		return "mime type not supported"
	case ErrCodeNamedTransformationNotFound:
		return "unable to find the named transformation"
	case ErrCodeTransformationNotFound:
		return "unable to find the transformation"
	case ErrCodeUnauthorizedRequest:
		return "unauthroized request"
	case ErrCodeServiceUnavailable:
		return "service (temporarely) unavailable"
	case ErrBadRequest:
		return "bad request"
	case ErrCodeInternal:
		return "internal error"
	default:
		return "unknown error"
	}
}

type Error struct {
	ErrCode ErrCode
	Msg     error
}

func New(errCode ErrCode) *Error {
	return &Error{
		ErrCode: errCode,
	}
}

func (e *Error) Error() string {
	if e.Msg != nil {
		return e.Msg.Error()
	}
	return e.ErrCode.RawMessage()
}

type ApiError struct {
	RawMessage string `json:"message"`
	Code       string `json:"code"`
}

func NewApiError(err Error) *ApiError {
	return &ApiError{
		RawMessage: err.Error(),
		Code:       err.ErrCode.Code(),
	}
}
