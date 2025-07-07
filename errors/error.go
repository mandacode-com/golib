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
