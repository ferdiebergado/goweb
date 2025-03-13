package message

var messages = map[string]string{
	"regSuccess": "Thank you for registering. Please check your email for the verification link.",
}

func Get(key string) string {
	msg, ok := messages[key]
	if !ok {
		return "Message not found"
	}
	return msg
}
