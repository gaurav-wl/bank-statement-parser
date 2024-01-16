# Bank Statement Parser Microservice

This microservice is designed to parse bank statements and deduce the total value of payments made on a specific date for each currency.

## Getting Started

Follow the instructions below to run the server and interact with the API.

### Prerequisites

- Go installed on your machine: [Download Go](https://golang.org/dl/)
- Make sure you have `curl` installed for making HTTP requests or you can use postman.

### Setup

1. Clone the repository:

    ```bash
    git clone https://github.com/gaurav-wl/bank-statement-parser.git
    cd bank-statement-parser
    ```

2. Build the application:

    ```bash
    makefile build
    ```

3. Run the server:

    ```bash
    makefile run PORT=8080
    ```

The server will be running at `http://localhost:8080`. or at the provided port.

### API Endpoint:

The API endpoint for parsing bank statements is /parse, and it expects a CSV file as form data and date to filter out the result.

#### Request:
- Method: POST
- Endpoint: /parse
- Content Type: multipart/form-data

#### Request Parameters
- file (type: file): The CSV file containing the bank statement.
- date (type: string, format: "dd/MM/yyyy"): The date for which payments will be considered. (e.g., "06/03/2011")

Example Request:
```bash
curl -X POST \
  -F "file=<your_file.csv>" \
  -F "date=06/03/2011" \
  http://localhost:8080/parse
```

Replace <your_file.csv> with the actual path to your CSV file.

#### Example Response:
```json
[
{"currency": "EUR", "total": 43543.23},
{"currency": "GBP", "total": 140263.45},
{"currency": "CAD", "total": 52789.21},
{"currency": "USD", "total": 126445.34}
]
```

The response will contain the total value of payments made on the specified date for each currency.

#### Notes

- Make sure to replace placeholders like <your_file.csv> with actual values.
- Adjust the port and other parameters in the examples as needed.
- Ensure that the server is running (make run PORT=8080) before making API requests.




