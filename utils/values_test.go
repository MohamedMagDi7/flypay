package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnumUnknownToString(t *testing.T) {
	unknownString := Unknown.String()
	assert.Equal(t, "unknown", unknownString, "string should be 'unknown'")
}

func TestEnumAuthorizedToString(t *testing.T) {
	authorizedString := Authorised.String()
	assert.Equal(t, "authorised", authorizedString, "string should be 'authorised'")
	authorizedString = B_Authorised.String()
	assert.Equal(t, "authorised", authorizedString, "string should be 'authorised'")
}

func TestEnumDeclineToString(t *testing.T) {
	declineString := Decline.String()
	assert.Equal(t, "decline", declineString, "string should be 'decline'")
	declineString = B_Decline.String()
	assert.Equal(t, "decline", declineString, "string should be 'decline'")
}

func TestEnumRefundedToString(t *testing.T) {
	refundedString := Refunded.String()
	assert.Equal(t, "refunded", refundedString, "string should be 'refunded'")
	refundedString = B_Refunded.String()
	assert.Equal(t, "refunded", refundedString, "string should be 'refunded'")
}
