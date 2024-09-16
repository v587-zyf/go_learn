package db

import (
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"sync"
)

type Shop struct {
	Data map[pb.ShopType]int `bson:"data" json:"data"`
}

func NewShop() *Shop {
	m := make(map[pb.ShopType]int)
	for k := range pb.ShopType_name {
		if k == 0 {
			continue
		}
		m[pb.ShopType(k)] = 0
	}

	return &Shop{Data: m}
}

var shopLock sync.RWMutex

func (s *Shop) Has(t pb.ShopType) bool {
	shopLock.RLock()
	defer shopLock.RUnlock()

	if _, ok := s.Data[t]; ok {
		return true
	}

	return false
}

func (s *Shop) Reset() {
	shopLock.Lock()
	defer shopLock.Unlock()

	for k := range pb.ShopType_name {
		if k == 0 {
			continue
		}
		s.Data[pb.ShopType(k)] = 0
	}
}

func (s *Shop) AddNum(t pb.ShopType) {
	shopLock.Lock()
	defer shopLock.Unlock()

	s.Data[t]++
}

func (s *Shop) SetNum(t pb.ShopType, num int) {
	shopLock.Lock()
	defer shopLock.Unlock()

	s.Data[t] = num
}

func (s *Shop) GetNum(t pb.ShopType) int {
	shopLock.Lock()
	defer shopLock.Unlock()

	return s.Data[t]
}

func (s *Shop) GetAllShopStrength() int {
	shopLock.RLock()
	defer shopLock.RUnlock()

	return s.Data[pb.ShopType_Extra_Strength] * t_tdb.Conf().Buy_extra_strength
}

func (s *Shop) ToPb() []*pb.ShopUnit {
	shopLock.RLock()
	defer shopLock.RUnlock()

	pbData := make([]*pb.ShopUnit, len(s.Data))
	num := 0
	for t := range s.Data {
		pbData[num] = s.ToPbByType(t)
		num++
	}

	return pbData
}

func (s *Shop) ToPbByType(t pb.ShopType) *pb.ShopUnit {
	return &pb.ShopUnit{
		ShopBuyType: t,
		BuyNum:      int32(s.Data[t]),
	}
}
