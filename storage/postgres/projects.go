package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

func (r *postgresRepo) ProjectCreate(ctx context.Context, req *models.ProjectCreateReq) (*models.ProjectResponse, error) {
	res := &models.ProjectResponse{}
	query := r.Db.Builder.Insert("projects").Columns(
		"id",
		"title",
		"access_key",
		"secret_key",
		"company_id",
		"owner_id",
		"status",
		"storage_id",
	).Values(req.ID, req.Title, req.AccessKey, req.SecretKey, req.CompanyID, req.OwnerID, req.Status, req.StorageID).Suffix(
		"RETURNING id, project_id, title, access_key, secret_key,company_id, owner_id, status, storage_id, created_at, updated_at")
	var projectID int
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &projectID, &res.Title,
		&res.AccessKey, &res.SecretKey,
		&res.CompanyID, &res.OwnerID,
		&res.Status, &res.StorageID,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "CompanyCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.ProjectID = "TR-" + strconv.Itoa(projectID)
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) ProjectGet(ctx context.Context, req *models.ProjectGetReq) (*models.ProjectResponse, error) {
	query := r.Db.Builder.Select("id, project_id, title, access_key, secret_key,company_id, owner_id,status, storage_id,created_at, updated_at").
		From("projects")

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"id": req.ID})
	} else if req.OwnerId != "" {
		query = query.Where(squirrel.Eq{"owner_id": req.OwnerId})
	} else if req.ProjectId != 0 {
		query = query.Where(squirrel.Eq{"project_id": req.ProjectId})
	} else {
		return &models.ProjectResponse{}, fmt.Errorf("at least one filter should be exists")
	}
	var projectID int
	res := &models.ProjectResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &projectID, &res.Title,
		&res.AccessKey, &res.SecretKey,
		&res.CompanyID, &res.OwnerID,
		&res.Status, &res.StorageID,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectCreate: query.RunWith(r.Db.Db).Scan()")
	}
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectGet:query.RunWith(r.Db.Db).QueryRow()")
	}
	res.ProjectID = "TR-" + strconv.Itoa(projectID)
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) ProjectGetID(ctx context.Context, ID int) (*models.ProjectResponse, error) {
	query := r.Db.Builder.Select("id, project_id, title, access_key, secret_key,company_id, owner_id, status, storage_id, created_at, updated_at").
		From("projects")

	if ID >= 1000000 {
		query = query.Where(squirrel.Eq{"project_id": ID})
	} else {
		return &models.ProjectResponse{}, fmt.Errorf("at least one filter should be exists")
	}
	var projectID int
	res := &models.ProjectResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &projectID, &res.Title,
		&res.AccessKey, &res.SecretKey,
		&res.CompanyID, &res.OwnerID,
		&res.Status, &res.StorageID,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectCreate: query.RunWith(r.Db.Db).Scan()")
	}
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectGet:query.RunWith(r.Db.Db).QueryRow()")
	}
	res.ProjectID = "TR-" + strconv.Itoa(projectID)
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) ProjectFind(ctx context.Context, req *models.ProjectsFindReq) (*models.ProjectsFindResponse, error) {
	var (
		res            = &models.ProjectsFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"base_url": req.Search + "%"})
	}
	if strings.TrimSpace(req.CompanyId) != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"company_id": req.CompanyId})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("projects").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, project_id, title, access_key, secret_key, company_id, owner_id, status, storage_id, created_at, updated_at").
		From("projects").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.ProjectResponse{}
		var projectID int
		err := rows.Scan(
			&temp.ID, &projectID, &temp.Title,
			&temp.AccessKey, &temp.SecretKey,
			&temp.CompanyID, &temp.OwnerID,
			&temp.Status, &temp.StorageID,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "CompanyFind: rows.Scan()")
		}
		temp.ProjectID = "TR-" + strconv.Itoa(projectID)
		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Projects = append(res.Projects, temp)
	}

	return res, nil
}

func (r *postgresRepo) ProjectUpdate(ctx context.Context, req *models.ProjectUpdateReq) (*models.ProjectResponse, error) {
	mp := make(map[string]interface{})
	mp["title"] = req.Title
	mp["access_key"] = req.AccessKey
	mp["secret_key"] = req.SecretKey
	mp["status"] = req.Status
	mp["updated_at"] = time.Now()
	query := r.Db.Builder.Update("projects").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix("RETURNING id, title, access_key, secret_key,company_id, owner_id,status, storage_id,created_at, updated_at")

	res := &models.ProjectResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Title,
		&res.AccessKey, &res.SecretKey,
		&res.CompanyID, &res.OwnerID,
		&res.Status, &res.StorageID,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) ProjectUpdateName(ctx context.Context, req *models.ProjectNameUpdateReq) (*models.ProjectResponse, error) {
	mp := make(map[string]interface{})
	mp["title"] = req.Title
	query := r.Db.Builder.Update("projects").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix("RETURNING id, project_id, title, access_key, secret_key,company_id, owner_id, status, storage_id,created_at, updated_at")

	var projectID int
	res := &models.ProjectResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &projectID, &res.Title,
		&res.AccessKey, &res.SecretKey,
		&res.CompanyID, &res.OwnerID,
		&res.Status, &res.StorageID,
		&CreatedAt, &UpdatedAt,
	)

	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "ProjectUpdateName: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}
	res.ProjectID = "TR-" + strconv.Itoa(projectID)
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)
	return res, nil
}

func (r *postgresRepo) ProjectDelete(ctx context.Context, req *models.ProjectDeleteReq) error {
	query := r.Db.Builder.Delete("projects").Where(squirrel.Eq{"id": req.ID})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "ProjectDelete: query.RunWith(r.Db.Db).Exec()")
}
