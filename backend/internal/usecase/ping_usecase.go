package usecase

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"time"
)

type PingUsecase interface {
	GetAllPings() ([]domain.PingData, error)
	AddPing(ping domain.PingData) error
}

type pingUsecase struct {
	pingRepo     repository.PingRepository
	pingInterval time.Duration
}

func NewPingUsecase(pingRepo repository.PingRepository, pingInterval time.Duration) PingUsecase {
	return &pingUsecase{pingRepo: pingRepo, pingInterval: pingInterval}
}

func (uc *pingUsecase) GetAllPings() ([]domain.PingData, error) {
	return uc.pingRepo.GetAllPings()
}

func (uc *pingUsecase) AddPing(ping domain.PingData) error {
	return uc.pingRepo.AddPing(ping)
}
