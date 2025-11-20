package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

func (r *postgresRepo) WebhookCreate(ctx context.Context, req *models.WebhookCreateReq) (*models.WebhookResponse, error) {
	res := &models.WebhookResponse{}
	query := r.Db.Builder.Insert("webhooks").Columns(
		"id",
		"project_id",
		"title",
		"webhook_type",
		"url",
		"active",
	).Values(req.ID, req.ProjectID, req.Title, req.WebhookType, req.URL, req.Active).Suffix(
		"RETURNING id, project_id, title, webhook_type, url, active,created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.ProjectID,
		&res.Title, &res.WebhookType,
		&res.URL, &res.Active,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "WebhookCreate: query.RunWith(r.Db.Db).Scan()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) WebhookGet(ctx context.Context, req *models.WebhookGetReq) (*models.WebhookResponse, error) {
	query := r.Db.Builder.Select("id, project_id, title, webhook_type, url, active, created_at, updated_at").
		From("webhooks")

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"id": req.ID})
	} else {
		return &models.WebhookResponse{}, fmt.Errorf("at least one filter should be exists")
	}
	res := &models.WebhookResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.ProjectID,
		&res.Title, &res.WebhookType,
		&res.URL, &res.Active,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "WebhookCreate: query.RunWith(r.Db.Db).Scan()")
	}
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "WebhookGet:query.RunWith(r.Db.Db).QueryRow()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) WebhookFind(ctx context.Context, req *models.WebhookFindReq) (*models.WebhookFindResponse, error) {
	var (
		res            = &models.WebhookFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"webhook_name": req.Search + "%"})
	}
	if strings.TrimSpace(req.ProjectId) != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"project_id": req.ProjectId})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("webhooks").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "WebhookFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, project_id, title, webhook_type, url, active, created_at, updated_at").
		From("webhooks").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "WebhookFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.WebhookResponse{}
		err := rows.Scan(
			&temp.ID, &temp.ProjectID,
			&temp.Title, &temp.WebhookType,
			&temp.URL, &temp.Active,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "WebhookFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Webhooks = append(res.Webhooks, temp)
	}

	return res, nil
}

func (r *postgresRepo) WebhookUpdate(ctx context.Context, req *models.WebhookUpdateReq) (*models.WebhookResponse, error) {
	mp := make(map[string]interface{})
	mp["title"] = req.Title
	mp["webhook_type"] = req.WebhookType
	mp["url"] = req.URL
	mp["active"] = req.Active
	mp["updated_at"] = time.Now()
	query := r.Db.Builder.Update("webhooks").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix("RETURNING id, project_id, title, webhook_type, url, active, created_at, updated_at")

	res := &models.WebhookResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.ProjectID,
		&res.Title, &res.WebhookType,
		&res.URL, &res.Active,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "WebhookUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) WebhookDelete(ctx context.Context, req *models.WebhookDeleteReq) error {
	query := r.Db.Builder.Delete("webhooks").Where(squirrel.Eq{"id": req.ID})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "WebhookDelete: query.RunWith(r.Db.Db).Exec()")
}
