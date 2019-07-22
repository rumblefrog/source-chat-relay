package protocol

const (
	MAX_BUFFER_LENGTH = 1024
)

type MessageType uint8

const (
	MessageInvalid MessageType = iota
	MessageAuthenticate
	MessageAuthenticateResponse
	MessageChat
	MessageEvent
	MessageTypeCount
)

type AuthenticateResponse uint8

const (
	AuthenticateInvalid AuthenticateResponse = iota
	AuthenticateSuccess
	AuthenticateDenied
	AuthenticateResponseCount
)

type IdentificationType uint8

const (
	IdentificationInvalid IdentificationType = iota
	IdentificationSteam
	IdentificationDiscord
	IdentificationTypeCount
)

func ParseMessageType(t uint8) MessageType {
	if t >= uint8(MessageTypeCount) {
		return MessageInvalid
	}

	return MessageType(t)
}

func ParseAuthenticateResponse(t uint8) AuthenticateResponse {
	if t >= uint8(AuthenticateResponseCount) {
		return AuthenticateInvalid
	}

	return AuthenticateResponse(t)
}
func ParseIdentificationType(t uint8) IdentificationType {
	if t >= uint8(IdentificationTypeCount) {
		return IdentificationInvalid
	}

	return IdentificationType(t)
}
