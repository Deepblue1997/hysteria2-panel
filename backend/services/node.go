package services

import (
	"errors"
	"hysteria2-panel/models"
	"time"

	"gorm.io/gorm"
)

type NodeService struct {
	db *gorm.DB
}

func NewNodeService(db *gorm.DB) *NodeService {
	return &NodeService{db: db}
}

// 创建节点
func (s *NodeService) CreateNode(node *models.Node) error {
	return s.db.Create(node).Error
}

// 获取节点列表
func (s *NodeService) GetNodes() ([]models.Node, error) {
	var nodes []models.Node
	err := s.db.Find(&nodes).Error
	return nodes, err
}

// 更新节点状态
func (s *NodeService) UpdateNodeStatus(nodeID uint, status *models.NodeStatus) error {
	updates := map[string]interface{}{
		"status":    1, // 在线
		"last_ping": time.Now(),
	}

	result := s.db.Model(&models.Node{}).Where("id = ?", nodeID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("节点不存在")
	}

	// 保存节点状态到缓存或其他存储中
	// TODO: 实现状态存储逻辑

	return nil
}

// 获取节点状态
func (s *NodeService) GetNodeStatus(nodeID uint) (*models.NodeStatus, error) {
	var node models.Node
	if err := s.db.First(&node, nodeID).Error; err != nil {
		return nil, err
	}

	// TODO: 从缓存或其他存储中获取节点状态
	status := &models.NodeStatus{
		LastReportAt: node.LastPing.Unix(),
	}

	return status, nil
}

// 检查节点在线状态
func (s *NodeService) CheckNodesStatus() error {
	// 将超过5分钟未更新的节点标记为离线
	timeout := time.Now().Add(-5 * time.Minute)
	return s.db.Model(&models.Node{}).
		Where("last_ping < ? AND status = ?", timeout, 1).
		Update("status", 0).Error
}
