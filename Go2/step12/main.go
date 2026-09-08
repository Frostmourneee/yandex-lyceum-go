package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"
)

const (
	ticketPrefix = "TICKET-"
	dateLayout   = "2006-01-02"
)

var errTimeout = errors.New("timeout")

var validStatuses = map[string]struct{}{
	"Готово":           {},
	"В работе":         {},
	"Не будет сделано": {},
}

func main() {

}

type Ticket struct {
	Ticket string
	User   string
	Status string
	Date   time.Time
}

func ParseTicket(text string) (Ticket, bool) {
	fields := strings.Split(text, "_")
	if len(fields) != 4 {
		return Ticket{}, false
	}
	id, user, status, rawDate := fields[0], fields[1], fields[2], fields[3]

	if !strings.HasPrefix(id, ticketPrefix) {
		return Ticket{}, false
	}

	if _, ok := validStatuses[status]; !ok {
		return Ticket{}, false
	}

	date, err := time.Parse(dateLayout, rawDate)
	if err != nil {
		return Ticket{}, false
	}

	return Ticket{
		Ticket: id,
		User:   user,
		Status: status,
		Date:   date,
	}, true
}

func filterUser(ticket Ticket, user *string) bool {
	if user == nil {
		return true
	}
	return ticket.User == *user
}

func filterStatus(ticket Ticket, status *string) bool {
	if status == nil {
		return true
	}
	return ticket.Status == *status
}

func readTickets(ctx context.Context, r io.Reader, user, status *string) ([]Ticket, error) {
	var tickets []Ticket

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		row := strings.TrimSpace(scanner.Text())
		if row == "" {
			continue
		}

		ticket, ok := ParseTicket(row)
		if !ok {
			continue
		}

		if !filterUser(ticket, user) || !filterStatus(ticket, status) {
			continue
		}

		tickets = append(tickets, ticket)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tickets, nil
}

func GetTasks(
	ctx context.Context,
	r io.Reader,
	w io.Writer,
	user *string,
	status *string,
	timeout time.Duration,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	type result struct {
		tickets []Ticket
		err     error
	}

	// Buffered so the goroutine can always finish, even if nobody reads the result.
	done := make(chan result, 1)
	go func() {
		tickets, err := readTickets(ctx, r, user, status)
		done <- result{tickets: tickets, err: err}
	}()

	select {
	case <-ctx.Done():
		return asTimeout(ctx.Err())
	case res := <-done:
		if res.err != nil {
			return asTimeout(res.err)
		}

		data, err := json.Marshal(res.tickets)
		if err != nil {
			return err
		}

		_, err = w.Write(data)
		return err
	}
}

func asTimeout(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errTimeout
	}
	return err
}
