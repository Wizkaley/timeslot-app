package repository

import (
	"fmt"
	"timeslot-app/models"

	"github.com/jackc/pgx"
)

type TimeslotRepoImplementation struct {
	db *pgx.Conn
}

func NewTimeslotRepository(dbconn *pgx.Conn) TimeslotRepo {
	return &TimeslotRepoImplementation{
		db: dbconn,
	}
}

type TimeslotRepo interface {
	Create(timeSlots []models.TimeSlot) error
	DeleteTimeSlotsByUserName(userName, timeSlot string) error
	GetTimeSlotsByUserName(userName string) ([]string, error)
}

func (ts *TimeslotRepoImplementation) Create(timeSlots []models.TimeSlot) error {

	for _, slot := range timeSlots {
		insertQuery := `INSERT INTO time_slots (id, user_id, time_slot) VALUES ($1, $2, $3)`
		_, err := ts.db.Exec(insertQuery, slot.ID, slot.UserID, slot.TimeSlot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TimeslotRepoImplementation) GetTimeSlotsByUserName(userName string) ([]string, error) {
	qry := `select ts.time_slot from users u 
	join time_slots ts on u.id=ts.user_id
	where u.name = $1`

	rows, err := ts.db.Query(qry, userName)
	if err != nil {
		return []string{}, err
	}
	var timeSlots []string
	for rows.Next() {

		var timeSlot string
		err := rows.Scan(&timeSlot)
		if err != nil {
			return nil, err
		}

		timeSlots = append(timeSlots, timeSlot)
	}

	fmt.Println("timeSlots for UserName:", userName, "timeslots: ", timeSlots)
	return timeSlots, nil
}

func (ts *TimeslotRepoImplementation) DeleteTimeSlotsByUserName(userName, timeSlot string) error {

	deleteQuery := `delete from time_slots where user_id in (select id from users where name=$1) and time_slot=$2`
	_, err := ts.db.Exec(deleteQuery, userName, timeSlot)
	if err != nil {
		return err
	}
	return nil
}
