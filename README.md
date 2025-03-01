# GoExpert Multithreading

This project consists of a Go application that fetches an address based on a provided  zip code in two different APIs (ViaCep and BrasilCep) and returns the fastest response.

## How to Run
- Run the application:
    ```sh
    go run main.go [CEP]
    ```

Replace [CEP] with the CEP you want to fetch. The application will fetch the CEP in both APIs and return the fastest response. The response will be printed in the command line.
