package v1

import "github.com/xybor/xychat/helpers/ws"

// NewEmptyWSResponse creates a success and non-data WS response.
func NewEmptyWSResponse(data interface{}) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   0,
	}
	return ws.CreateWSResponse(nil, &meta)
}

// NewWSResponse creates a success WS response with a data.
func NewWSResponse(data interface{}) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   0,
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
