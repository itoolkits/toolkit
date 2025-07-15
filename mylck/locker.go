// Implement some locks, lock type:
//	- TxLocker: use db tx implements the attribute system lock
//	- KeyLocker: lock a string

package mylck

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"math"
	"math/rand"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// Constants
const (
	lockNameLenLimit = 250           // Compress lock name when over this length
	lockDefaultVer   = 0             // Lock default version
	lockVerMaxLimit  = math.MaxInt64 // Lock max version
	lockVerGrowStep  = 1             // Version grow step when lock

	retryWaitMinTime = 10 // Sleep a while when Retry to get lock
	retryWaitMaxTime = 200
	retryTimesLimit  = 10 // Retry many times, return error
)

// Error var
var (
	ErrLockNameDuplicate = errors.New("lock name duplicate error")
	ErrLockVerLag        = errors.New("lock version lag error")

	ErrConcurrentModifyLock = errors.New("concurrent modify lock error")
)

// Global var
var (
	uniqueDuplicateErrorReg, _ = regexp.Compile("\\[1062\\](.*)Duplicate(.*)UNIQUE")

	// Store KeyLocker
	keyLockerMap = &sync.Map{}
)

// KeyLocker - string lock, one string one lock
type KeyLocker struct {
	key   string
	mutex *sync.Mutex
	cnt   *int32
}

// NewKeyLocker - create key locker
func NewKeyLocker(key string) *KeyLocker {
	v, ok := keyLockerMap.Load(key)
	if ok {
		kl := v.(*KeyLocker)
		atomic.AddInt32(kl.cnt, 1)
		return kl
	}
	cnt := int32(1)
	tl := &KeyLocker{
		key:   key,
		mutex: &sync.Mutex{},
		cnt:   &cnt,
	}

	v, ok = keyLockerMap.LoadOrStore(key, tl)
	kl := v.(*KeyLocker)
	if ok {
		atomic.AddInt32(kl.cnt, 1)
	}

	return kl
}

// Lock - lock the key
func (kl *KeyLocker) Lock() {
	kl.mutex.Lock()
}

// Unlock - unlock the key
func (kl *KeyLocker) Unlock() {
	defer kl.mutex.Unlock()

	if atomic.AddInt32(kl.cnt, -1) <= 0 {
		keyLockerMap.Delete(kl.key)
	}
}

// Lock - Lock interface
type Lock interface {
	Lock() (*gorm.DB, error)

	RetryLock(retryTimes int) (*gorm.DB, error)

	BeginTx() *gorm.DB
	CommitTx()
	RollbackTx()
}

// TxLocker - implements Lock, use db tx
type TxLocker struct {
	name         string
	info         string
	localLock    *KeyLocker
	mutexLockDao *MutexLockDao
}

var _ Lock = (*TxLocker)(nil)

// NewTxLocker - create a tx lock
func NewTxLocker(mutexLockDao *MutexLockDao, lockName string) *TxLocker {
	// Set lock info
	info := lockName

	// Compress lock name
	if len(lockName) > lockNameLenLimit {
		lockName = strMD5(lockName)
	}

	// return a locker
	return &TxLocker{
		mutexLockDao: mutexLockDao,
		name:         lockName,
		info:         info,
		localLock:    NewKeyLocker(lockName),
	}
}

// Lock - tx lock do lock, no retry when error
func (tl *TxLocker) Lock() (*gorm.DB, error) {
	return tl.RetryLock(1)
}

// RetryLock - tx lock do lock, retry few times
func (tl *TxLocker) RetryLock(retryTimes int) (*gorm.DB, error) {
	// Fix arguments
	if retryTimes < 1 {
		retryTimes = 1
	}
	if retryTimes > retryTimesLimit {
		retryTimes = retryTimesLimit
	}

	// Get local lock, use KeyLocker
	tl.localLock.Lock()

	var err error
	var db *gorm.DB

retryLoop:
	for i := 0; i < retryTimes; i++ {

		// Sleep a while when begin retry
		if err != nil && err == ErrConcurrentModifyLock {
			sleepRandomMs(retryWaitMinTime, retryWaitMaxTime)
		}

		// Begin db transaction
		db = tl.BeginTx()

		// Try lock
		err = tl.lock()

		// Some error need retry
		switch err {
		case ErrLockNameDuplicate:
			err = ErrConcurrentModifyLock
			continue retryLoop
		case ErrLockVerLag:
			err = ErrConcurrentModifyLock
			continue retryLoop
		default:
			break retryLoop
		}
	}

	return db, err
}

// lock - do lock operation
func (tl *TxLocker) lock() error {
	mutexLock, err := tl.mutexLockDao.GetByName(tl.name)
	if err != nil {
		return err
	}

	// Insert lock
	if mutexLock == nil || mutexLock.ID < 1 {
		mutexLock = &MutexLock{
			Name: tl.name,
			Info: tl.info,
			Ver:  lockDefaultVer,
		}
		err = tl.mutexLockDao.Save(mutexLock)
		// SQL error [1062] [23000]: Duplicate entry 'a' for key 'mutex_lock.name_UNIQUE'
		if err != nil {
			if uniqueDuplicateErrorReg.MatchString(err.Error()) {
				return ErrLockNameDuplicate
			}
			return err
		}
		return err
	}

	// Grow lock version, set default value when grow max
	hisVer := mutexLock.Ver
	if mutexLock.Ver == lockVerMaxLimit {
		mutexLock.Ver = lockDefaultVer
	} else {
		mutexLock.Ver = mutexLock.Ver + lockVerGrowStep
	}

	// Update lock version
	rst, err := tl.mutexLockDao.updateVer(mutexLock.Name, hisVer, mutexLock.Ver)
	if err != nil {
		return err
	}

	// Version lag
	if rst <= 0 {
		return ErrLockVerLag
	}

	return nil
}

// BeginTx - begin db tx
func (tl *TxLocker) BeginTx() *gorm.DB {
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}
	return tl.mutexLockDao.BeginTx(txOptions)
}

// CommitTx - commit db tx, release local lock
func (tl *TxLocker) CommitTx() {
	defer tl.localLock.Unlock()
	_ = tl.mutexLockDao.CommitTx()
}

// RollbackTx - rollback db tx, release local lock
func (tl *TxLocker) RollbackTx() {
	defer tl.localLock.Unlock()
	_ = tl.mutexLockDao.RollbackTx()
}

// strMD5 - string md5
func strMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// waitRandomMs - goroutine sleep random duration
func sleepRandomMs(min, max int) {
	rst := rand.Intn(max - min + 1)
	time.Sleep(time.Duration(rst+min) * time.Millisecond)
}
