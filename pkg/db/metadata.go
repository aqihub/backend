package db

import (
	"encoding/json"
	"fmt"
	"go-tropic-thunder/pkg/models"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"

	"go-tropic-thunder/pkg/storage"

	"golang.org/x/net/context"
)

type MetadataManager struct {
	mu          sync.Mutex
	redisClient *redis.Client
	ipfsClient  *storage.IPFSClient
}

func NewMetadataManager() *MetadataManager {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PORT"), // Adjust your Redis server address here
		Password: "",                      // No password by default
		DB:       0,                       // Default DB
	})

	// Test the Redis connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &MetadataManager{
		redisClient: client,
		ipfsClient:  storage.NewIPFSClient("localhost:5001"),
	}
}

func (m *MetadataManager) InsertDocument(deviceId string, doc map[string]interface{}) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	db := os.Getenv("IPFS_DATABASE")

	// Check if database exists in Redis, if not, create it
	_, err := m.redisClient.HGet(context.Background(), db, deviceId).Result()
	if err == redis.Nil {
		// No such field (deviceId) exists, create it
		m.redisClient.HSet(context.Background(), db, deviceId, "[]")
	} else if err != nil {
		return "", fmt.Errorf("failed to check Redis for existing data: %v", err)
	}

	// Add the document to IPFS
	cid, err := m.ipfsClient.Add(doc)
	if err != nil {
		return "", fmt.Errorf("failed to store document in IPFS: %v", err)
	}

	// Retrieve existing documents for the deviceId from Redis
	collsStr, err := m.redisClient.HGet(context.Background(), db, deviceId).Result()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve data from Redis: %v", err)
	}

	// Unmarshal the documents for the deviceId
	var colls []string
	if collsStr != "" {
		if err := json.Unmarshal([]byte(collsStr), &colls); err != nil {
			return "", fmt.Errorf("failed to unmarshal Redis data: %v", err)
		}
	}

	// Append the new CID to the list of documents
	colls = append(colls, cid)

	// Store the updated list of documents back to Redis
	collsData, err := json.Marshal(colls)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data for Redis: %v", err)
	}

	if err := m.redisClient.HSet(context.Background(), db, deviceId, collsData).Err(); err != nil {
		return "", fmt.Errorf("failed to store updated documents in Redis: %v", err)
	}

	return cid, nil
}

func (m *MetadataManager) GetDocuments(collection string) (models.DocumentsResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	db := os.Getenv("IPFS_DATABASE")
	collsStr, err := m.redisClient.HGet(context.Background(), db, collection).Result()
	if err == redis.Nil {
		return models.DocumentsResponse{}, fmt.Errorf("collection %s does not exist", collection)
	} else if err != nil {
		return models.DocumentsResponse{}, fmt.Errorf("failed to retrieve collection data from Redis: %v", err)
	}

	var colls []string
	if err := json.Unmarshal([]byte(collsStr), &colls); err != nil {
		return models.DocumentsResponse{}, fmt.Errorf("failed to unmarshal Redis data for collection %s: %v", collection, err)
	}

	latest_document := colls[len(colls)-1]
	response := models.DocumentsResponse{
		LatestDocument: latest_document,
		Documents:      colls,
	}
	return response, nil
}

func (m *MetadataManager) GetAllDocuments() (models.AllDocumentsResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	db := os.Getenv("IPFS_DATABASE")
	keys, err := m.redisClient.HKeys(context.Background(), db).Result()
	if err != nil {
		return models.AllDocumentsResponse{}, fmt.Errorf("failed to retrieve collections from Redis: %v", err)
	}

	// Initialize a slice to hold all collection data
	allDocuments := models.AllDocumentsResponse{
		Collections: []models.CollectionData{},
	}

	for _, collectionName := range keys {
		collsStr, err := m.redisClient.HGet(context.Background(), db, collectionName).Result()
		if err != nil {
			continue // Skip collections that have no data
		}

		var colls []string
		if err := json.Unmarshal([]byte(collsStr), &colls); err != nil {
			continue // Skip invalid collections
		}

		if len(colls) > 0 {
			latestDocument := colls[len(colls)-1]

			allDocuments.Collections = append(allDocuments.Collections, models.CollectionData{
				CollectionName: collectionName,
				CollectionData: models.DocumentResponse{
					LatestDocument: latestDocument,
					Documents:      colls,
				},
			})
		}
	}

	if len(allDocuments.Collections) == 0 {
		return models.AllDocumentsResponse{}, fmt.Errorf("no collections found in database %s", db)
	}

	return allDocuments, nil
}

func (m *MetadataManager) GetCIDData(cid string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := m.ipfsClient.Get(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from IPFS: %v", err)
	}

	// Decode JSON into a map
	decodedData := map[string]interface{}{}
	if err := json.Unmarshal(data, &decodedData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON data for CID %s: %v", cid, err)
	}

	return decodedData, nil
}
