package db

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/enums"
	"sync"
)

type RedPoint struct {
	Data map[pb.RedPointType]int `bson:"data" json:"data"`
}

func NewRedPoint() *RedPoint {
	m := make(map[pb.RedPointType]int)
	for k := range pb.RedPointType_name {
		if k == 0 {
			continue
		}
		m[pb.RedPointType(k)] = enums.YES
	}

	return &RedPoint{Data: m}
}

var redPointLock sync.RWMutex

func (s *RedPoint) Has(t pb.RedPointType) bool {
	redPointLock.RLock()
	defer redPointLock.RUnlock()

	if _, ok := s.Data[t]; ok {
		return true
	}

	return false
}

func (s *RedPoint) Reset() {
	redPointLock.Lock()
	defer redPointLock.Unlock()

	for k := range pb.RedPointType_name {
		if k == 0 {
			continue
		}
		s.Data[pb.RedPointType(k)] = enums.YES
	}
}

func (s *RedPoint) Look(t pb.RedPointType) {
	redPointLock.Lock()
	defer redPointLock.Unlock()

	s.Data[t] = enums.YES
}

func (s *RedPoint) UnLook(t pb.RedPointType) {
	redPointLock.Lock()
	defer redPointLock.Unlock()

	s.Data[t] = enums.NO
}

func (s *RedPoint) ToPb() []*pb.RedPointUnit {
	redPointLock.RLock()
	defer redPointLock.RUnlock()

	pbData := make([]*pb.RedPointUnit, len(s.Data))
	num := 0
	for t := range s.Data {
		pbData[num] = s.ToPbByType(t)
		num++
	}

	return pbData
}

func (s *RedPoint) ToPbByType(t pb.RedPointType) *pb.RedPointUnit {
	return &pb.RedPointUnit{
		RedPointType: t,
		Flag:         s.Data[t] == enums.YES,
	}
}
