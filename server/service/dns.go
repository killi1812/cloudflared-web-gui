package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/killi1812/cloudflared-web-gui/app"
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/cerror"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IDnsSrv interface {
	GetDnsRecords(uuid uuid.UUID) ([]model.DnsRecord, error)
}

func NewDnsSrv() IDnsSrv {
	var service IDnsSrv
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &DnsSrv{
			db:     db,
			logger: logger,
		}
	})

	return service
}

type respT struct {
	Result     []model.DnsRecord `json:"result"`
	Success    bool              `json:"success"`
	ResultInfo struct {
		Page        int `json:"page"`
		PerPage     int `json:"per_page"`
		Count       int `json:"count"`
		Total_count int `json:"total_count"`
		Total_pages int `json:"total_pages"`
	} `json:"result_info"`
}

type DnsSrv struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

// GetDnsRecords implements IDnsSrv.
func (d *DnsSrv) GetDnsRecords(uuid uuid.UUID) ([]model.DnsRecord, error) {
	if app.ZoneId == "" {
		d.logger.Error(cerror.ErrZoneIdNotSet)
		return nil, cerror.ErrZoneIdNotSet
	}

	if app.CloudflaredApiKey == "" {
		d.logger.Error(cerror.ErrCloudflaredApiKeyNotSet)
		return nil, cerror.ErrCloudflaredApiKeyNotSet
	}

	baseURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", app.ZoneId)

	tunnelContent := uuid.String() + ".cfargotunnel.com"
	// 2. Create a new HTTP request
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		d.logger.Errorf("Error creating request, err = %v", err)
		return nil, err
	}

	// 3. Add query parameters to the request URL
	q := req.URL.Query()
	q.Add("type", "CNAME")
	q.Add("content", tunnelContent)
	req.URL.RawQuery = q.Encode()

	// 4. Set the required headers
	req.Header.Set("Authorization", "Bearer "+app.CloudflaredApiKey)
	req.Header.Set("Content-Type", "application/json")

	// 5. Create a client and execute the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		d.logger.Errorf("Error sending request, err = %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 6. Read and print the response
	d.logger.Debugf("Cloudflared response Status Code: %d", resp.StatusCode)

	var res respT
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		d.logger.Errorf("Error reading response body, err = %v", err)
		return nil, err
	}

	d.logger.Debugf("Cloudflared response = %+v", res)

	return res.Result, nil
}
