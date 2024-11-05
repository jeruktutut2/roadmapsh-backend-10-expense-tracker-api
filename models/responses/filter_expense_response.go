package modelresponses

type FilterExpenseResponse struct {
	Filter string  `json:"filter"`
	Total  float64 `json:"total"`
}

func ToFilterExpenseResponse(filter string, total float64) FilterExpenseResponse {
	return FilterExpenseResponse{
		Filter: filter,
		Total:  total,
	}
}
