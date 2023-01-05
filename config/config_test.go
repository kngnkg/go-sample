package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	assert.Equal(t, wantPort, got.Port)
	wantEnv := "dev"
	assert.Equal(t, wantEnv, got.Env)
}
