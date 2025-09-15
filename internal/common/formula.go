package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type IdGenerator struct {
	counter uint64
}

// NewOrderIDGenerator creates a new order ID generator
func NewIDGenerator() *IdGenerator {
	return &IdGenerator{
		counter: 0,
	}
}

func (g *IdGenerator) GenerateCommonID(identifier string) (string, error) {
	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixMilli()

	// Atomic increment counter to ensure uniqueness even with same timestamp
	counter := atomic.AddUint64(&g.counter, 1)

	// Generate 4 random bytes for additional entropy
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		logrus.Error(err)
		return "", NewError(err, ErrConflict)
	}
	randomHex := hex.EncodeToString(randomBytes)

	// Format: ORD + timestamp + counter + random
	id := fmt.Sprintf("%v%d%04d%s", identifier, timestamp, counter%10000, randomHex)

	return id, nil
}
