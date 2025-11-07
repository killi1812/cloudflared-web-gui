package model

import (
	"time"

	"github.com/google/uuid"
)

type Tunnel struct {
	Id          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Connections []Connection `json:"connections"`
	Token       string       `json:"token,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Connection struct {
	// "colo_name": "zag01",
	//
	//	"id": "8877da37-d646-45c1-9e7e-2de1e14711d7",
	//	"is_pending_reconnect": false,
	//	"origin_ip": "141.136.241.49",
	//	"opened_at": "2025-11-06T08:29:57.963765Z"
}

// {"level":"warn","message":"Your version 2025.9.0 is outdated. We recommend upgrading it to 2025.10.1","time":"2025-11-06T12:22:04Z"}
