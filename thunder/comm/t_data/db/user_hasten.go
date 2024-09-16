package db

import (
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"sync"
	"time"
)

type HastenUnit struct {
	StartTime time.Time `bson:"start_time" json:"start_time,omitempty"`
	EndTime   time.Time `bson:"end_time" json:"end_time,omitempty"`
}
type Hasten struct {
	Data    map[pb.HastenType]*HastenUnit `bson:"data" json:"data"`
	EndTime time.Time                     `bson:"end_time" json:"end_time,omitempty"`
}

func NewHasten() *Hasten {
	m := make(map[pb.HastenType]*HastenUnit)
	for k := range pb.HastenType_name {
		m[pb.HastenType(k)] = &HastenUnit{
			StartTime: time.Time{},
			EndTime:   time.Time{},
		}
	}

	return &Hasten{Data: m}
}

var hastenLock sync.RWMutex

func (s *Hasten) AddEndTime(d time.Duration) {
	timeNow := time.Now()
	if !s.EndTime.IsZero() && s.EndTime.After(timeNow) {
		s.EndTime = s.EndTime.Add(d)
	} else {
		s.EndTime = timeNow.Add(d)
	}
}

func (s *Hasten) IsMax() bool {
	if !s.EndTime.IsZero() && s.EndTime.After(time.Now().Add(time.Second*time.Duration(t_tdb.Conf().Hasten_max_second))) {
		return true
	}
	return false
}

func (s *Hasten) GetMultiply() float64 {
	multiply := 1.0
	if s.EndTime.After(time.Now()) {
		multiply = 2.0
	}
	return multiply
}

func (s *Hasten) Has(t pb.HastenType) bool {
	hastenLock.RLock()
	defer hastenLock.RUnlock()

	if _, ok := s.Data[t]; ok {
		return true
	}

	return false
}

func (s *Hasten) Reset() {
	hastenLock.Lock()
	defer hastenLock.Unlock()

	for k := range pb.HastenType_name {
		s.Data[pb.HastenType(k)] = &HastenUnit{
			StartTime: time.Time{},
			EndTime:   time.Time{},
		}
	}
	s.EndTime = time.Time{}
}

func (s *Hasten) Set(t pb.HastenType, unit *HastenUnit) {
	hastenLock.Lock()
	defer hastenLock.Unlock()

	s.Data[t] = unit
}

func (s *Hasten) Get(t pb.HastenType) *HastenUnit {
	hastenLock.Lock()
	defer hastenLock.Unlock()

	return s.Data[t]
}

func (s *Hasten) ToPb() *pb.Hasten {
	hastenLock.RLock()
	defer hastenLock.RUnlock()

	return &pb.Hasten{
		Hasten:  s.ToPbSlice(),
		EndTime: s.EndTime.Unix(),
	}
}

func (s *Hasten) ToPbSlice() []*pb.HastenUnit {
	hastenLock.RLock()
	defer hastenLock.RUnlock()

	pbData := make([]*pb.HastenUnit, len(s.Data))
	num := 0
	for t := range s.Data {
		pbData[num] = s.ToPbByType(t)
		num++
	}

	return pbData
}

func (s *Hasten) ToPbByType(t pb.HastenType) *pb.HastenUnit {
	return &pb.HastenUnit{
		HastenType: t,
		StartTime:  s.Data[t].StartTime.Unix(),
		EndTime:    s.Data[t].EndTime.Unix(),
	}
}
