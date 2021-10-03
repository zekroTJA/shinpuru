package vote

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/bwmarrin/snowflake"
)

func HashUserID(userID string, salt []byte) (hash string, err error) {
	sid, err := snowflake.ParseString(userID)
	if err != nil {
		return
	}

	idb := big.NewInt(sid.Int64() & int64(^uint(0)>>(64-48))).Bytes()
	comb := append(idb, salt...)
	hash = fmt.Sprintf("%x", sha256.Sum256(comb))

	return
}
