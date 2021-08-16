package ws

import "github.com/xybor/xychat/xyerrors"

type wsResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// createWSResponse is a constructor for create an WS response.  WSResponse is
// a struct with two field: data (store the WS's result) and meta (store the
// WS's metadata)
func createWSResponse(
	data *interface{},
	meta *map[string]interface{},
) wsResponse {
	return wsResponse{Data: data, Meta: meta}
}

// NewWSResponse creates a success WS response with a data.
func NewWSResponse(data interface{}) wsResponse {
	meta := map[string]interface{}{
		"errno": 0,
	}

	if data == nil {
		return createWSResponse(nil, &meta)
	}

	return createWSResponse(&data, &meta)
}

// NewWSError creates a failure WS response.
func NewWSError(xerr xyerrors.XyError) wsResponse {
	meta := map[string]interface{}{
		"errno": xerr.Errno(),
		"error": xerr.Error(),
	}

	return createWSResponse(nil, &meta)
}
