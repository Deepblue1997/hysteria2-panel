package services

import (
	"errors"
	"hysteria2-panel/models"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

type TrafficService struct {
	db    *gorm.DB
	mutex sync.RWMutex
	// 用户实时流量统计，key 为用户ID
	stats map[uint]*TrafficStat
}

type TrafficStat struct {
	Upload   int64
	Download int64
	LastSync time.Time
}

func NewTrafficService(db *gorm.DB) *TrafficService {
	service := &TrafficService{
		db:    db,
		stats: make(map[uint]*TrafficStat),
	}
	go service.syncTrafficPeriodically()
	return service
}

// 记录流量使用
func (s *TrafficService) RecordTraffic(userID uint, upload, download int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stat, exists := s.stats[userID]
	if !exists {
		stat = &TrafficStat{LastSync: time.Now()}
		s.stats[userID] = stat
	}

	stat.Upload += upload
	stat.Download += download

	// 检查是否需要同步到数据库
	if time.Since(stat.LastSync) > 5*time.Minute {
		if err := s.syncUserTraffic(userID); err != nil {
			return err
		}
	}

	return nil
}

// 检查用户是否超出流量限制
func (s *TrafficService) CheckTrafficLimit(userID string) (bool, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return false, errors.New("无效的用户ID")
	}

	var user models.User
	if err := s.db.First(&user, uid).Error; err != nil {
		return false, err
	}

	// 检查账户是否过期
	if !user.ExpireAt.IsZero() && time.Now().After(user.ExpireAt) {
		return false, errors.New("账户已过期")
	}

	// 获取实时流量统计
	s.mutex.RLock()
	stat, exists := s.stats[uint(uid)]
	s.mutex.RUnlock()

	if exists {
		currentTraffic := user.Traffic + stat.Upload + stat.Download
		// 如果流量限制为0，表示不限制
		if user.TrafficLimit > 0 && currentTraffic >= user.TrafficLimit {
			return false, errors.New("已超出流量限制")
		}
	}

	return true, nil
}

// 定期同步流量数据到数据库
func (s *TrafficService) syncTrafficPeriodically() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mutex.Lock()
		for userID := range s.stats {
			_ = s.syncUserTraffic(userID)
		}
		s.mutex.Unlock()
	}
}

// 同步单个用户的流量数据到数据库
func (s *TrafficService) syncUserTraffic(userID uint) error {
	stat := s.stats[userID]
	if stat.Upload == 0 && stat.Download == 0 {
		return nil
	}

	err := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("traffic", gorm.Expr("traffic + ?", stat.Upload+stat.Download)).
		Error

	if err == nil {
		stat.Upload = 0
		stat.Download = 0
		stat.LastSync = time.Now()
	}

	return err
}
