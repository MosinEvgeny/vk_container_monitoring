package repository

import "backend/internal/domain"

type PingRepository interface {
	GetAllPings() ([]domain.PingData, error)
	AddPing(ping domain.PingData) error
}
