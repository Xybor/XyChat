package api

type APIResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// CreateAPIResponse is a constructor for create an API response.  APIResponse
// is a struct with two field: data (store the API's result) and meta (store 
// the API's metadata)
func CreateAPIResponse(
	data *interface{},
	meta *map[string]interface{},
) APIResponse {
	return APIResponse{Data: data, Meta: meta}
}
