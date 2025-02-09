package usecase

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PingUsecase interface {
	SendPingData(ip string, pingTime int) error
	GetPingInterval() time.Duration
}

type pingUsecase struct {
	backendURL   string
	pingInterval time.Duration
}

func NewPingUsecase(backendURL string, pingInterval time.Duration) PingUsecase {
	return &pingUsecase{backendURL: backendURL, pingInterval: pingInterval}
}

func (uc *pingUsecase) SendPingData(ip string, pingTime int) error {
	endpoint := uc.backendURL + "/pings/add"
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("ip_address", ip)
	q.Set("ping_time", strconv.Itoa(pingTime))
	u.RawQuery = q.Encode()

	resp, err := http.Post(u.String(), "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("backend вернул статус: %s", resp.Status)
	}

	return nil
}

func (uc *pingUsecase) GetPingInterval() time.Duration {
	return uc.pingInterval
}
