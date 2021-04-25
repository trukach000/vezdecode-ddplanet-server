package support

import (
	"context"
	"ddplanet-server/pkg/database"
	"fmt"
	"time"
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
