package identity

import (
	"crypto/sha1"
	"fmt"

	"github.com/google/uuid"
)

func New() string {
	return uuid.NewString()
}

func DeterministicSlotID(roomID string, startISO string, endISO string) string {
	seed := []byte(fmt.Sprintf("%s|%s|%s", roomID, startISO, endISO))
	return uuid.NewHash(sha1.New(), uuid.NameSpaceURL, seed, 5).String()
}
