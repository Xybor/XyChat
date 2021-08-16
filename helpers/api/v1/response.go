package v1

import (
	"github.com/xybor/xychat/helpers/api"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
)

// NewAPIResponse creates a success API response with a data.
func NewAPIResponse(data interface{}) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   0,
	}

	if data == nil {
		return api.CreateAPIResponse(nil, &meta)
	}

	return api.CreateAPIResponse(&data, &meta)
}

// NewAPIError creates a failure API response.
func NewAPIError(se xyerrors.XyError) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   se.Errno(),
		"error":   se.Error(),
	}

	return api.CreateAPIResponse(nil, &meta)
}
