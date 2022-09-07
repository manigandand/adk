package respond

// Msg returns new json with message
func Msg(msg interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
	}
}

// response holds the handlerfunc response
type response struct {
	Data interface{} `json:"data,omitempty"`
	Meta Meta        `json:"meta"`
}

// Meta holds the status of the request informations
// TODO: add meta information for paginations
type Meta struct {
	Status int `json:"status_code"`
}
