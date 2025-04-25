// md5 checker

package ctxt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

var akSKMap = make(map[string]string)

// AddAKSKCache add ak sk key
func AddAKSKCache(ak, sk string) {
	akSKMap[ak] = sk
}

// MD5Check check md5
func MD5Check(ak, sk string, ts int64) error {
	if ak == "" {
		return fmt.Errorf("ak can not blank")
	}

	if sk == "" {
		return fmt.Errorf("sk can not blank")
	}

	cur := time.Now()
	endTime := cur.Add(time.Minute * 5).Unix()
	startTime := cur.Add(-time.Minute * 5).Unix()

	if ts > endTime || ts < startTime {
		return fmt.Errorf("timestamp out of range")
	}
	rsk, h := akSKMap[ak]
	if !h {
		return fmt.Errorf("no authorization")
	}
	hs := hopeAuth(ak, rsk, ts)
	if hs != sk {
		return fmt.Errorf("no authorization")
	}
	return nil
}

// hopeAuth - get md5 str
func hopeAuth(ak, sk string, ts int64) string {
	str := fmt.Sprintf("%s%s%d", ak, sk, ts)
	hash := md5.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
