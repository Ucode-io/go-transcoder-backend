package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

func (r *postgresRepo) UserCreate(ctx context.Context, req *models.UserCreateReq) (*models.UserResponse, error) {
	res := &models.UserResponse{}
	query := r.Db.Builder.Insert("users").Columns(
		"id",
		"username",
		"password",
		"refresh_token",
		"user_role",
	).Values(req.ID, req.Username, req.Password, "", req.Role).Suffix(
		"RETURNING id, username, refresh_token, user_role, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Username,
		&res.RefreshToken, &res.Role,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "UserCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) UserGet(ctx context.Context, req *models.UserGetReq) (*models.UserResponse, error) {
	query := r.Db.Builder.Select("id, username, password, user_role, created_at, updated_at").
		From("users")

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"id": req.ID})
	} else if req.Username != "" {
		query = query.Where(squirrel.Eq{"username": req.Username})
	} else {
		return &models.UserResponse{}, fmt.Errorf("at least one filter should be exists")
	}

	res := &models.UserResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Username,
		&res.Password,
		&res.Role,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "UserCreate: query.RunWith(r.Db.Db).Scan()")
	}
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "UserGet:query.RunWith(r.Db.Db).QueryRow()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) UserFind(ctx context.Context, req *models.UserFindReq) (*models.UserFindResponse, error) {
	var (
		res            = &models.UserFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"username": req.Search + "%"})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("users").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "UserFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, username, user_role, created_at, updated_at").
		From("users").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "UserFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.UserResponse{}
		err := rows.Scan(
			&temp.ID, &temp.Username,
			&temp.Role, &CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "UserFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)

		res.Users = append(res.Users, temp)
	}

	return res, nil
}

func (r *postgresRepo) UserUpdate(ctx context.Context, req *models.UserUpdateReq) (*models.UserResponse, error) {
	res := &models.UserResponse{}
	mp := make(map[string]interface{})

	if req.Username != "" {
		mp["username"] = req.Username
	}
	if req.Password != "" {
		mp["password"] = req.Password
	}
	if req.RefreshToken != "" {
		mp["refresh_token"] = req.RefreshToken
	}
	if req.Role != "" {
		mp["user_role"] = req.Role
	}
	mp["updated_at"] = time.Now()

	query := r.Db.Builder.Update("users").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix("RETURNING id, username, user_role, refresh_token, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.Username,
		&res.AccessToken,
		&res.RefreshToken,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "UserUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) UserDelete(ctx context.Context, req *models.UserDeleteReq) error {
	query := r.Db.Builder.Delete("users").Where(squirrel.Eq{"id": req.ID})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "UserDelete: query.RunWith(r.Db.Db).Exec()")
}
