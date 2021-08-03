package v1

import "github.com/xybor/xychat/helpers/ws"

func NewWSResponse(data interface{}) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno": 0,
	}
	return ws.CreateWSResponse(&data, &meta)
}

func NewWSError(errno int, err string) ws.WSResponse {
	meta := map[string]interface{}{
		"version": 1,
		"errno": errno,
		"error": err,
	}

	return ws.CreateWSResponse(nil, &meta)
}