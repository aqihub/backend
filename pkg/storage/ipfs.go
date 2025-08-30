package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFSClient struct {
	shell *shell.Shell
}

func NewIPFSClient(apiURL string) *IPFSClient {
	return &IPFSClient{shell: shell.NewShell(apiURL)}
}

func (c *IPFSClient) Add(doc map[string]interface{}) (string, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to serialize document: %w", err)
	}

	cid, err := c.shell.Add(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to add document to IPFS: %w", err)
	}
	log.Println("Added document to IPFS:", cid)
	return cid, nil
}

func (c *IPFSClient) Get(cid string) ([]byte, error) {
	reader, err := c.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from IPFS for CID %s: %w", cid, err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from IPFS response: %w", err)
	}

	return data, nil
}
