package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
)

const (
	motesToCSPRRate = 1_000_000_000
	getRatesRetries = 10
)

type SyncRates struct {
	di.CasperClientAware

	networkName string
	httpClient  http.Client

	setRateDeployerPrivateKey     casper.PrivateKey
	csprRatesProviderContractHash casper.ContractHash
	contractExecutionAmount       int64
	rateAPIUrl                    string
	syncDuration                  time.Duration
}

func NewSyncRates() *SyncRates {
	return &SyncRates{
		httpClient: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *SyncRates) SetNetworkName(networkName string) {
	s.networkName = networkName
}

func (s *SyncRates) SetRateDeployerPrivateKey(privateKey casper.PrivateKey) {
	s.setRateDeployerPrivateKey = privateKey
}

func (s *SyncRates) SetRateAPIUrl(url string) {
	s.rateAPIUrl = url
}

func (s *SyncRates) SetCSPRRatesProviderContractHash(contractHash casper.ContractHash) {
	s.csprRatesProviderContractHash = contractHash
}

func (s *SyncRates) SetContractExecutionAmount(executionAmount int64) {
	s.contractExecutionAmount = executionAmount
}

func (s *SyncRates) SetSyncDuration(duration time.Duration) {
	s.syncDuration = duration
}

func (s *SyncRates) Execute(ctx context.Context) error {
	zap.S().Info("Rates oracle started...")
	if err := s.syncRates(ctx); err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to sync rates")
		return err
	}

	ticker := time.NewTicker(s.syncDuration)
	for {
		select {
		case <-ticker.C:
			if err := s.syncRates(ctx); err != nil {
				zap.S().With(zap.Error(err)).Error("Failed to sync rates")
				return err
			}
		case <-ctx.Done():
			zap.S().Info("Exit on context signal")
			return ctx.Err()
		}
	}
}

func (s *SyncRates) syncRates(ctx context.Context) error {
	var (
		rates float32
		err   error
	)
	for i := 1; i <= getRatesRetries; i++ {
		rates, err = s.getRates(ctx)
		if err == nil {
			zap.S().Info("Got rates successfully")
			break
		}
		zap.S().With(zap.Error(err)).With(zap.String("url", s.rateAPIUrl)).Error("Failed to get rates, retrying...")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	if rates == 0 {
		zap.S().With(zap.Error(err)).Error("Failed to get rates")
		return fmt.Errorf("unable to get rates from the URL - %s", s.rateAPIUrl)
	}

	if err := s.putSetRateDeploy(ctx, rates); err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to get rates")
		return err
	}

	zap.S().Info("Rates successfully synchronized")
	return nil
}

func (s *SyncRates) getRates(ctx context.Context) (float32, error) {
	req, err := http.NewRequest(http.MethodGet, s.rateAPIUrl, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("request is not successful - response code: %d", resp.StatusCode)
	}

	var rateResponse struct {
		Data float32 `json:"data"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&rateResponse); err != nil {
		return 0, err
	}
	return rateResponse.Data, nil
}

func (s *SyncRates) putSetRateDeploy(ctx context.Context, rates float32) error {
	rate := int64(float32(1/rates) * motesToCSPRRate)

	args := (&casper.Args{}).
		AddArgument("rate", *clvalue.NewCLUInt512(big.NewInt(rate)))

	session := casper.ExecutableDeployItem{
		StoredContractByHash: &casper.StoredContractByHash{
			Hash:       s.csprRatesProviderContractHash,
			EntryPoint: "set_rate",
			Args:       args,
		},
	}

	deployHeader := casper.DefaultHeader()
	deployHeader.Account = s.setRateDeployerPrivateKey.PublicKey()
	deployHeader.ChainName = s.networkName

	payment := casper.StandardPayment(big.NewInt(s.contractExecutionAmount))

	deploy, err := casper.MakeDeploy(deployHeader, payment, session)
	if err != nil {
		return err
	}

	err = deploy.SignDeploy(s.setRateDeployerPrivateKey)
	if err != nil {
		return err
	}

	_, err = s.GetCasperClient().PutDeploy(ctx, *deploy)
	return nil
}
