package v1

import "github.com/xybor/xychat/helpers/api"

func NewEmptyAPIResponse() api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   0,
	}
	return api.CreateAPIResponse(nil, &meta)
}

func NewAPIResponse(data interface{}) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   0,
	}
	return api.CreateAPIResponse(&data, &meta)
}

func NewAPIError(errno int, err string) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno":   errno,
		"error":   err,
	}

	return api.CreateAPIResponse(nil, &meta)
}
