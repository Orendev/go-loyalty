package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Orendev/go-loyalty/internal/config"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrorGetAccrual = errors.New("get in accrual system was invalid")
)

type HTTPClient struct {
	repo          repository.Storage
	accrualSystem config.AccrualSystem
	accrualChan   chan models.Accrual
}

func NewHTTPClient(ctx context.Context, repo repository.Storage, accrualSystem config.AccrualSystem, accrualChan chan models.Accrual) (*HTTPClient, error) {
	instance := &HTTPClient{repo: repo, accrualSystem: accrualSystem, accrualChan: accrualChan}

	go instance.worker(ctx)

	return instance, nil
}

func (h *HTTPClient) GetAccrual(order int) (models.AccrualResponse, error) {

	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%v", h.accrualSystem.Addr, order))
	var accrualResponse models.AccrualResponse
	if err != nil {
		return accrualResponse, err
	}

	switch resp.StatusCode {

	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return accrualResponse, err
		}

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				return
			}
		}()

		err = json.Unmarshal(body, &accrualResponse)
		if err != nil {
			return accrualResponse, err
		}

		return accrualResponse, nil

	case http.StatusTooManyRequests:
		header := resp.Header
		a := header["Retry-After"][0]
		sec, err := strconv.Atoi(a)
		if err != nil {
			return accrualResponse, err
		}
		accrualResponse.RetryAfterDuration = time.Duration(sec)

		return accrualResponse, nil //
	}

	return accrualResponse, ErrorGetAccrual
}
