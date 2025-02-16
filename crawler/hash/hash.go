//---------------------------------------//
// A minimal abstraction over crypto/md5 //
//---------------------------------------//

package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type Hash [16]byte

// Returns md5 hash of a string
func Hashs(s string) Hash {
	return md5.Sum([]byte(s))
}

func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}
