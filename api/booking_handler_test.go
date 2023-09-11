package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexeyAndryushin/reservations/db/fixtures"
	"github.com/AlexeyAndryushin/reservations/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

		var (
		nonAuthUser    = fixtures.AddUser(db.Store, "admin", "admin", false)
		user           = fixtures.AddUser(db.Store, "alex", "noadmin", false)
		hotel          = fixtures.AddHotel(db.Store, "good hotel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
		)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "alex", "alex", false)
		user           = fixtures.AddUser(db.Store, "alex", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "hotel hotel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}