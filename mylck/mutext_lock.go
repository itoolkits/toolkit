// mutex lock dao

package mylck

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type MutexLock struct {
	ID        int64     `gorm:"column:id;primary_key" json:"-"`
	Name      string    `gorm:"column:name;unique" json:"name"`
	Info      string    `gorm:"column:info" json:"info"`
	Ver       int64     `gorm:"column:ver" json:"ver"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:null" json:"updatedAt"`
}

func (d *MutexLock) TableName() string {
	return "mutex_lock"
}

type MutexLockDao struct {
	db *gorm.DB
}

// NewMutexLockDao - create mutex lock dao
func NewMutexLockDao(db *gorm.DB) *MutexLockDao {
	return &MutexLockDao{db}
}

// GetByName - get mutex lock by name
func (m *MutexLockDao) GetByName(name string) (*MutexLock, error) {
	var ret *MutexLock
	err := m.db.Where(" name = ? ", name).First(&ret).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if ret.ID < 1 {
		return nil, nil
	}
	return ret, nil
}

// Save - save lock
func (m *MutexLockDao) Save(mut *MutexLock) error {
	return m.db.Save(mut).Error
}

// UpdateByInfo - update mutex lock by info
func (m *MutexLockDao) UpdateByInfo(name, info string, mut *MutexLock) (int64, error) {
	db := m.db.Where(" name = ? and info = ? ", name, info).Updates(mut)
	rows := db.RowsAffected
	err := db.Error
	return rows, err
}

// updateVer - update lock version, return the number of update rows
func (m *MutexLockDao) updateVer(name string, hisVer, preVer int64) (int64, error) {
	d := m.db.Model(&MutexLock{}).Where(" name = ? and ver = ? ", name, hisVer).UpdateColumn("ver", preVer)
	if d.Error != nil {
		return 0, d.Error
	}
	return d.RowsAffected, nil
}

// BeginTx - begin transaction
func (m *MutexLockDao) BeginTx(opt *sql.TxOptions) *gorm.DB {
	return m.db.Begin(opt)
}

// CommitTx - commit transaction
func (m *MutexLockDao) CommitTx() error {
	return m.db.Commit().Error
}

// RollbackTx - rollback transaction
func (m *MutexLockDao) RollbackTx() error {
	return m.db.Rollback().Error
}
