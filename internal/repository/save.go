package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
)

func (r *Repository) Save(ctx context.Context, img model.Image) (uuid.UUID, error) {
	query := `
		INSERT INTO images (filename, path, action, params, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
   `

	paramsJSON, err := json.Marshal(img.Action.Params)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal action params: %w", err)
	}

	var id uuid.UUID
	err = r.store.DB.QueryRowContext(
		ctx, query, img.Filename, img.Path, img.Action.Name, paramsJSON, img.Status,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to save in db: %w", err)
	}

	return id, nil

}
