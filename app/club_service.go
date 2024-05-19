package app

import (
	"fmt"
	"time"
)

type Club struct {
	tables            int
	openTime          time.Time
	closeTime         time.Time
	pricePerHour      int
	clients           map[string]*Client
	tableOccupancy    map[int]*Client
	eventLog          []string
	revenue           map[int]int
	totalTimeOccupied map[int]time.Duration
	queue             []*Client
}

type Client struct {
	name      string
	startTime time.Time
	endTime   time.Time
	table     int
}

func (club *Club) calculateRevenue(client *Client, endTime time.Time) {
	duration := endTime.Sub(client.startTime)
	hours := int(duration.Hours())

	if duration.Minutes() > float64(hours*60) {
		hours++
	}

	club.revenue[client.table] += hours * club.pricePerHour
	club.totalTimeOccupied[client.table] += duration
}

func (club *Club) printResults() {
	fmt.Printf("%02d:%02d\n", club.openTime.Hour(), club.openTime.Minute())
	for _, event := range club.eventLog {
		fmt.Println(event)
	}
	fmt.Printf("%02d:%02d\n", club.closeTime.Hour(), club.closeTime.Minute())

	for table := 1; table <= club.tables; table++ {
		fmt.Printf("%d %d %s\n", table, club.revenue[table], club.convertTime(club.totalTimeOccupied[table]))
	}
}

func (club *Club) convertTime(time time.Duration) string {
	hours := int(time.Hours())
	minutes := int(time.Minutes()) % 60

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func (club *Club) RemoveUsersAfterTime() {
	for client := range club.clients {
		club.calculateRevenue(club.clients[client], club.closeTime)
		delete(club.clients, client)
		outLog := club.closeTime.Format(TemplateFormatTime) + " 11 " + client
		club.eventLog = append(club.eventLog, outLog)
	}
}

func (club *Club) TakeAnEmptySeat(eventTime time.Time) (int, string) {
	for table := range club.tableOccupancy {
		if club.tableOccupancy[table] == nil {
			clientName := club.queue[0].name

			club.queue[0].startTime = eventTime
			club.queue[0].table = table

			club.clients[clientName] = club.queue[0]
			club.tableOccupancy[table] = club.queue[0]
			club.queue = club.queue[1:]
			return table, clientName
		}
	}
	return 0, ""
}
