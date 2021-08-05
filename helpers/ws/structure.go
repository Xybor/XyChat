package ws

type WSResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// CreateWSResponse is a constructor for create an WS response.  WSResponse is
// a struct with two field: data (store the WS's result) and meta (store the
// WS's metadata)
func CreateWSResponse(
	data *interface{},
	meta *map[string]interface{},
) WSResponse {
	return WSResponse{Data: data, Meta: meta}
}
