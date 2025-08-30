package models

type CommonResponse struct {
	Status uint64      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status uint64 `json:"status"`
	Error  string `json:"error"`
}

type AllDocumentsResponse struct {
	Collections []CollectionData `json:"collections"`
}

type CollectionData struct {
	CollectionName string           `json:"collection_name"`
	CollectionData DocumentResponse `json:"collection_data"`
}

type DocumentResponse struct {
	LatestDocument string   `json:"latest_document"`
	Documents      []string `json:"documents"`
}

type DocumentsResponse struct {
	LatestDocument string   `json:"latest_document"`
	Documents      []string `json:"documents"`
}
