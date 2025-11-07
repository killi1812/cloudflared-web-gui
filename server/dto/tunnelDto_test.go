package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	// Import the packages your DTOs depend on
	// I am assuming the paths based on your provided code
	"github.com/killi1812/cloudflared-web-gui/dto"
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/format"
)

// --- Mocks & Test Setups ---

// Helper function to create a string pointer
func strPtr(s string) *string {
	return &s
}

// --- Tests ---

func TestTunnelDto_FromModel(t *testing.T) {
	// 1. Arrange
	testUUID := uuid.New()
	testTime, _ := time.Parse(format.DateTimeFormat, "2025-11-07 10:30:00")

	// This is our mock model data
	mockTunnel := model.Tunnel{
		Id:        testUUID,
		Name:      "my-test-tunnel",
		IsRunning: true,
		CreatedAt: testTime,
		DeletedAt: testTime,
	}

	// This is the DTO we want to populate
	var dto dto.TunnelDto

	dto.FromModel(mockTunnel)

	// 3. Assert
	assert.Equal(t, testUUID.String(), dto.Id)
	assert.Equal(t, "my-test-tunnel", dto.Name)
	assert.Equal(t, "2025-11-07 10:30:00", dto.CreatedAt)
	assert.Equal(t, "2025-11-07 10:30:00", dto.DeletedAt)
	assert.True(t, dto.IsRunning)
	assert.Empty(t, dto.DnsRecords) // Ensure DNS records are not unexpectedly populated
}

func TestDnsRecordDto_FromModel(t *testing.T) {
	// 1. Arrange
	testTime, _ := time.Parse(format.DateTimeFormat, "2025-11-07 11:00:00")
	mockSettings := map[string]string{"setting": "value"}
	mockMeta := map[string]any{"meta": true}
	mockTags := []any{"tag1", "tag2"}

	mockRecord := model.DnsRecord{
		Id:         "dns-id-123",
		Name:       "test.example.com",
		Type:       "A",
		Content:    "1.2.3.4",
		Proxiable:  true,
		Proxied:    true,
		Ttl:        300,
		Settings:   mockSettings,
		Meta:       mockMeta,
		Commnet:    strPtr("test comment"),
		Tags:       mockTags,
		CreatedAt:  testTime,
		ModifiedOn: testTime,
	}

	var dto dto.DnsRecordDto

	// 2. Act
	dto.FromModel(mockRecord)

	// 3. Assert
	assert.Equal(t, "dns-id-123", dto.Id)
	assert.Equal(t, "test.example.com", dto.Name)
	assert.Equal(t, "A", dto.Type)
	assert.Equal(t, "1.2.3.4", dto.Content)
	assert.True(t, dto.Proxiable)
	assert.True(t, dto.Proxied)
	assert.Equal(t, 300, dto.Ttl)
	assert.Equal(t, mockSettings, dto.Settings)
	assert.Equal(t, mockMeta, dto.Meta)
	assert.NotNil(t, dto.Commnet)
	assert.Equal(t, "test comment", *dto.Commnet)
	assert.Equal(t, mockTags, dto.Tags)
	assert.Equal(t, "2025-11-07 11:00:00", dto.CreatedAt)
	assert.Equal(t, "2025-11-07 11:00:00", dto.ModifiedOn)
}

func TestDnsRecordDto_FromModel_NilComment(t *testing.T) {
	// 1. Arrange
	// Test the specific case where the comment pointer is nil
	mockRecord := model.DnsRecord{
		Id:      "dns-id-456",
		Commnet: nil,
	}

	var dto dto.DnsRecordDto

	// 2. Act
	dto.FromModel(mockRecord)

	// 3. Assert
	// Ensure the nil pointer wasn't incorrectly handled
	assert.Nil(t, dto.Commnet)
}

func TestArrDnsRecordDto_FromModel(t *testing.T) {
	// 1. Arrange
	mockRecords := []model.DnsRecord{
		{Id: "id-1", Name: "record1.com"},
		{Id: "id-2", Name: "record2.com"},
	}

	// Initialize the slice. This is important for testing the fix.
	var dtoArr dto.ArrDnsRecordDto

	// 2. Act
	dtoArr.FromModel(mockRecords)

	// 3. Assert
	assert.NotNil(t, dtoArr, "The slice should be initialized, not nil")
	assert.Len(t, dtoArr, 2)
	assert.Equal(t, "id-1", dtoArr[0].Id)
	assert.Equal(t, "record1.com", dtoArr[0].Name)
	assert.Equal(t, "id-2", dtoArr[1].Id)
	assert.Equal(t, "record2.com", dtoArr[1].Name)
}

func TestArrDnsRecordDto_FromModel_Empty(t *testing.T) {
	// 1. Arrange
	mockRecords := []model.DnsRecord{} // Empty slice
	var dtoArr dto.ArrDnsRecordDto     // Nil slice

	// 2. Act
	dtoArr.FromModel(mockRecords)

	// 3. Assert
	assert.NotNil(t, dtoArr, "The slice should be initialized, not nil")
	assert.Len(t, dtoArr, 0)
}
