package support

import (
	"context"
	"ddplanet-server/app/supporting"
	"ddplanet-server/pkg/database"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

func CreateSupportTask(ctx context.Context, firstName, secondName, lastName, phone, message string) (int64, error) {
	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return -1, err
	}

	result, err := db.ExecContext(
		ctx,
		`INSERT INTO ddplanet.support_requests
				(
				ts_created,
				first_name,
				second_name,
				last_name,
				phone,
				message)
				VALUES
				(?,?,?,?,?,?);
				`,
		time.Now().Unix(),
		firstName,
		secondName,
		lastName,
		phone,
		message,
	)
	if err != nil {
		return -1, err
	}

	requestId, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("can;t get last insert ID for new request task")
	}

	return requestId, nil
}

func GetRequests(
	ctx context.Context,
	statuses []supporting.SupportRequestStatus,
	tsCreatedFrom, tsCreatedTo, limit int64,
	searchPhoneOrID string,
) ([]supporting.SupportRequest, error) {
	requests := make([]supporting.SupportRequest, 0)

	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return requests, err
	}

	if limit == 0 {
		limit = 100
	}

	sq := squirrel.
		Select(
			"id",
			"ts_created",
			"first_name",
			"second_name",
			"last_name",
			"phone",
			"message",
			"status",
			"ts_closed",
		).
		From("support_requests")

	if len(statuses) > 0 {
		sq = sq.Where(squirrel.Eq{"status": statuses})
	}

	if tsCreatedFrom > 0 {
		sq = sq.Where(squirrel.GtOrEq{"ts_created": tsCreatedFrom})
	}

	if tsCreatedTo > 0 {
		sq = sq.Where(squirrel.LtOrEq{"ts_created": tsCreatedTo})
	}

	if searchPhoneOrID != "" {
		sq = sq.Where(
			squirrel.Or{
				squirrel.Like{"phone": fmt.Sprint("%", searchPhoneOrID, "%")},
				squirrel.Eq{"id": searchPhoneOrID},
			},
		)

	}

	sq = sq.Limit(uint64(limit))

	sql, args, err := sq.ToSql()
	logrus.Infof("sql: %s", sql)

	if err != nil {
		return requests, fmt.Errorf("can't prepare sql query in `GetRequests`: %s", err)
	}

	rows, err := db.QueryContext(
		ctx,
		sql,
		args...,
	)
	if err != nil {
		return requests, fmt.Errorf("sql query error in `GetRequests`: %s", err)
	}

	for rows.Next() {
		var req supporting.SupportRequest
		scanErr := rows.Scan(
			&req.ID,
			&req.CreatedTs,
			&req.FirstName,
			&req.SecondName,
			&req.LastName,
			&req.Phone,
			&req.Message,
			&req.Status,
			&req.ClosedTs,
		)
		if scanErr != nil {
			logrus.Warningf("can't scan row of support requests: %s", scanErr)
			return requests, scanErr
		}

		requests = append(requests, req)
	}

	return requests, nil
}
