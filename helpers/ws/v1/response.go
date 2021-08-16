package v1

import ("github.com/xybor/xychat/helpers/ws"
xyerrors "github.com/xybor/xychat/xyerrors/v1")

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
func NewWSError(xerr xyerrors.XyError) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   xerr.Errno(),
		"error":   xerr.Error(),
	}

	return ws.CreateWSResponse(nil, &meta)
}
