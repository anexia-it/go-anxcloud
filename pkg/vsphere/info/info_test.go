package info

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiskSize_UnmarshalJSON(t *testing.T) {
	diskJson := []byte("{ \"disk_gb\": 10.00 }")
	var diskInfo DiskInfo
	var expectedSize int = 10

	err := json.Unmarshal(diskJson, &diskInfo)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedSize, diskInfo.DiskGB)
}