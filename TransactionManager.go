package core

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type TransactionManager struct {
	once         sync.Once
	db           *gorm.DB
	tx           *gorm.DB
	transCounter int64
}

func (t *TransactionManager) Transaction(callback func()) {

	t.begin()

	defer func() {
		if err := recover(); err != nil {
			t.rollback()
		}
	}()

	// get the error and
	callback()

	t.commit()

}

func (t *TransactionManager) begin() {
	// first time no transaction start yet
	if t.transCounter == 0 {
		// create a internal ref tx
		t.tx = t.db.Begin()
	} else if t.transCounter >= 1 && t.supportSavePoint() {
		// after the first time we create a savepoint if the db were supported
		t.createSavePoint()
	}

	// increase arc var
	t.transCounter++

	// @todo maybe fire [beganTransaction] event
}

func (t *TransactionManager) commit() {
	if t.transCounter == 1 {
		t.tx.Commit()
	}

	// trigger this to be maintains the ref counting
	t.transCounter = t.max(0, t.transCounter-1)

	// @todo maybe fire [committed] event
}

func (t *TransactionManager) rollback() {
	if t.transCounter == 0 {
		// create a internal ref tx
		t.tx = t.db.Rollback()
	} else if t.transCounter >= 1 && t.supportSavePoint() {
		t.removeSavePoint()
	}

}

// get the number of active transactions.
func (t *TransactionManager) Level() int64 {
	return t.transCounter
}

// create a save point within the database
func (t *TransactionManager) createSavePoint() {
	// @todo db.execRaw('added up savepoint')

}

func (t *TransactionManager) removeSavePoint() {
	// @todo db.exeRaw('rollback savepoint')
}

func (t *TransactionManager) supportSavePoint() bool {
	// @todo db check
	return true
}

func (t *TransactionManager) max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: NewDbManager().GetDb(),
	}
}
