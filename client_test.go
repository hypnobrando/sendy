package sendy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	TestTodoObject struct {
		UserID int    `json:"userId`
		Title  string `json:"title`
	}

	TestPostObject struct {
		ID int `json:"id"`
	}
)

func TestGet(t *testing.T) {
	var response TestTodoObject

	err := Get("https://jsonplaceholder.typicode.com").
		Path("/todos/1").
		SendIt().
		JSON(&response).
		Error()

	assert.Nil(t, err)
	assert.Equal(t, 1, response.UserID)
	assert.Equal(t, "delectus aut autem", response.Title)
}

func TestPost(t *testing.T) {
	var response TestPostObject

	err := Post("https://jsonplaceholder.typicode.com").
		Path("/posts").
		JSON(map[string]int{
			"id": 101,
		}).
		SendIt().
		JSON(&response).
		Error()

	assert.Nil(t, err)
	assert.Equal(t, 101, response.ID)
}

type (
	testBody struct {
		ExternalID      string                 `json:"external_id"`
		TicketProceeds  testBodyTicketProceeds `json:"ticket_proceeds"`
		Seating         testBodySeating        `json:"seating"`
		TicketType      string                 `json:"ticket_type"`
		SplitType       string                 `json:"split_type"`
		NumberOfTickets int                    `json:"number_of_tickets"`
		Event           testBodyEvent          `json:"event"`
		Venue           testBodyVenue          `json:"venue"`
	}

	testBodyTicketProceeds struct {
		Amount       int    `json:"amount"`
		CurrencyCode string `json:"currency_code"`
		Display      string `json:"display"`
	}

	testBodySeating struct {
		Section string `json:"section"`
		Row     string `json:"row"`
	}

	testBodyEvent struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
	}

	testBodyVenue struct {
		Name string `json:"name"`
	}
)

func TestContentLength(t *testing.T) {
	request := NewClient().
		Host("https://google.com").
		Post().
		JSON(testBody{
			ExternalID: "1239-22112",
			TicketProceeds: testBodyTicketProceeds{
				Amount:       30,
				CurrencyCode: "USD",
				Display:      "$30.00",
			},
			Seating: testBodySeating{
				Section: "General Admission",
				Row:     "General Admission",
			},
			TicketType:      "ETicketUrl",
			SplitType:       "Any",
			NumberOfTickets: 99,
			Event: testBodyEvent{
				Name:      "Keystone Revisited",
				StartDate: "2023-09-08T21:30:00",
			},
			Venue: testBodyVenue{
				Name: "The Lariat in Buena Vista",
			},
		})

	for _, h := range request.headers {
		if h.Key == "Content-Length" {
			assert.Equal(t, "356", h.Value)
			return
		}
	}

	t.Fatal("Content-Length header not found")
}
