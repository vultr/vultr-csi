// Package vultruserdata provides functionality calling the vultr virtual
// machine user data server and using the results to check if a node is a part
// of a vultr kubernetes engine cluster
package vultruserdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	vultrUserDataURL = "http://169.254.169.254/latest/user-data"
	requestTimeout   = 5 * time.Second
)

func IsVKE() bool {
	ud := NewUserData()
	if err := ud.get(); err != nil {
		return false
	}

	if ud.Data.VKE.NodeID != "" {
		return true
	}

	return false
}

type UserData struct {
	Data struct {
		VKE struct {
			NodeID string `json:"node_id"`
		} `json:"vke"`
	} `json:"data"`
}

func NewUserData() *UserData {
	return &UserData{}
}

func (u *UserData) get() error {
	req, err := http.NewRequest("GET", vultrUserDataURL, nil)
	if err != nil {
		return fmt.Errorf("error creating http request : %v", err)
	}

	client := &http.Client{Timeout: requestTimeout}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error in http client request : %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status %q not ok : %v", resp.StatusCode, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading http request body : %v", err)
	}

	if err := json.Unmarshal(body, &u); err != nil {
		return fmt.Errorf("error unmarshalling request body : %v", err)
	}

	return nil
}
