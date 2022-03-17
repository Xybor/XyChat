package helpers

import xyerrors "github.com/xybor/xychat/xyerrors"

type response struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// createResponse is a constructor for create an API response.  APIResponse
// is a struct with two field: data (store the API's result) and meta (store
// the API's metadata)
func createResponse(
	data *interface{},
	meta *map[string]interface{},
) response {
	return response{Data: data, Meta: meta}
}

// NewResponse creates a success API response with a data.
func NewResponse(data interface{}) response {
	meta := map[string]interface{}{
		"errno": 0,
	}

	if data == nil {
		return createResponse(nil, &meta)
	}

	return createResponse(&data, &meta)
}

// NewErrorResponse creates a failure API response.
func NewErrorResponse(xerr xyerrors.XyError) response {
	meta := map[string]interface{}{
		"errno": xerr.Errno(),
		"error": xerr.Error(),
	}

	return createResponse(nil, &meta)
}
