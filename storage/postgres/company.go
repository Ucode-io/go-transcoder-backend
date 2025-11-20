package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

func (r *postgresRepo) CompanyCreate(ctx context.Context, req *models.CompanyCreateReq) (*models.CompanyResponse, error) {
	res := &models.CompanyResponse{}
	query := r.Db.Builder.Insert("companies").Columns(
		"id",
		"title",
		"owner_id",
		"status",
	).Values(req.ID, req.Title, req.OwnerID, req.Status).Suffix(
		"RETURNING id, title, owner_id, status, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Title,
		&res.OwnerID, &res.Status,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) CompanyGet(ctx context.Context, req *models.CompanyGetReq) (*models.CompanyResponse, error) {
	query := r.Db.Builder.Select("id, title, owner_id, status, created_at, updated_at").
		From("companies")

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"id": req.ID})
	} else if req.OwnerId != "" {
		query = query.Where(squirrel.Eq{"owner_id": req.OwnerId})
	} else {
		return &models.CompanyResponse{}, fmt.Errorf(" hey at least one filter should be exists")
	}
	
	res := &models.CompanyResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Title,
		&res.OwnerID, &res.Status,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyCreate: query.RunWith(r.Db.Db).Scan()")
	}
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyGet:query.RunWith(r.Db.Db).QueryRow()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) CompanyFind(ctx context.Context, req *models.CompaniesFindReq) (*models.CompaniesFindResponse, error) {
	var (
		res            = &models.CompaniesFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"title": req.Search + "%"})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("companies").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompaniesFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, title, owner_id, status, created_at, updated_at").
		From("companies").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompaniesFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.CompanyResponse{}
		err := rows.Scan(
			&temp.ID, &temp.Title,
			&temp.OwnerID, &temp.Status,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "CompaniesFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)

		res.Companies = append(res.Companies, temp)
	}

	return res, nil
}

func (r *postgresRepo) CompanyUpdate(ctx context.Context, req *models.CompanyUpdateReq) (*models.CompanyResponse, error) {
	mp := make(map[string]interface{})
	if req.Status != "" {
		mp["status"] = req.Status
	}
	if req.Title != "" {
		mp["title"] = req.Title
	}
	mp["updated_at"] = time.Now()
	query := r.Db.Builder.Update("companies").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix("RETURNING id, title, owner_id, status, created_at, updated_at")

	res := &models.CompanyResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Title,
		&res.OwnerID, &res.Status,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) CompanyDelete(ctx context.Context, req *models.CompanyDeleteReq) error {
	query := r.Db.Builder.Delete("companies").Where(squirrel.Eq{"id": req.ID})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "CompanyDelete: query.RunWith(r.Db.Db).Exec()")
}
