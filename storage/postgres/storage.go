package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

func (r *postgresRepo) StorageCreate(ctx context.Context, req *models.StorageCreateReq) (*models.StorageResponse, error) {
	res := &models.StorageResponse{}
	query := r.Db.Builder.Insert("storages").Columns(
		"id",
		"type",
		"domain_name",
		"access_key",
		"secret_key",
		"region",
	).Values(req.ID, req.Type, req.DomainName, req.AccessKey, req.SecretKey, req.Region).Suffix(
		"RETURNING id, type, domain_name, access_key, secret_key, region, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Type,
		&res.DomainName, &res.AccessKey,
		&res.SecretKey, &res.Region, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) StorageGet(ctx context.Context, req *models.StorageGetReq) (*models.StorageResponse, error) {
	query := r.Db.Builder.Select("id, type, domain_name, access_key, secret_key, region, created_at, updated_at").
		From("storages")

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"id": req.ID})
	} else {
		return &models.StorageResponse{}, fmt.Errorf("at least one filter should be exists")
	}
	res := &models.StorageResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Type,
		&res.DomainName, &res.AccessKey,
		&res.SecretKey, &res.Region, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) StorageFind(ctx context.Context, req *models.StorageFindReq) (*models.StorageFindResponse, error) {
	var (
		res            = &models.StorageFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"domain_name": req.Search + "%"})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("storages").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, type, domain_name, access_key, secret_key, region, created_at, updated_at").
		From("storages").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.StorageResponse{}
		err := rows.Scan(
			&temp.ID, &temp.Type,
			&temp.DomainName, &temp.AccessKey,
			&temp.SecretKey, &temp.Region, &CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "CompanyFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)

		res.Storages = append(res.Storages, temp)
	}

	return res, nil
}

func (r *postgresRepo) StorageUpdate(ctx context.Context, req *models.StorageUpdateReq) (*models.StorageResponse, error) {
	mp := make(map[string]interface{})
	mp["type"] = req.Type
	mp["domain_name"] = req.DomainName
	mp["access_key"] = req.AccessKey
	mp["secret_key"] = req.SecretKey
	mp["region"] = req.Region
	mp["updated_at"] = time.Now()
	query := r.Db.Builder.Update("storages").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix("RETURNING id, type, domain_name, access_key, secret_key, created_at, updated_at")

	res := &models.StorageResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Type,
		&res.DomainName, &res.AccessKey,
		&res.SecretKey, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) StorageDelete(ctx context.Context, req *models.StorageDeleteReq) error {
	query := r.Db.Builder.Delete("storages").Where(squirrel.Eq{"id": req.ID})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "CompanyDelete: query.RunWith(r.Db.Db).Exec()")
}
