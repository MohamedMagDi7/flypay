package utils

import (
	"github.com/callicoder/go-docker/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage(t *testing.T) {
	status := 200
	message := "this is an ok message"
	res := Message(status, message)
	assert.Equal(t, 200, res["status"], "status should be identical")
	assert.Equal(t, "this is an ok message", res["message"], "messages should be identical")
}

func TestPerpareFilter(t *testing.T) {
	var filter models.Filter
	filter.Status = "authorised"
	PerpareFilter(&filter)
	assert.Equal(t, 1, filter.AStatusCode, "statusCode should be 1")
	assert.Equal(t, 100, filter.BStatusCode, "statusCode should be 100")

	filter.Status = "decline"
	PerpareFilter(&filter)
	assert.Equal(t, 2, filter.AStatusCode, "statusCode should be 2")
	assert.Equal(t, 200, filter.BStatusCode, "statusCode should be 200")

	filter.Status = "refunded"
	PerpareFilter(&filter)
	assert.Equal(t, 3, filter.AStatusCode, "statusCode should be 3")
	assert.Equal(t, 300, filter.BStatusCode, "statusCode should be 300")
}

/**************************BenchMark****************************************/
func BenchmarkMessage(b *testing.B) {
	status := 200
	message := "this is an ok message"
	for i := 0; i < b.N; i++ {
		Message(status, message)
	}
}

func BenchmarkPerpareFilterPerpareFilter(b *testing.B) {
	var filter models.Filter
	filter.Status = "authorised"
	for i := 0; i < b.N; i++ {
		PerpareFilter(&filter)
	}
}
