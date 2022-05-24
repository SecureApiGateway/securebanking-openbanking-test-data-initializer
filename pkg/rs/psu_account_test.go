package rs

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestUnmarshalling(t *testing.T) {
	data := []byte(`
            {"_id":"d5c3dbbd-4803-45d5-accc-1b03bed1b63c","_rev":"0000000038c1c572","userName":"psu","accountStatus":"active","givenName":"PSU","sn":"Payment Services User","mail":"psu@acme.com"}
        `)
	dummy := &UserResponse{}
	err := json.Unmarshal(data, dummy)
	if err != nil {
		fmt.Printf("An error occured: %v\n", err)
		os.Exit(1)
	}

	// we want to print the field names as well
	fmt.Printf("Result %v\n", dummy.UserId)
}
