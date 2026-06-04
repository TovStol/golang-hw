package sqlstorage

import (
	"fmt"
	"strings"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	db           *sqlx.DB
	dbDriverName string
	dsn          string
}

func New(dbDriverName string, dsn string) *SQLStorage {
	dbDriverName = strings.ToLower(dbDriverName)
	dsn = strings.TrimSpace(dsn)
	return &SQLStorage{dbDriverName: dbDriverName, dsn: dsn}
}

func (r *SQLStorage) Connect() error {
	r.db = sqlx.MustConnect(r.dbDriverName, r.dsn)
	// Настройки ниже конфигурируют пулл подключений к базе данных. Их названия стандартны для большинства библиотек.
	// Ознакомиться с их описанием можно на примере документации Hikari pool:
	// https://github.com/brettwooldridge/HikariCP?tab=readme-ov-file#gear-configuration-knobs-baby
	r.db.SetMaxIdleConns(5)
	r.db.SetMaxOpenConns(20)
	r.db.SetConnMaxLifetime(1 * time.Minute)
	r.db.SetConnMaxIdleTime(10 * time.Minute)
	return nil
}

func (r *SQLStorage) Close() error {
	err := r.db.Close()
	if err != nil {
		fmt.Print(err)
	}
	return nil
}

func (r *SQLStorage) Create(event models.Event) (id int64, err error) {
	result, err := r.db.Query(
		`INSERT INTO event (title, date_time, end_date_time, description, user_id, notify_before)
		 VALUES($1, $2, $3, $4, $5, $6) RETURNING id`,
		event.Title, event.DateTime, event.EndDateTime, event.Description, event.UserID, int64(event.NotifyBefore))
	if err != nil {
		return 0, err
	}
	if result.Next() {
		err = result.Scan(&id)
		return id, err
	}
	return id, err
}

func (r *SQLStorage) Update(event models.Event) {
	r.db.Query( //nolint:errcheck
		`UPDATE event SET title=$1, date_time=$2, end_date_time=$3, description=$4, user_id=$5, notify_before=$6
		 WHERE id=$7`,
		event.Title, event.DateTime, event.EndDateTime, event.Description, event.UserID, int64(event.NotifyBefore), event.ID)
}

func (r *SQLStorage) DeleteByID(eventID int64) (err error) {
	_, err = r.db.Query(
		"DELETE from event where id = $1", eventID)
	return err
}

func (r *SQLStorage) FindEventsByDay(date time.Time) (res []models.Event, err error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	err = r.db.Select(&res,
		"SELECT * FROM event WHERE date_time >= $1 AND date_time < $2", start, end)
	return res, err
}

func (r *SQLStorage) FindEventsByWeek(date time.Time) (res []models.Event, err error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 7)
	err = r.db.Select(&res,
		"SELECT * FROM event WHERE date_time >= $1 AND date_time < $2", start, end)
	return res, err
}

func (r *SQLStorage) FindEventsByMonth(date time.Time) (res []models.Event, err error) {
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0)
	err = r.db.Select(&res,
		"SELECT * FROM event WHERE date_time >= $1 AND date_time < $2", start, end)
	return res, err
}

func (r *SQLStorage) FindByID(id int64) (res models.Event, err error) {
	err = r.db.Get(&res,
		"SELECT * FROM event WHERE id = $1", id)
	return res, err
}

func (r *SQLStorage) FindAll() (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event")
	return res, err
}

func (r *SQLStorage) FindByIDs(ids []int64) (res []models.Event, err error) {
	query, args, err := sqlx.In(
		"SELECT * from event where id IN(?)", ids)
	if err == nil {
		query = r.db.Rebind(query)
		var rows *sqlx.Rows
		rows, err = r.db.Queryx(query, args...)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var event models.Event
			err = rows.StructScan(&event)
			if err != nil {
				return nil, err
			}
			res = append(res, event)
		}
		return res, err
	}
	return make([]models.Event, 0), err
}

func (r *SQLStorage) DeleteByIDs(ids []int64) (err error) {
	query, args, err := sqlx.In(
		"DELETE from event where id IN(?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Query(query, args...)
	return err
}

func (r *SQLStorage) ExecuteQuery(query string) {
	_, err := r.db.Exec(query)
	if err != nil {
		return
	}
}

func (r *SQLStorage) FindEventsToNotify(now time.Time) ([]models.Event, error) {
	var res []models.Event
	err := r.db.Select(&res,
		`SELECT * FROM event
		 WHERE notify_before > 0
		   AND date_time - (notify_before / 1000 * interval '1 microsecond') <= $1
		   AND date_time > $1`,
		now)
	return res, err
}

func (r *SQLStorage) DeleteOldEvents(before time.Time) error {
	_, err := r.db.Exec(
		"DELETE FROM event WHERE date_time < $1", before)
	return err
}

func (r *SQLStorage) SaveNotification(n models.Notification) error {
	_, err := r.db.Exec(
		`INSERT INTO notification (event_id, title, event_date, user_id)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (event_id) DO NOTHING`,
		n.EventID, n.Title, n.EventDate, n.UserID)
	return err
}
