package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	
	assert.True(t, true)
}

type MockedDao struct {
	mock.Mock
}
