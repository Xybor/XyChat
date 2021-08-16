package api

import xyerrors "github.com/xybor/xychat/xyerrors"

type apiResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// createAPIResponse is a constructor for create an API response.  APIResponse
// is a struct with two field: data (store the API's result) and meta (store
// the API's metadata)
func createAPIResponse(
	data *interface{},
	meta *map[string]interface{},
) apiResponse {
	return apiResponse{Data: data, Meta: meta}
}

// NewAPIResponse creates a success API response with a data.
func NewAPIResponse(data interface{}) apiResponse {
	meta := map[string]interface{}{
		"errno": 0,
	}

	if data == nil {
		return createAPIResponse(nil, &meta)
	}

	return createAPIResponse(&data, &meta)
}

// NewAPIError creates a failure API response.
func NewAPIError(se xyerrors.XyError) apiResponse {
	meta := map[string]interface{}{
		"errno": se.Errno(),
		"error": se.Error(),
	}

	return createAPIResponse(nil, &meta)
}
