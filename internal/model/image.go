package model

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
	ID        uuid.UUID `json:"id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"file_path"`
	Action    Action    `json:"actions"` // action to perform
	Status    string    `json:"status"`  // pending / processed / failed
	CreatedAt time.Time `json:"created_at"`
}

type Action struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}
