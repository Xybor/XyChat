package v1

import "github.com/xybor/xychat/helpers/ws"

// NewWSResponse creates a success WS response with a data.
func NewWSResponse(data interface{}) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   0,
	}

	if data == nil {
		return ws.CreateWSResponse(nil, &meta)
	}

	return ws.CreateWSResponse(&data, &meta)
}

// NewWSError creates a failure WS response.
func NewWSError(errno int, err string) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   errno,
		"error":   err,
	}

	return ws.CreateWSResponse(nil, &meta)
}
