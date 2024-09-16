package db

import (
	"github.com/v587-zyf/gc/enums"
	"sync"
)

type Invite struct {
	Invite   uint64           `bson:"invite" json:"invite,omitempty"` // 邀请人id
	Flag     int              `bson:"flag" json:"flag,omitempty"`     // 是否发送奖励
	Invitees map[uint64]int   `bson:"invitees" json:"invitees"`       // UID:diamond
	Reward   map[int]struct{} `bson:"reward" json:"reward"`           // 已领取任务奖励表id
}

func NewInvite(invite uint64) *Invite {
	return &Invite{
		Invite:   invite,
		Flag:     enums.NO,
		Invitees: make(map[uint64]int),
		Reward:   make(map[int]struct{}),
	}
}

var inviteLock sync.RWMutex

func (i *Invite) Reset() {
	inviteLock.Lock()
	defer inviteLock.Unlock()

	i.Reward = make(map[int]struct{})
}

func (i *Invite) AddInvitees(uID uint64, diamond int) {
	inviteLock.Lock()
	defer inviteLock.Unlock()

	i.Invitees[uID] += diamond
}

func (i *Invite) AddReward(id int) {
	inviteLock.Lock()
	defer inviteLock.Unlock()

	i.Reward[id] = struct{}{}
}

//
//func (i *Invite) ToPbInvitees(sid int64) []*pb.InviteUnit {
//	inviteLock.RLock()
//	defer inviteLock.RUnlock()
//
//	var (
//		num    = 0
//		pbData = make([]*pb.InviteUnit, len(i.Invitees))
//
//		err     error
//		rdbUser *redis.User
//		locker  *rdb_cluster.Locker
//	)
//	for userID, diamond := range i.Invitees {
//		locker = redis.LockUser(userID, sid)
//		if locker == nil {
//			log.Error("get locker err", zap.Uint64("userID", userID))
//			continue
//		}
//		var unlock = func() { locker.Unlock() }
//		rdbUser, err = redis.GetUser(userID, locker)
//		if err != nil {
//			unlock()
//			log.Error("get redis user err", zap.Error(err), zap.Uint64("userID", userID))
//			continue
//		}
//		pbData[num] = &pb.InviteUnit{
//			Uid:       userID,
//			Head:      rdbUser.Telegram.Head,
//			FirstName: rdbUser.Telegram.FirstName,
//			LastName:  rdbUser.Telegram.LastName,
//			UserName:  rdbUser.Telegram.UserName,
//			Lv:        int32(rdbUser.Basic.Lv),
//			Diamond:   int32(diamond),
//		}
//		num++
//	}
//	return pbData
//}

func (i *Invite) ToPbReward() []int32 {
	inviteLock.RLock()
	defer inviteLock.RUnlock()

	pbData := make([]int32, len(i.Reward))
	num := 0
	for k := range i.Reward {
		pbData[num] = int32(k)
		num++
	}

	return pbData
}
