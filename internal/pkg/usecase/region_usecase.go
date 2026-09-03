package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"evermos-api/internal/pkg/model"
)

type RegionUseCase interface {
	GetProvinces(ctx context.Context) ([]model.Province, error)
	GetRegencies(ctx context.Context, provinceID string) ([]model.Regency, error)
	GetDistricts(ctx context.Context, regencyID string) ([]model.District, error)
	GetVillages(ctx context.Context, districtID string) ([]model.Village, error)
}

type regionUseCase struct {
	baseURL string
	client  *http.Client
	mu      sync.RWMutex
	cache   map[string]interface{}
}

func NewRegionUseCase() RegionUseCase {
	return &regionUseCase{
		baseURL: "https://emsifa.github.io/api-wilayah-indonesia/api",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		cache: make(map[string]interface{}),
	}
}

func (u *regionUseCase) GetProvinces(ctx context.Context) ([]model.Province, error) {
	cacheKey := "provinces"
	u.mu.RLock()
	if val, ok := u.cache[cacheKey]; ok {
		u.mu.RUnlock()
		return val.([]model.Province), nil
	}
	u.mu.RUnlock()

	url := fmt.Sprintf("%s/provinces.json", u.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var provinces []model.Province
	if err := json.NewDecoder(resp.Body).Decode(&provinces); err != nil {
		return nil, fmt.Errorf("failed to decode provinces response: %w", err)
	}

	u.mu.Lock()
	u.cache[cacheKey] = provinces
	u.mu.Unlock()

	return provinces, nil
}

func (u *regionUseCase) GetRegencies(ctx context.Context, provinceID string) ([]model.Regency, error) {
	cacheKey := fmt.Sprintf("regencies_%s", provinceID)
	u.mu.RLock()
	if val, ok := u.cache[cacheKey]; ok {
		u.mu.RUnlock()
		return val.([]model.Regency), nil
	}
	u.mu.RUnlock()

	url := fmt.Sprintf("%s/regencies/%s.json", u.baseURL, provinceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch regencies from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var regencies []model.Regency
	if err := json.NewDecoder(resp.Body).Decode(&regencies); err != nil {
		return nil, fmt.Errorf("failed to decode regencies response: %w", err)
	}

	u.mu.Lock()
	u.cache[cacheKey] = regencies
	u.mu.Unlock()

	return regencies, nil
}

func (u *regionUseCase) GetDistricts(ctx context.Context, regencyID string) ([]model.District, error) {
	cacheKey := fmt.Sprintf("districts_%s", regencyID)
	u.mu.RLock()
	if val, ok := u.cache[cacheKey]; ok {
		u.mu.RUnlock()
		return val.([]model.District), nil
	}
	u.mu.RUnlock()

	url := fmt.Sprintf("%s/districts/%s.json", u.baseURL, regencyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch districts from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var districts []model.District
	if err := json.NewDecoder(resp.Body).Decode(&districts); err != nil {
		return nil, fmt.Errorf("failed to decode districts response: %w", err)
	}

	u.mu.Lock()
	u.cache[cacheKey] = districts
	u.mu.Unlock()

	return districts, nil
}

func (u *regionUseCase) GetVillages(ctx context.Context, districtID string) ([]model.Village, error) {
	cacheKey := fmt.Sprintf("villages_%s", districtID)
	u.mu.RLock()
	if val, ok := u.cache[cacheKey]; ok {
		u.mu.RUnlock()
		return val.([]model.Village), nil
	}
	u.mu.RUnlock()

	url := fmt.Sprintf("%s/villages/%s.json", u.baseURL, districtID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch villages from external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var villages []model.Village
	if err := json.NewDecoder(resp.Body).Decode(&villages); err != nil {
		return nil, fmt.Errorf("failed to decode villages response: %w", err)
	}

	u.mu.Lock()
	u.cache[cacheKey] = villages
	u.mu.Unlock()

	return villages, nil
}
