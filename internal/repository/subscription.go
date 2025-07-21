package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Headliner38/Subscription_Service/internal/model"
)

func CreateSubscription(db *sql.DB, id, serviceName string, price int, userID string, startDate time.Time, endDate *time.Time) error {
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
	 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, id, serviceName, price, userID, startDate, endDate)
	if err != nil {
		return err
	}
	return err
}

func GetSubscription(db *sql.DB, id string) (*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	row := db.QueryRow(query, id)

	var sub model.Subscription
	var endDate sql.NullTime

	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, err
	}

	if endDate.Valid {
		sub.EndDate = &endDate.Time
	} else {
		sub.EndDate = nil
	}

	return &sub, nil
} // реализовать если успею GetSubscriptionByID, ByServiceName и GetByPrice

func UpdateSubscription(db *sql.DB, id, serviceName string, price int, userID string, startDate time.Time, endDate *time.Time) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 
	WHERE id = $6`

	result, err := db.Exec(query, serviceName, price, userID, startDate, endDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func DeleteSubscription(db *sql.DB, id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func ListSubscriptions(db *sql.DB) ([]model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {
		var sub model.Subscription
		var endDate sql.NullTime

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			sub.EndDate = &endDate.Time
		} else {
			sub.EndDate = nil
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func CalculateTotalCost(db *sql.DB, userID, serviceName string, startDate, endDate time.Time) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if userID != "" {
		query += ` AND user_id = $` + fmt.Sprint(argIdx)
		args = append(args, userID)
		argIdx++
	}
	if serviceName != "" {
		query += ` AND service_name = $` + fmt.Sprint(argIdx)
		args = append(args, serviceName)
		argIdx++
	}
	if !startDate.IsZero() {
		query += ` AND start_date >= $` + fmt.Sprint(argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if !endDate.IsZero() {
		query += ` AND start_date <= $` + fmt.Sprint(argIdx)
		args = append(args, endDate)
		argIdx++
	}

	var totalCost int
	err := db.QueryRow(query, args...).Scan(&totalCost)
	if err != nil {
		return 0, err
	}

	return totalCost, nil
}
