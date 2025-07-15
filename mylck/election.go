// elect node master, user mysql

package mylck

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	isMasterNode = false
)

// Election - election
type Election struct {
	NodeName      string
	MutexLockDao  *MutexLockDao
	MasterKey     string
	ElectInterval string
	Timeout       time.Duration
}

// IsMaster - is master node
func IsMaster() bool {
	return isMasterNode
}

// WrapJob - schedule election job
func (n *Election) WrapJob() {
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	_, err := c.AddFunc(n.ElectInterval, func() {
		n.elect()
		slog.Debug("[Schedule] [NodeElection] Done", "isMaster", IsMaster())
	})
	if err != nil {
		slog.Error("node manager schedule create error", "error", err)
		return
	}
	c.Start()
}

// elect - elect master
func (n *Election) elect() {
	n.MutexLockDao.BeginTx(nil)
	defer func() {
		_ = n.MutexLockDao.RollbackTx()
	}()

	mut, err := n.MutexLockDao.GetByName(n.MasterKey)
	if err != nil {
		slog.Error("get mutex lock error", "masterKey", n.MasterKey, "error", err)
		return
	}

	if mut == nil || mut.ID < 1 {
		err = n.MutexLockDao.Save(&MutexLock{
			Name:      n.MasterKey,
			Ver:       1,
			Info:      n.NodeName,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("create mutex lock error", "nodeName", n.NodeName, "error", err)
			return
		}

		err = n.MutexLockDao.CommitTx()
		if err != nil {
			slog.Error("mutex tx commit error", "nodeName", n.NodeName, "error", err)
			return
		}

		isMasterNode = true
		return
	}
	var rows int64
	// lock time out
	if time.Since(mut.UpdatedAt) > n.Timeout {
		rows, err = n.MutexLockDao.UpdateByInfo(n.MasterKey, mut.Info, &MutexLock{
			Info:      n.NodeName,
			UpdatedAt: time.Now(),
		})
	} else {
		// update by self
		rows, err = n.MutexLockDao.UpdateByInfo(n.MasterKey, n.NodeName, &MutexLock{
			Info:      n.NodeName,
			UpdatedAt: time.Now(),
		})
	}
	if err != nil {
		slog.Error("update mutex lock error", "nodeName", n.NodeName, "error", err)
		return
	}

	// update fail is slave node
	if rows < 1 {
		isMasterNode = false
		return
	}

	err = n.MutexLockDao.CommitTx()
	if err != nil {
		slog.Error("mutex tx commit error", "nodeName", n.NodeName, "error", err)
		return
	}
	isMasterNode = true
}
