# Web Apps in Go

In this project, I wanted to try web development.

I built 5 small servers each building on top of the next one in order to learn new concepts gradually.
  - Time Server: Tells the current time.
  - Echo Server: Echoes back the requested URL.
  - QR Server: Serves the QR code of the html form.
  - Search Server: Posts a GET request to Google CSE api to return the top 5 search results of the html form.
  - Books Server: An in memory books database, that implements RESTful API, so that the user can interact with the database by posting HTTP requests.
