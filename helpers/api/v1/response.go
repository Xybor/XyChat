package v1

import "github.com/xybor/xychat/helpers/api"

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
func NewAPIError(errno int, err string) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   errno,
		"error":   err,
	}

	return api.CreateAPIResponse(nil, &meta)
}
