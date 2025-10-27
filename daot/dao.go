// dao util

package daot

import (
	"errors"

	"gorm.io/gorm"
)

type Dao[T any] struct {
	db *gorm.DB
}

// NewDao create dao
func NewDao[T any](db *gorm.DB) *Dao[T] {
	return &Dao[T]{db: db}
}

// Get first record
func (d *Dao[T]) Get(condition *T) (*T, error) {
	var rst T
	err := d.db.Where(condition).First(&rst).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &rst, nil
}

// GetList get record list
func (d *Dao[T]) GetList(condition *T) ([]*T, error) {
	var rst []*T
	err := d.db.Where(condition).Find(&rst).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return rst, nil
}

// Save record
func (d *Dao[T]) Save(ele *T) error {
	return d.db.Save(ele).Error
}

// Delete record
func (d *Dao[T]) Delete(condition *T) error {
	return d.db.Delete(condition).Error
}

// BatchInsert batch insert
func (d *Dao[T]) BatchInsert(eleList []*T, batchSize int) error {
	return d.db.CreateInBatches(eleList, batchSize).Error
}

// BatchDelete batch delete
func (d *Dao[T]) BatchDelete(eleList []*T, batchSize int) error {
	for i := 0; i < len(eleList); i += batchSize {
		end := i + batchSize
		if end > len(eleList) {
			end = len(eleList)
		}
		err := d.db.Delete(eleList[i:end]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// Where find
func (d *Dao[T]) Where(query interface{}, args ...interface{}) ([]*T, error) {
	var rst []*T
	err := d.db.Where(query, args...).Find(&rst).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return rst, nil
}
