package ws

type WSResponse struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

func CreateWSResponse(
	data *interface{},
	meta *map[string]interface{},
) WSResponse {
	return WSResponse{Data: data, Meta: meta}
}
