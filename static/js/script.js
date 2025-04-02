document.addEventListener("DOMContentLoaded", function () {
    loadMeetings(); 
});

//  Meetings von API werden im Dropdown/Stellenübersicht angezeigt
function loadMeetings() {
        fetch("/meetings")
            .then(response => response.json())
            .then(meetings => {
                const select = document.getElementById("meetingSelect");
                select.innerHTML = "";
                const container = document.getElementById("meetingsContainer");
                container.innerHTML = ""; 

                meetings.forEach(meeting => {
                    const option = document.createElement("option");
                    option.value = meeting.id;
                    option.textContent = `${meeting.title} (Tickets verfügbar: ${meeting.available_tickets}/${meeting.total_tickets})`;
                    select.appendChild(option);

                    const meetingBox = document.createElement("div");
                    meetingBox.classList.add("meeting-box");
                    meetingBox.innerHTML = `
                        <h3>${meeting.title}</h3>
                        <p>Verfügbare Tickets: <strong>${meeting.available_tickets} / ${meeting.total_tickets}</strong></p>
                    `;
                    container.appendChild(meetingBox);
                });
            })
            .catch(error => console.error("Fehler beim Laden der Meetings:", error));
    }

//  Verarbeitet Buchung, wenn button gedrückt wird
document.getElementById("bookingForm").addEventListener("submit", function (event) {
    console.log("Submit-Button wurde geklickt!"); 
    event.preventDefault();

    const meetingID = parseInt(document.getElementById("meetingSelect").value);
    const firstName = document.getElementById("firstName").value.trim();
    const lastName = document.getElementById("lastName").value.trim();
    const email = document.getElementById("email").value.trim();
    const tickets = parseInt(document.getElementById("tickets").value);

    // Alle Felder müssen ausgefüllt sein
    if (!firstName || !lastName || !email || isNaN(tickets) || tickets <= 0) {
        showError("Bitte fülle alle Felder korrekt aus.");
        return;
    }

    //  Daten für die Buchung
    const bookingData = {
        meeting_id: meetingID,
        first_name: firstName,
        last_name: lastName,
        email: email,
        tickets: tickets
    };

    console.log("Daten werden übermittels",bookingData);

    fetch("/book", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(bookingData)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {
                 throw new Error(text); 
                }
            );
        }
        return response.json();
    })
    .then(data => {
        showSuccess(data.message);
        loadMeetings(loadMeetings,500); 
    })
    .catch(error => showError(error.message));
});

// Fehlermeldung 
function showError(message) {
    const responseMessage = document.getElementById("responseMessage");
    responseMessage.innerText = "❌ " + message;
    responseMessage.style.color = "red";
    responseMessage.style.opacity = "1";
    setTimeout(() => responseMessage.style.opacity = "0", 4000);
}

// Erfolgreich gebucht
function showSuccess(message) {
    const responseMessage = document.getElementById("responseMessage");
    responseMessage.innerText = "✅ " + message;
    responseMessage.style.color = "green";
    responseMessage.style.opacity = "1";
    setTimeout(() => responseMessage.style.opacity = "0", 4000);
}

//  Smooth Animationen für Bilder und Inhalte
document.addEventListener("DOMContentLoaded", function() {
    const fadeElements = document.querySelectorAll(".fade-in");

    function handleScroll() {
        fadeElements.forEach(el => {
            const rect = el.getBoundingClientRect();
            if (rect.top < window.innerHeight * 0.8) {
                el.classList.add("show");
            }
        });
    }

    window.addEventListener("scroll", handleScroll);
    handleScroll();
});
