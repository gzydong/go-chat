package service

import (
	"context"

	"go-chat/internal/repository/model"
)

type ContactGroupService struct {
	*BaseService
}

func NewContactGroupService(base *BaseService) *ContactGroupService {
	return &ContactGroupService{BaseService: base}
}

func (c *ContactGroupService) Create(ctx context.Context, value *model.ContactGroup) (int, error) {
	return 0, nil
}
