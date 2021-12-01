package repository

import (
	"time"

	"github.com/adrialopezbou/bookings-go/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time)  ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
}
