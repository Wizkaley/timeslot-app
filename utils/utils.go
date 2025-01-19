package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"timeslot-app/models"
)

func CreateFile(b []byte) {
	fmt.Println("Writing file")
	file, err := os.Create("swagger/swagger.json")

	if err != nil {
		panic(err)
	}

	length, err := file.WriteString(string(b))

	if err != nil {
		panic(err)
	}
	fmt.Printf("File name: %s", file.Name())
	fmt.Printf("\nfile length: %d\n", length)
}

// for this example we will assume the timezone is EST
func ValidateTimeStamp(ts string) bool {
	// layout := "2 Jan 2025 2-4 PM EST"

	sp := strings.Split(ts, " ")
	if len(sp) < 6 {
		fmt.Println("Invalid timestamp1")
		return false

	}

	date := strings.Join(sp[0:3], " ")
	partOfDay := strings.ToUpper(sp[4])
	timezone := sp[5] // "EST"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return false
	}

	// split the time slots

	timeSlots := strings.Split(sp[3], "-")
	if len(timeSlots) < 2 {
		fmt.Println("Invalid timestamp2")
		return false
	}
	slotStart := fmt.Sprintf("%s %s %s", date, timeSlots[0], partOfDay)
	slotEnd := fmt.Sprintf("%s %s %s", date, timeSlots[1], partOfDay)

	layout := "02 Jan 2006 3 PM"
	ss, err := time.ParseInLocation(layout, slotStart, loc)
	if err != nil {
		fmt.Println("Invalid timestamp3")
		return false
	}

	se, err := time.ParseInLocation(layout, slotEnd, loc)
	if err != nil {
		fmt.Println("Invalid timestamp4")
		return false
	}

	if ss.After(se) {
		fmt.Println("Invalid timestamp5")
		return false
	}

	fmt.Println(ss, se)
	return true
}

func ValidateAndFormatTimeStamp(ts string) (startTime, endTime time.Time, valid bool) {

	sp := strings.Split(ts, " ")
	if len(sp) < 6 {
		fmt.Println("Invalid timestamp1")
		return time.Time{}, time.Time{}, false

	}

	date := strings.Join(sp[0:3], " ")
	partOfDay := strings.ToUpper(sp[4])
	timezone := sp[5] // "EST"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}

	// split the time slots

	timeSlots := strings.Split(sp[3], "-")
	if len(timeSlots) < 2 {
		fmt.Println("Invalid timestamp2")
		return time.Time{}, time.Time{}, false
	}
	slotStart := fmt.Sprintf("%s %s %s", date, timeSlots[0], partOfDay)
	slotEnd := fmt.Sprintf("%s %s %s", date, timeSlots[1], partOfDay)

	layout := "02 Jan 2006 3 PM"
	ss, err := time.ParseInLocation(layout, slotStart, loc)
	if err != nil {
		fmt.Println("Invalid timestamp3")
		return time.Time{}, time.Time{}, false
	}

	se, err := time.ParseInLocation(layout, slotEnd, loc)
	if err != nil {
		fmt.Println("Invalid timestamp4")
		return time.Time{}, time.Time{}, false
	}

	if ss.After(se) {
		fmt.Println("Invalid timestamp5")
		return time.Time{}, time.Time{}, false
	}

	fmt.Println(ss, se)
	return ss, se, true
}

func CheckIfTimeSlotsOverlap(timeSlot1, timeSlot2 models.TimeSlotStartAndEnd, eventDuration time.Duration) bool {
	slotStart := func(a, b time.Time) time.Time {
		if a.After(b) {
			return a
		}
		return b
	}(timeSlot1.StartTime, timeSlot2.StartTime)

	slotEnd := func(a, b time.Time) time.Time {
		if a.Before(b) {
			return a
		}
		return b
	}(timeSlot1.EndTime, timeSlot2.EndTime)
	duration := slotEnd.Sub(slotStart)
	return duration >= eventDuration
}

func DisplayTimeSlots(timeSlots []models.TimeSlotStartAndEnd) {
	for _, ts := range timeSlots {
		fmt.Printf("Perfect Slots Start time: %s, End time: %s\n", ts.StartTime, ts.EndTime)
	}
}

func DisplayPartialMatchSlots(timeSlots []models.MatchingEventSlots) {
	for _, ts := range timeSlots {
		fmt.Printf("Partial Slots Start time: %s, End time: %s\n", ts.Slot.StartTime, ts.Slot.EndTime)
		fmt.Printf("Available Participants: %v\n", ts.AvailableParticipants)
		fmt.Printf("Unavailable Participants: %v\n", ts.UnavailableParticipants)
	}
}

func ErrorHelper(message string, err error) []byte {
	servErr := models.ServiceError{Message: "Bad request", Error: err.Error()}
	marshalledErr, _ := json.Marshal(servErr)
	return marshalledErr
}

func SuccessHelper(message string) []byte {
	servErr := models.ServiceMessage{Message: message}
	marshalledErr, _ := json.Marshal(servErr)
	return marshalledErr
}

func SearchString(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
