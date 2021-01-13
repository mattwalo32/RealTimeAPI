package messages

func DoesMessageRequireResponse(msg Message) bool {
	if msg.GetMessageType() == MESSAGE_ACKNOWLEDGE {
		return false
	}

	return msg.IsResponseRequired()
}

func IsVerbose(msg Message) bool {
	switch msg.GetMessageType() {
	case MESSAGE_ACKNOWLEDGE:
		return true

	case MESSAGE_JOIN_SERVER:
		return true

	case MESSAGE_FIND_ROOM:
		return true
	}

	return false
}
