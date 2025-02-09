package repository

import (
	"database/sql"
	"fmt"

	"backend/internal/domain"
)

type PostgresPingRepository struct {
	db *sql.DB
}

func NewPostgresPingRepository(db *sql.DB) *PostgresPingRepository {
	return &PostgresPingRepository{db: db}
}

func (r *PostgresPingRepository) GetAllPings() ([]domain.PingData, error) {
	rows, err := r.db.Query("SELECT ip_address, ping_time, last_successful_ping FROM pings")
	if err != nil {
		return nil, fmt.Errorf("failed to query pings: %w", err)
	}
	defer rows.Close()

	var pings []domain.PingData
	for rows.Next() {
		var ping domain.PingData
		err := rows.Scan(&ping.IPAddress, &ping.PingTime, &ping.LastSuccessfulPing)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ping %w", err)
		}
		pings = append(pings, ping)
	}

	return pings, nil
}

func (r *PostgresPingRepository) AddPing(ping domain.PingData) error {
	_, err := r.db.Exec("INSERT INTO pings (ip_address, ping_time, last_successful_ping) VALUES ($1, $2, $3) ON CONFLICT (ip_address) DO UPDATE SET ping_time = $2, last_successful_ping = $3", ping.IPAddress, ping.PingTime, ping.LastSuccessfulPing)
	if err != nil {
		return fmt.Errorf("failed to insert ping %w", err)
	}
	return nil
}
