package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const TemplateFormatTime = "15:04"
const InitTime = "00:00"

func Run(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	club := &Club{
		clients:           make(map[string]*Client),
		tableOccupancy:    make(map[int]*Client),
		eventLog:          make([]string, 0),
		revenue:           make(map[int]int),
		totalTimeOccupied: make(map[int]time.Duration),
	}

	var events []string

	scanner.Scan()
	club.tables, err = strconv.Atoi(scanner.Text())
	if err != nil {
		formatError(scanner.Text())
	}

	scanner.Scan()
	times := strings.Split(scanner.Text(), " ")
	club.openTime, err = time.Parse(TemplateFormatTime, times[0])
	if err != nil {
		formatError(scanner.Text())
	}
	club.closeTime, err = time.Parse(TemplateFormatTime, times[1])
	if err != nil {
		formatError(scanner.Text())
	}

	scanner.Scan()
	club.pricePerHour, err = strconv.Atoi(scanner.Text())
	if err != nil {
		formatError(scanner.Text())
	}

	for scanner.Scan() {
		events = append(events, scanner.Text())
	}

	prevTimeEvent, _ := time.Parse(TemplateFormatTime, InitTime)

	for _, event := range events {
		if processEvent(event, club, &prevTimeEvent) {
			formatError(event)
		}
	}

	club.RemoveUsersAfterTime()
	club.printResults()
}

func formatError(line string) {
	fmt.Println("Input format error: " + line)
	os.Exit(1)
}

func logEvent(eventTime time.Time, eventID int, args ...string) string {
	return fmt.Sprintf("%s %d %s", eventTime.Format(TemplateFormatTime), eventID, strings.Join(args, " "))
}

func processEvent(event string, club *Club, prevTimeEvent *time.Time) bool {
	parts := strings.Split(event, " ")
	eventTime, err := time.Parse(TemplateFormatTime, parts[0])
	if err != nil {
		return true
	}
	if prevTimeEvent.After(eventTime) {
		return true
	}
	*prevTimeEvent = eventTime

	eventID, err := strconv.Atoi(parts[1])
	if err != nil || eventID > 4 || eventID < 1 {
		return true
	}
	clientName := parts[2]

	if eventTime.After(club.closeTime) {
		club.RemoveUsersAfterTime()
	}

	switch eventID {
	case 1:
		club.eventLog = append(club.eventLog, logEvent(eventTime, 1, clientName))
		if _, exists := club.clients[clientName]; exists {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 13, "YouShallNotPass"))
		} else if eventTime.Before(club.openTime) || eventTime.After(club.closeTime) {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 13, "NotOpenYet"))
		} else {
			club.clients[clientName] = &Client{name: clientName}
		}
	case 2:
		tableNumber, err := strconv.Atoi(parts[3])
		if err != nil || tableNumber > club.tables {
			return true
		}
		club.eventLog = append(club.eventLog, logEvent(eventTime, 2, clientName, strconv.Itoa(tableNumber)))
		client, exists := club.clients[clientName]
		if !exists {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 13, "ClientUnknown"))
		} else if currentClient, occupied := club.tableOccupancy[tableNumber]; occupied && currentClient.name != clientName {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 13, "PlaceIsBusy"))
		} else {
			if client.table != 0 {
				club.calculateRevenue(client, eventTime)
				club.tableOccupancy[client.table] = nil
			}
			client.table = tableNumber
			client.startTime = eventTime
			club.tableOccupancy[tableNumber] = client
		}
	case 3:
		club.eventLog = append(club.eventLog, logEvent(eventTime, 3, clientName))
		if len(club.tableOccupancy) < club.tables {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 13, "ICanWaitNoLonger!"))
		} else if len(club.queue) < club.tables {
			club.queue = append(club.queue, &Client{name: clientName})
		} else {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 11, clientName))
		}
	case 4:
		club.eventLog = append(club.eventLog, logEvent(eventTime, 4, clientName))
		client, exists := club.clients[clientName]
		if !exists {
			club.eventLog = append(club.eventLog, logEvent(eventTime, 13, "ClientUnknown"))
		} else {
			club.calculateRevenue(client, eventTime)
			club.tableOccupancy[client.table] = nil
			delete(club.clients, clientName)

			if len(club.queue) != 0 {
				table, clientNameQueue := club.TakeAnEmptySeat(eventTime)
				club.eventLog = append(club.eventLog, logEvent(eventTime, 12, clientNameQueue, strconv.Itoa(table)))
			}
		}
	}
	return false
}
