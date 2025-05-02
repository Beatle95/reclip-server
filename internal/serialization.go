package internal

func SerializeSync(thisData ClientData, otherData []ClientData) []byte {
	// TODO:
	return nil
}

func SerializeClientId(id string) []byte {
	// TODO:
	return nil
}

func SerializeClientData(clientData *ClientData) []byte {
	// TODO:
	return nil
}

func SerializeTextUpdate(id string, text string) []byte {
	// TODO:
	return nil
}

func SerializeError(errorText string) []byte {
	// TODO:
	return nil
}

func DeserializeClientId(data []byte) (string, error) {
	// TODO:
	return "", nil
}

func DeserializeText(data []byte) (string, error) {
	// TODO:
	return "", nil
}

func DeserializeClientData(data []byte) (ClientData, error) {
	// TODO:
	return ClientData{}, nil
}
