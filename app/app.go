package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Run() {
	filename := "test_file.txt"
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
	club.tables, _ = strconv.Atoi(scanner.Text())

	scanner.Scan()
	times := strings.Split(scanner.Text(), " ")
	club.openTime, _ = time.Parse("15:04", times[0])
	club.closeTime, _ = time.Parse("15:04", times[1])

	scanner.Scan()
	club.pricePerHour, _ = strconv.Atoi(scanner.Text())

	for scanner.Scan() {
		events = append(events, scanner.Text())
	}

	for _, event := range events {
		parts := strings.Split(event, " ")
		eventTime, err := time.Parse("15:04", parts[0])
		if err != nil {
			fmt.Println(event)
			os.Exit(1)
		}
		eventId, err := strconv.Atoi(parts[1])
		if err != nil || eventId > 4 || eventId < 1 {
			fmt.Println(event)
			os.Exit(1)
		}
		clientName := parts[2]

		if eventTime.After(club.closeTime) {
			club.RemoveUsersAfterTime()
		}

		switch eventId {
		case 1:
			if _, exists := club.clients[clientName]; exists {
				outLog := eventTime.Format("15:04") + " 13 " + "YouShallNotPass"
				club.eventLog = append(club.eventLog, outLog)
			} else if eventTime.Before(club.openTime) || eventTime.After(club.closeTime) {
				outLog := eventTime.Format("15:04") + " 1 " + clientName
				club.eventLog = append(club.eventLog, outLog)
				outLog = eventTime.Format("15:04") + " 13 " + "NotOpenYet"
				club.eventLog = append(club.eventLog, outLog)
			} else {
				club.clients[clientName] = &Client{name: clientName}
				outLog := eventTime.Format("15:04") + " 1 " + clientName
				club.eventLog = append(club.eventLog, outLog)
			}
		case 2:
			tableNumber, _ := strconv.Atoi(parts[3])
			if tableNumber > club.tables {
				fmt.Println(event)
				os.Exit(1)
			}
			if client, exists := club.clients[clientName]; !exists {
				outLog := eventTime.Format("15:04") + " 13 " + "ClientUnknown"
				club.eventLog = append(club.eventLog, outLog)
			} else if currentClient, occupied := club.tableOccupancy[tableNumber]; occupied && currentClient.name != clientName {
				outLog := eventTime.Format("15:04") + " 13 " + "PlaceIsBusy"
				club.eventLog = append(club.eventLog, outLog)
			} else {
				if client.table != 0 {
					club.tableOccupancy[client.table] = nil
					club.tableOccupancy[tableNumber] = client
					client.table = tableNumber
				}
				client.table = tableNumber
				client.startTime = eventTime
				club.tableOccupancy[tableNumber] = client
				outLog := eventTime.Format("15:04") + " 2 " + clientName + " " + strconv.Itoa(tableNumber)
				club.eventLog = append(club.eventLog, outLog)
			}
		case 3:
			if len(club.tableOccupancy) < club.tables {
				outLog := eventTime.Format("15:04") + " 13 " + "ICanWaitNoLonger!"
				club.eventLog = append(club.eventLog, outLog)
				continue
			}
			outLog := eventTime.Format("15:04") + " 3 " + clientName
			if len(club.queue) < club.tables {
				club.queue = append(club.queue, &Client{name: clientName})
			}
			club.eventLog = append(club.eventLog, outLog)
		case 4:
			clientName := parts[2]
			if client, exists := club.clients[clientName]; !exists {
				outLog := eventTime.Format("15:04") + " 13 " + "ClientUnknown"
				club.eventLog = append(club.eventLog, outLog)
			} else {
				club.calculateRevenue(client, eventTime)
				club.tableOccupancy[client.table] = nil
				delete(club.clients, clientName)

				outLog := eventTime.Format("15:04") + " 4 " + clientName
				club.eventLog = append(club.eventLog, outLog)

				if len(club.queue) != 0 {
					table, clientNameQueue := club.TakeAnEmptySeat(eventTime)
					outLog := eventTime.Format("15:04") + " 12 " + clientNameQueue + " " + strconv.Itoa(table)
					club.eventLog = append(club.eventLog, outLog)
				}

			}
		}
	}

	club.RemoveUsersAfterTime()
	club.printResults()
}
