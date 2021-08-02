package v1

import "github.com/xybor/xychat/helpers/api"

func NewAPIResponse(data interface{}) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno": 0,
	}
	return api.APIResponse{Data: &data, Meta: &meta}
}

func NewAPIError(errno int, err string) api.APIResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno": errno,
		"error": err,
	}

	return api.APIResponse{Meta: &meta}
}
