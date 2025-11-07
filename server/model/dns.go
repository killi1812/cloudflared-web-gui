package model

import "time"

type DnsRecord struct {
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

	CreatedAt  time.Time `json:"created_at"`
	ModifiedOn time.Time `json:"modified_on"`

	/*
	   "id": "6add6cf92fb83351b8ff32efe67a9ef5",
	   	"name": "cmd.francvok.from.hr",
	   	"type": "CNAME",
	   	"content": "6fe4ac0c-4d13-499e-b031-31065f16b611.cfargotunnel.com",
	   	"proxiable": true,
	   	"proxied": true,
	   	"ttl": 1,
	   	"settings": {
	   	  "flatten_cname": false
	   	},
	   	"meta": {},
	   	"comment": null,
	   	"tags": [],
	   	"created_on": "2025-11-05T10:18:11.403687Z",
	   	"modified_on": "2025-11-05T10:18:11.403687Z"
	*/
}
