# Go Client-Server Application

This Go project consists of a simple client-server application to fetch and store the current USD exchange rate.

## Description

- The **server** retrieves the USD exchange rate from an external API and responds with the current rate of the dollar to the Brazilian real. This value is exposed as JSON via an HTTP endpoint.
- The **client** makes a request to the server, obtains the exchange rate value, and saves it in a file called `cotacao.txt`, in the format `Dólar: {value}`.

## Features

- **Server**:
  - Sends a request to an external API to get the USD exchange rate.
  - Returns the exchange rate in JSON format via an HTTP endpoint.
  - Persists exchange rate data in an SQLite database.

- **Client**:
  - Requests the server to retrieve the exchange rate.
  - Saves the exchange rate in the file `cotacao.txt` for local storage.

## Technologies

- Language: Go
- Database: SQLite (for server-side persistence of exchange rates)
- HTTP: Client-server communication via REST API

Let's #GO
