package errors

type PublicError interface {
	error
	Public() string
}

type AppError struct {
	msg       string
	publicMsg string
}

func (e *AppError) Error() string {
	return e.msg
}

func (e *AppError) Public() string {
	if e.publicMsg != "" {
		return e.publicMsg
	}
	return "internal error"
}

func NewPublicError(msg string, publicMsg ...string) error {
	if len(publicMsg) > 0 {
		return &AppError{
			msg:       msg,
			publicMsg: publicMsg[0],
		}
	}
	return &AppError{
		msg: msg,
	}
}

func Join(err error, msg string) error {
	if err == nil {
		return nil
	}
	if publicErr, ok := err.(PublicError); ok {
		return &AppError{
			msg:       msg + "\n\t" + err.Error(),
			publicMsg: publicErr.Public(),
		}
	}
	return &AppError{
		msg: msg + "\n\t" + err.Error(),
	}
}

func IsPublicError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(PublicError)
	return ok
}

func IsAppError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*AppError)
	return ok
}
