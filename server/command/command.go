package command

type CommandType string

const (
	NONE            CommandType = "none"
	SUCCESS_MESSAGE             = "success_msg"
	ERROR_MESSAGE               = "error_msg"
	REDIRECT                    = "redirect"
)

type StringPayload struct {
	Payload string
}

type RedirectPayload struct {
	Target string
}

type JSONCommand struct {
	ReturnCode int
	Command    CommandType
	Payload    interface{}
}

func NewRedirectCommand(target string) *JSONCommand {
	return &JSONCommand{
		ReturnCode: 200,
		Command:    REDIRECT,
		Payload: &RedirectPayload{
			Target: target,
		},
	}
}

func NewValidationErrorCommand(message string) *JSONCommand {
	return &JSONCommand{
		ReturnCode: 500,
		Command:    ERROR_MESSAGE,
		Payload: &StringPayload{
			message,
		},
	}
}
