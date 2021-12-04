package dbrepo

import (
	"errors"
	"time"

	"github.com/adrialopezbou/bookings-go/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pas
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestrictions inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDates returns true if availability exists
func (m *testDBRepo) SearchAvailabilityByDatesAndRoomId(start, end time.Time, roomID int) (bool, error) {
	if roomID == 2 {
		return false, errors.New("some error")
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time)  ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}
	
	return room, nil
}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}


func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}


func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	return 0, "", nil
}