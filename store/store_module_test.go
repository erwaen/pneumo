package store
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStoreService = &StorageService{}

func init(){
	testStoreService = CreateStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.valkeyClient != nil, "Valkey client should be initialized") 
}

func TestStorePneumo(t *testing.T) {
	fullUrl := "http://example.com/pneumo"
	pneumo := "http://pneumo.com/pneumo123"

	err := testStoreService.storePneumo(fullUrl, pneumo)
	assert.NoError(t, err, "Should store pneumo without error")

	retrievedUrl, err := testStoreService.retrieveFromPneumo(pneumo)
	assert.NoError(t, err, "Should retrieve full URL without error")
	assert.Equal(t, fullUrl, retrievedUrl, "Retrieved URL should match stored URL")

	// Try to retrieve a non-existent pneumo
	_, err = testStoreService.retrieveFromPneumo("nonexistent")
	assert.Error(t, err, "Should return error for non-existent pneumo")
}

func TestGetNoExist(t *testing.T) {
	_, err := testStoreService.retrieveFromPneumo("nonexistent")
	assert.Error(t, err, "Should return error for non-existent key")
}
