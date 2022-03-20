package model

type Level string

const (
	Alert  Level = "ALERT"
	Logout Level = "LOGOUT"
	Input  Level = "INPUT"
)

type Error struct {
	Message     string `json:"error"`
	Level       Level  `json:"level"`
	Description string `json:"description"`
}

func AccessError() Error {
	return Error{
		Message: "Dazu hast du keine Berechtigung",
		Level:   Alert,
	}
}

func InternalError() Error {
	return Error{
		Message: "Es ist ein interner Fehler aufgetreten. Bitte wende dich an einen Administrator",
		Level:   Alert,
	}
}

func ProcessError() Error {
	return Error{
		Message: "Diese Anfrage kann nicht verarbeitet werden.",
		Level:   Alert,
	}
}

func LogoutError() Error {
	return Error{
		Level: Logout,
	}
}

func AlertError(message string) Error {
	return Error{
		Message: message,
		Level:   Alert,
	}
}

func InputError(message string, input string) Error {
	return Error{
		Message:     message,
		Level:       Input,
		Description: input,
	}
}
