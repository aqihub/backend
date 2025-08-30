package api

import (
	"encoding/json"
	"fmt"
	"go-tropic-thunder/pkg/db"
	"go-tropic-thunder/pkg/models"
	"net/http"
)

var metaManager = db.NewMetadataManager()

func InsertDocumentHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		DeviceId  string  `json:"device_id"`
		GpsLat    float64 `json:"gps_lat"`
		GpsLng    float64 `json:"gps_lng"`
		Timestamp int64   `json:"timestamp"`
		TempCel   float64 `json:"temp_cel"`
		Humidity  float64 `json:"humidity"`
		TvocPpb   float64 `json:"tvoc_ppb"`
		Eco2Ppm   float64 `json:"eco2_ppm"`
		AQI       int64   `json:"aqi"`
		IsPublic  bool    `json:"is_public"`
	}

	var request requestBody

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Invalid input: %v", err),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	deviceId := request.DeviceId
	// Convert struct to map using json marshal/unmarshal
	jsonData, _ := json.Marshal(request)
	var requestMap map[string]interface{}
	json.Unmarshal(jsonData, &requestMap)

	cid, err := metaManager.InsertDocument(deviceId, requestMap)
	if err != nil {
		response := models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Failed to insert document: %v", err),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response := models.CommonResponse{
		Status: http.StatusOK,
		Data:   map[string]string{"cid": cid},
	}

	// Set Content-Type header and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		response := models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Failed to encode response: %v", err),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
}

func GetDocumentsAndCIDHandler(w http.ResponseWriter, r *http.Request) {
	var docs models.DocumentsResponse
	var err error
	var data map[string]interface{}

	collection := r.URL.Query().Get("collection_name")
	documentId := r.URL.Query().Get("document_id")
	if documentId != "" {
		data, err = metaManager.GetCIDData(documentId)
		if err != nil {
			response := models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  fmt.Sprintf("Failed to fetch document: %v", err),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		response := models.CommonResponse{
			Status: http.StatusOK,
			Data:   data,
		}
		// Set Content-Type header and encode the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			response := models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Failed to encode response: %v", err),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
	} else if collection != "" {
		docs, err = metaManager.GetDocuments(collection)
		if err != nil {
			response := models.ErrorResponse{
				Status: http.StatusNotFound,
				Error:  fmt.Sprintf("Failed to fetch collections: %v", err),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		response := models.CommonResponse{
			Status: http.StatusOK,
			Data:   docs,
		}
		// Set Content-Type header and encode the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			response := models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Failed to encode response: %v", err),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// Fetch and display all collection names
		collections, err := metaManager.GetAllDocuments()
		if err != nil {
			response := models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Failed to fetch collection names: %v", err),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		response := models.CommonResponse{
			Status: http.StatusOK,
			Data:   collections,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			response := models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Failed to encode response: %v", err),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
	}
}

//func GetPublicDataHandler(w http.ResponseWriter, r *http.Request) {
//	// Fetch all documents
//	documents, err := metaManager.GetAllDocuments()
//	if err != nil {
//		response := models.ErrorResponse{
//			Status: http.StatusInternalServerError,
//			Error:  fmt.Sprintf("Failed to fetch documents: %v", err),
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusInternalServerError)
//		_ = json.NewEncoder(w).Encode(response)
//		return
//	}
//
//	// Initialize a map to store public data for each collection
//	publicData := make(map[string][]map[string]interface{})
//
//	for _, document := range documents.Collections {
//		collectionName := document.CollectionName
//		// Fetch collection data for the current collection name
//		collectionData, err := metaManager.GetDocuments(collectionName)
//		if err != nil {
//			response := models.ErrorResponse{
//				Status: http.StatusInternalServerError,
//				Error:  fmt.Sprintf("Failed to fetch collection data for collection '%s': %v", collectionName, err),
//			}
//
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusInternalServerError)
//			_ = json.NewEncoder(w).Encode(response)
//			return
//		}
//
//		// Get the CID details of the latest document of collection data
//		latestDocumentCID, err := metaManager.GetCIDData(collectionData.LatestDocument)
//		if err != nil {
//			response := models.ErrorResponse{
//				Status: http.StatusInternalServerError,
//				Error:  fmt.Sprintf("Failed to fetch CID for latest document of collection '%s': %v", collectionName, err),
//			}
//
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusInternalServerError)
//			_ = json.NewEncoder(w).Encode(response)
//			return
//		}
//
//		// Store the CID details in the publicData map
//		publicData[collectionName] = append(publicData[collectionName], map[string]interface{}{
//			"cid": latestDocumentCID,
//		})
//	}
//	//}
//
//	// Prepare and send the response
//	response := models.CommonResponse{
//		Status: http.StatusOK,
//		Data:   publicData,
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	if err := json.NewEncoder(w).Encode(response); err != nil {
//		response := models.ErrorResponse{
//			Status: http.StatusInternalServerError,
//			Error:  fmt.Sprintf("Failed to encode response: %v", err),
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusInternalServerError)
//		_ = json.NewEncoder(w).Encode(response)
//		return
//	}
//}
