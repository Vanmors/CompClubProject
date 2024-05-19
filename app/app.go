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
		fmt.Println("Input format error: " + scanner.Text())
		os.Exit(1)
	}

	scanner.Scan()
	times := strings.Split(scanner.Text(), " ")
	club.openTime, err = time.Parse(TemplateFormatTime, times[0])
	if err != nil {
		fmt.Println("Input format error: " + times[0] + " " + times[1])
		os.Exit(1)
	}
	club.closeTime, err = time.Parse(TemplateFormatTime, times[1])
	if err != nil {
		fmt.Println("Input format error: " + times[0] + " " + times[1])
		os.Exit(1)
	}

	scanner.Scan()
	club.pricePerHour, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Input format error: " + scanner.Text())
		os.Exit(1)
	}

	for scanner.Scan() {
		events = append(events, scanner.Text())
	}

	for _, event := range events {
		parts := strings.Split(event, " ")
		eventTime, err := time.Parse(TemplateFormatTime, parts[0])
		if err != nil {
			fmt.Println("Input format error: " + event)
			os.Exit(1)
		}
		eventId, err := strconv.Atoi(parts[1])
		if err != nil || eventId > 4 || eventId < 1 {
			fmt.Println("Input format error: " + event)
			os.Exit(1)
		}
		clientName := parts[2]

		if eventTime.After(club.closeTime) {
			club.RemoveUsersAfterTime()
		}

		switch eventId {
		case 1:
			outLog := eventTime.Format(TemplateFormatTime) + " 1 " + clientName
			club.eventLog = append(club.eventLog, outLog)
			if _, exists := club.clients[clientName]; exists {
				outLog := eventTime.Format(TemplateFormatTime) + " 13 " + "YouShallNotPass"
				club.eventLog = append(club.eventLog, outLog)
			} else if eventTime.Before(club.openTime) || eventTime.After(club.closeTime) {
				outLog = eventTime.Format(TemplateFormatTime) + " 13 " + "NotOpenYet"
				club.eventLog = append(club.eventLog, outLog)
			} else {
				club.clients[clientName] = &Client{name: clientName}
			}
		case 2:
			tableNumber, _ := strconv.Atoi(parts[3])
			if tableNumber > club.tables {
				fmt.Println(event)
				os.Exit(1)
			}
			outLog := eventTime.Format(TemplateFormatTime) + " 2 " + clientName + " " + strconv.Itoa(tableNumber)
			club.eventLog = append(club.eventLog, outLog)
			if client, exists := club.clients[clientName]; !exists {
				outLog := eventTime.Format(TemplateFormatTime) + " 13 " + "ClientUnknown"
				club.eventLog = append(club.eventLog, outLog)
			} else if currentClient, occupied := club.tableOccupancy[tableNumber]; occupied && currentClient.name != clientName {
				outLog := eventTime.Format(TemplateFormatTime) + " 13 " + "PlaceIsBusy"
				club.eventLog = append(club.eventLog, outLog)
			} else {
				if client.table != 0 {
					club.calculateRevenue(client, eventTime)
					club.tableOccupancy[client.table] = nil
					club.tableOccupancy[tableNumber] = client
					client.table = tableNumber
				}
				client.table = tableNumber
				client.startTime = eventTime
				club.tableOccupancy[tableNumber] = client
			}
		case 3:
			outLog := eventTime.Format(TemplateFormatTime) + " 3 " + clientName
			club.eventLog = append(club.eventLog, outLog)
			if len(club.tableOccupancy) < club.tables {
				outLog := eventTime.Format(TemplateFormatTime) + " 13 " + "ICanWaitNoLonger!"
				club.eventLog = append(club.eventLog, outLog)
				continue
			}
			if len(club.queue) < club.tables {
				club.queue = append(club.queue, &Client{name: clientName})
			} else {
				outLog := eventTime.Format(TemplateFormatTime) + " 11 " + clientName
				club.eventLog = append(club.eventLog, outLog)
			}

		case 4:
			clientName := parts[2]
			outLog := eventTime.Format(TemplateFormatTime) + " 4 " + clientName
			club.eventLog = append(club.eventLog, outLog)
			if client, exists := club.clients[clientName]; !exists {
				outLog := eventTime.Format(TemplateFormatTime) + " 13 " + "ClientUnknown"
				club.eventLog = append(club.eventLog, outLog)
			} else {
				club.calculateRevenue(client, eventTime)
				club.tableOccupancy[client.table] = nil
				delete(club.clients, clientName)

				if len(club.queue) != 0 {
					table, clientNameQueue := club.TakeAnEmptySeat(eventTime)
					outLog := eventTime.Format(TemplateFormatTime) + " 12 " + clientNameQueue + " " + strconv.Itoa(table)
					club.eventLog = append(club.eventLog, outLog)
				}

			}
		}
	}

	club.RemoveUsersAfterTime()
	club.printResults()
}
