package utils

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

// Snowflake struct represents the Snowflake ID generator
type Snowflake struct {
	machineID     string
	hashedMachine uint32
	lastTimestamp int64
	sequence      int
	mu            sync.Mutex
}

// NewSnowflake creates a new Snowflake instance
func NewSnowflake(machineID string) *Snowflake {
	hashedMachine := hashString(machineID)
	return &Snowflake{
		machineID:     machineID,
		hashedMachine: hashedMachine,
	}
}

// hashString generates a hash value for a given string
func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// GenerateID generates a unique ID based on Snowflake algorithm
func (s *Snowflake) GenerateID() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTimestamp := time.Now().UnixNano() / 1e6

	if currentTimestamp < s.lastTimestamp {
		fmt.Println("Warning: Clock moved backwards. Waiting until the next millisecond.")
		for currentTimestamp <= s.lastTimestamp {
			currentTimestamp = time.Now().UnixNano() / 1e6
		}
	}

	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & ((1 << 12) - 1)
		if s.sequence == 0 {
			// Sequence overflow, wait until next millisecond
			for currentTimestamp <= s.lastTimestamp {
				currentTimestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	// ID format: [timestamp][hashedMachineID][sequence]
	id := (currentTimestamp << 22) | (int64(s.hashedMachine) << 10) | int64(s.sequence)

	return fmt.Sprintf("%s%d", s.machineID, id)
}
