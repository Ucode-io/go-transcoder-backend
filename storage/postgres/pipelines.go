package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"gitlab.com/transcodeuz/transcode-rest/models"
)

func (r *postgresRepo) PipelineCreate(ctx context.Context, req *models.PipelineCreateReq) (*models.PipelineResponse, error) {
	res := &models.PipelineResponse{}
	query := r.Db.Builder.Insert("pipelines").Columns(
		"id",
		"project_id",
		"fail_description",
		"input_url",
		"output_key",
		"output_path",
		"size_kb",
		"transcode_duration_seconds",
		"upload_duration_seconds",
		"max_resolution",
		"resolutions",
		"lang",
		"lang_code",
	).Values(req.ID, req.ProjectID, req.FailDescription, req.InputURL, req.OutputKey, req.OutputPath,
		req.SizeKB, req.TranscodeDurationSecs,
		req.UploadDurationSecs, req.MaxResolution, pq.Array(req.ResolutionsString), req.Language, req.LanguageCode).Suffix(
		`RETURNING  
		id, project_id, stage, stage_status , fail_description, 
		input_url, output_key, output_path, size_kb, transcode_duration_seconds, 
		upload_duration_seconds, max_resolution, 
		webhook_status, webhook_last_retry, webhook_retry_count, resolutions, lang, lang_code,
		created_at, updated_at`)

	nulTime := pq.NullTime{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.ProjectID,
		&res.Stage, &res.StageStatus, &res.FailDescription,
		&res.InputURL, &res.OutputKey,
		&res.OutputPath, &res.SizeKB,
		&res.TranscodeDurationSecs, &res.UploadDurationSecs, &res.MaxResolution,
		&res.WebhookStatus, &nulTime, &res.WebhookRetryCount, pq.Array(&res.ResolutionsString),
		&res.Language, &res.LanguageCode, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	query = r.Db.Builder.Insert("audio_tracks").Columns(
		"id",
		"pipeline_id",
		"size_kb",
		"input_url",
		"lang_code",
		"lang",
	)

	for _, e := range req.AudioTracks {
		query = query.Values(e.Id, res.ID, e.SizeKb, e.InputURL, e.LanguageCode, e.Language)
	}

	sqlQuery, args, _ := query.ToSql()
	_, err = r.Db.Db.Exec(sqlQuery, args...)
	if err != nil {
		return nil, HandleDatabaseError(err, r.Log, "creating audio tracks")
	}

	query = r.Db.Builder.Insert("subtitles").Columns(
		"id",
		"pipeline_id",
		"size_kb",
		"input_url",
		"lang_code",
		"lang",
	)

	for _, e := range req.Subtitle {
		query = query.Values(e.Id, res.ID, e.SizeKb, e.InputURL, e.LanguageCode, e.Language)
	}

	sqlQuerySubtitle, args, _ := query.ToSql()
	_, err = r.Db.Db.Exec(sqlQuerySubtitle, args...)
	if err != nil {
		return nil, HandleDatabaseError(err, r.Log, "creating audio tracks")
	}

	return res, nil
}

func (r *postgresRepo) PipelineGet(ctx context.Context, req *models.PipelineGetReq) (*models.PipelineResponse, error) {
	query := r.Db.Builder.Select(`id, project_id, stage, stage_status, 
	fail_description, input_url, output_key, output_path, size_kb, 
	transcode_duration_seconds, upload_duration_seconds, 
	max_resolution, resolutions, 
	video_duration, preparation_duration_seconds, 
	webhook_status, webhook_last_retry, webhook_retry_count, lang, lang_code, created_at, updated_at`).
		From("pipelines")

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"id": req.ID})
	} else {
		return &models.PipelineResponse{}, fmt.Errorf("at least one filter should be exists")
	}

	res := &models.PipelineResponse{}
	nulTime := pq.NullTime{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.ProjectID,
		&res.Stage, &res.StageStatus, &res.FailDescription,
		&res.InputURL, &res.OutputKey,
		&res.OutputPath, &res.SizeKB,
		&res.TranscodeDurationSecs, &res.UploadDurationSecs,
		&res.MaxResolution, pq.Array(&res.ResolutionsString),
		&res.VideoDuration, &res.PreparationDurationSeconds,
		&res.WebhookStatus, &nulTime, &res.WebhookRetryCount,
		&res.Language, &res.LanguageCode, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) PipelinesFind(ctx context.Context, req *models.PipelineFindReq) (*models.PipelinesFindResponse, error) {
	var (
		res            = &models.PipelinesFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"base_url": req.Search + "%"})
	}
	if strings.TrimSpace(req.ProjectId) != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"project_id": req.ProjectId})
	}
	if req.WebhookStatus != 0 {
		if req.WebhookStatus < 0 {
			whereCondition = append(whereCondition, squirrel.Eq{"webhook_status": false}, squirrel.Gt{"webhook_retry_count": 0}, squirrel.LtOrEq{"webhook_retry_count": 10})
		} else {
			whereCondition = append(whereCondition, squirrel.Eq{"webhook_status": true})
		}
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("pipelines").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select(`id, project_id, stage, stage_status, fail_description, 
	input_url, output_key, output_path, size_kb, transcode_duration_seconds, 
	upload_duration_seconds, max_resolution, resolutions, 
	video_duration, preparation_duration_seconds, 
	webhook_status, webhook_last_retry, webhook_retry_count, created_at, updated_at`).
		From("pipelines").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.PipelineResponse{}
		nulTime := pq.NullTime{}

		err := rows.Scan(
			&temp.ID, &temp.ProjectID,
			&temp.Stage, &temp.StageStatus, &temp.FailDescription,
			&temp.InputURL, &temp.OutputKey,
			&temp.OutputPath, &temp.SizeKB,
			&temp.TranscodeDurationSecs, &temp.UploadDurationSecs,
			&temp.MaxResolution, pq.Array(&temp.ResolutionsString),
			&temp.VideoDuration, &temp.PreparationDurationSeconds,
			&temp.WebhookStatus, &nulTime, &temp.WebhookRetryCount,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "PipelineFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Pipelines = append(res.Pipelines, temp)
	}

	return res, nil
}

func (r *postgresRepo) PipelineUpdate(ctx context.Context, req *models.PipelineUpdateReq) (*models.PipelineResponse, error) {
	mp := make(map[string]interface{})
	if req.Stage != "" {
		mp["stage"] = req.Stage
	}
	if req.StageStatus != "" {
		mp["stage_status"] = req.StageStatus
	}
	if req.PreparationDurationSeconds != 0 {
		mp["preparation_duration_seconds"] = req.PreparationDurationSeconds
	}
	if req.TranscodeDurationSecs != 0 {
		mp["transcode_duration_seconds"] = req.TranscodeDurationSecs
	}
	if req.UploadDurationSecs != 0 {
		mp["upload_duration_seconds"] = req.UploadDurationSecs
	}
	if req.FailDescription != "" {
		mp["fail_description"] = req.FailDescription
	}
	if req.VideoDuration != 0 {
		mp["video_duration"] = req.VideoDuration
	}
	if len(req.ResolutionsString) != 0 {
		mp["resolutions"] = pq.Array(req.ResolutionsString)
	}
	if req.InputURL != "" {
		mp["input_url"] = req.InputURL
	}
	if req.OutputKey != "" {
		mp["output_key"] = req.OutputKey
	}
	if req.OutputPath != "" {
		mp["output_path"] = req.OutputPath
	}
	if req.SizeKB != 0 {
		mp["size_kb"] = req.SizeKB
	}
	if req.MaxResolution != "" {
		mp["max_resolution"] = req.MaxResolution
	}
	if req.WebhookStatus {
		mp["webhook_status"] = req.WebhookStatus
	}
	if !req.WebhookLastRetry.IsZero() {
		mp["webhook_last_retry"] = req.WebhookLastRetry
	}
	if req.WebhookRetryCount != 0 {
		mp["webhook_retry_count"] = req.WebhookRetryCount
	}
	mp["updated_at"] = time.Now()
	// webhook_status, webhook_last_retry, webhook_retry_count,
	query := r.Db.Builder.Update("pipelines").SetMap(mp).
		Where(squirrel.Eq{"id": req.ID}).
		Suffix(`RETURNING 
		id, project_id, stage, stage_status, fail_description, 
		input_url, output_key, output_path, size_kb, transcode_duration_seconds, 
		upload_duration_seconds, max_resolution, resolutions, 
		video_duration, preparation_duration_seconds, 
		webhook_status, webhook_last_retry, webhook_retry_count, created_at, updated_at`)

	nulTime := pq.NullTime{}
	res := &models.PipelineResponse{}
	err := query.RunWith(r.Db.Db).Scan(
		&res.ID, &res.ProjectID,
		&res.Stage, &res.StageStatus, &res.FailDescription,
		&res.InputURL, &res.OutputKey,
		&res.OutputPath, &res.SizeKB,
		&res.TranscodeDurationSecs, &res.UploadDurationSecs,
		&res.MaxResolution, pq.Array(&res.ResolutionsString),
		&res.VideoDuration, &res.PreparationDurationSeconds,
		&res.WebhookStatus, &nulTime, &res.WebhookRetryCount,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)
	if nulTime.Valid {
		res.WebhookLastRetry = nulTime.Time
	}

	return res, nil
}

func (r *postgresRepo) PipelineDelete(ctx context.Context, req *models.PipelineDeleteReq) error {
	query := r.Db.Builder.Delete("pipelines").Where(squirrel.Eq{"id": req.ID})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "PipelineDelete: query.RunWith(r.Db.Db).Exec()")
}

func (r *postgresRepo) PipelineDashboarStatistics(ctx context.Context, req *models.DashboardStatisticsRequest) (*models.DashboardStatisticsResponse, error) {
	var (
		res            = &models.DashboardStatisticsResponse{}
		whereCondition = squirrel.And{}
	)

	if strings.TrimSpace(req.ProjectId) != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"project_id": req.ProjectId})
	}
	if strings.TrimSpace(req.CompanyId) != "" {
		whereCondition = append(whereCondition, squirrel.Eq{"company_id": req.ProjectId})
	}

	query := r.Db.Builder.Select(`
	COUNT(1) as total_count, 
    COUNT(CASE WHEN stage = 'upload' AND stage_status = 'success' THEN 1 END) as success_count,
    SUM(size_kb) as total_size `).From("pipelines").Where(whereCondition)

	err := query.RunWith(r.Db.Db).Scan(
		&res.TotalCount,
		&res.SuccessPercent,
		&res.TotalSize,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineDashboarStatistics: query.RunWith(r.Db.Db).Scan()")
	}
	res.SuccessPercent = uint32((100.0 / float32(res.TotalCount)) * float32(res.SuccessPercent))
	res.TotalSize /= 1024

	query2 := r.Db.Builder.Select(`
		DATE_TRUNC('day', created_at) AS date,
		COUNT(1) AS count,
		SUM(size_kb) AS size
	`).From("pipelines").Where(whereCondition).GroupBy("date").OrderBy("date DESC").Limit(uint64(req.LastNdays))

	rows, err := query2.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "PipelineDashboarStatistics: query2.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := models.PipelinePerDay{}

		err := rows.Scan(&temp.Date, &temp.Count, &temp.Size)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "PipelineDashboarStatistics: rows.Scan()")
		}
		temp.Size /= 1024
		temp.Date = strings.TrimSuffix(temp.Date, "T00:00:00Z")
		res.PipelinePerDay = append(res.PipelinePerDay, temp)
	}
	l := len(res.PipelinePerDay)
	for i := 0; i < l/2; i++ {
		res.PipelinePerDay[i], res.PipelinePerDay[l-1-i] = res.PipelinePerDay[l-1-i], res.PipelinePerDay[i]
	}
	return res, nil
}
