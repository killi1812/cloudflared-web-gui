package dto

import (
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/format"
)

type TunnelDto struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	DnsRecords ArrDnsRecordDto `json:"dnsRecords"`

	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
}

func (t *TunnelDto) FromModel(tnl model.Tunnel) {
	t.Id = tnl.Id.String()
	t.Name = tnl.Name
	t.CreatedAt = tnl.CreatedAt.Format(format.DateTimeFormat)
	t.DeletedAt = tnl.DeletedAt.Format(format.DateTimeFormat)
}

type DnsRecordDto struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Content   string  `json:"content"`
	Proxiable bool    `json:"proxiable"`
	Proxied   bool    `json:"proxied"`
	Ttl       int     `json:"ttl"`
	Settings  any     `json:"settings"`
	Meta      any     `json:"meta"`
	Commnet   *string `json:"commnet"`
	Tags      []any   `json:"tags"`

	CreatedAt  string `json:"created_at"`
	ModifiedOn string `json:"modified_on"`
}

func (d *DnsRecordDto) FromModel(dns model.DnsRecord) {
	d.Id = dns.Id
	d.Name = dns.Name
	d.Type = dns.Type
	d.Content = dns.Content
	d.Proxiable = dns.Proxiable
	d.Proxied = dns.Proxied
	d.Ttl = dns.Ttl
	d.Settings = dns.Settings
	d.Meta = dns.Meta
	if dns.Commnet != nil {
		d.Commnet = dns.Commnet
	}
	d.Tags = dns.Tags
	d.CreatedAt = dns.CreatedAt.Format(format.DateTimeFormat)
	d.ModifiedOn = dns.ModifiedOn.Format(format.DateTimeFormat)
}

type ArrDnsRecordDto []DnsRecordDto

func (a *ArrDnsRecordDto) FromModel(dnsRecords []model.DnsRecord) {
	tmp := make(ArrDnsRecordDto, len(dnsRecords))
	for i, dns := range dnsRecords {
		tmp[i].FromModel(dns)
	}
	*a = tmp
}
