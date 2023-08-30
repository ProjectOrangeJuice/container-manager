package shared

type EventData struct {
	Request string
	Result  []byte
}

type UpdateResult struct {
	ErrorReason string
}
