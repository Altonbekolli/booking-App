package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Struktur für Meetings
type Meeting struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	TotalTickets     uint   `json:"total_tickets"`
	AvailableTickets uint   `json:"available_tickets"`
}

// Struktur für eine Buchung
type Booking struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Tickets   uint   `json:"tickets"`
	MeetingID int    `json:"meeting_id"`
}

var meetings = []Meeting{
	{ID: 1, Title: "Go Conference 2024", TotalTickets: 30, AvailableTickets: 30},
	{ID: 2, Title: "Go Web Development Workshop", TotalTickets: 40, AvailableTickets: 40},
	{ID: 3, Title: "Go Founders Meetup", TotalTickets: 7, AvailableTickets: 7},
	{ID: 4, Title: "Go Performance Bootcamp", TotalTickets: 100, AvailableTickets: 100},
	{ID: 5, Title: "Go AI & Machine Learning Meetup", TotalTickets: 75, AvailableTickets: 75},
	{ID: 6, Title: "Go Open Source Contributors Gathering", TotalTickets: 1000, AvailableTickets: 1000},
}

var bookings = []Booking{}
var mu sync.Mutex

func main() {
	fmt.Println("Starting server on http://localhost:8080...")

	//  HTML, CSS & JS als statische Dateien bereitstellen
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	//  API-Routen
	http.HandleFunc("/meetings", getMeetings)         // Meetings abrufen
	http.HandleFunc("/create-meeting", createMeeting) // Meeting erstellen (Admin)
	http.HandleFunc("/book", handleBooking)           // Tickets buchen

	//  Server starten
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Fehler beim Starten des Servers:", err)
	}
}

// Zeigt alle Meetings an (GET)
func getMeetings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meetings)
}

// Admin kann ein neues Meeting erstellen (POST)
func createMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var newMeeting Meeting
	err := json.NewDecoder(r.Body).Decode(&newMeeting)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newMeeting.ID = len(meetings) + 1
	newMeeting.AvailableTickets = newMeeting.TotalTickets
	meetings = append(meetings, newMeeting)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Meeting created successfully!"})
}

// Verarbeitet Buchungen (POST)
func handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var bookingData Booking
	err := json.NewDecoder(r.Body).Decode(&bookingData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//  Prüfen, ob Meeting existiert
	var selectedMeeting *Meeting
	for i := range meetings {
		if meetings[i].ID == bookingData.MeetingID {
			selectedMeeting = &meetings[i]
			break
		}
	}

	if selectedMeeting == nil {
		http.Error(w, "Meeting not found", http.StatusNotFound)
		return
	}

	//  Validierung
	isValidEmail := strings.Contains(bookingData.Email, "@")
	isValidName := len(bookingData.FirstName) >= 2 && len(bookingData.LastName) >= 2
	isValidTicket := bookingData.Tickets > 0 && bookingData.Tickets <= selectedMeeting.AvailableTickets

	if !isValidEmail || !isValidName || !isValidTicket {
		http.Error(w, "Invalid booking data", http.StatusBadRequest)
		return
	}

	//  Buchung speichern
	mu.Lock()
	selectedMeeting.AvailableTickets -= bookingData.Tickets
	bookings = append(bookings, bookingData)
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"message": "Booking successful!"})
}
