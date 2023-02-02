# Stay4Long Challenge Code
[![CI Test](https://github.com/xSolrac87/stay4long/actions/workflows/test.yml/badge.svg)](https://github.com/xSolrac87/stay4long/actions/workflows/test.yml)

This challenge was about create an API with a unique namespace ( booking ) and under that with two different endpoints.  
````
POST /stats
POST /maximize
with the payload body being an slice of bookings
[
  {
    "request_id": "acme_AAAA",
    "check_in": "YYYY-MM-DD",
    "nights": 3,
    "selling_rate": 100,
    "margin": 20
  }
]
````

## My Approach
I implemented this solution purely with go language and following an architecture know as package pattern, where  
each one of them is responsible for one unique task, which make sure that the code is clear, easy to understand, and reusable. 
On top of that this project is made build in house which means that all the package that I'm using are from standard library,   
no extra dependencies are require.  
With this I'm able to create a clear, organized, and maintainable solution.

## Prerequisites
- Docker
- Docker-Compose

## Optional
- Go >= 1.19
- Make

## Getting Started
There is an .env.example file where you can configure ENV variables such as `SERVER_PORT` and `SERVER_TIMEOUT`.  
If there is no .env file, by default docker-compose use 7546 as `SERVER_PORT` and 60 as `SERVER_TIMEOUT`.  
In the docker-compose file `PORT: 3214` is used for port mapping, you can change it if you need.  
To start the server:
   - #### with docker-compose:
        ```bash
        docker-compose up
        ```
   - #### with make:
      ```bash
      make up
        ```
   - #### with go (needs ENV variables):
     ```bash
     go run cmd/api/main.go
     ```

## How it works
Once the container is running, if you started the container either `docker-compose` or `make` you can access it by (assuming default settings):  
``http://localhost:7546/``  

If you decided to start your server at local with go, it will be the same but changing the port to use the one on the ENV variables.

There are two endpoints available:
### stats
    Valid HTTP Method: POST
    Endpoint: http://localhost:7546/stats
    Body: Slice of Bookings
    Status Code: 200, 400, 405 and 500
    Response: {"avg_night":<float>,"min_night":<float>,"max_night":<float>}
    Example:
  ```bash
    curl -X POST \
    http://localhost:7546/stats \
    -H 'Content-Type: application/json' \
    -d '[
            {
            "request_id": "1234567890",
            "check_in": "2022-10-01",
            "nights": 7,
            "selling_rate": 200,
            "margin": 8
            },
            {
            "request_id": "0987654321",
            "check_in": "2022-10-08",
            "nights": 5,
            "selling_rate": 150,
            "margin": 7
            }
    ]'
  ```
### maximize
    Valid HTTP Method: POST
    Endpoint: http://localhost:7546/stats
    Body: Slice of Bookings
    Status Code: 200, 400, 405 and 500
    Response: {"request_ids":<array>,"total_profit":<int>,"avg_night":<float>,"min_night":<float>,"max_night":<float>}
    Example:
  ```bash
    curl -X POST \
    http://localhost:7546/maximize \
    -H 'Content-Type: application/json' \
    -d '[
            {
            "request_id": "1234567890",
            "check_in": "2022-10-01",
            "nights": 7,
            "selling_rate": 200,
            "margin": 8
            },
            {
            "request_id": "0987654321",
            "check_in": "2022-10-08",
            "nights": 5,
            "selling_rate": 150,
            "margin": 7
            }
    ]'
  ```

## Testing
### UNIT && E2E
- #### with docker-compose ( assuming your docker container is the same ):
     ```bash
     docker exec -it s4l_s4l-service-api_1 go test ./... -tags=unit -v
     ```
    ```bash
    docker exec -it s4l_s4l-service-api_1 go test ./... -tags=e2e -v
     ```
- #### with make:
   ```bash
    make test-unit
   ```
     ```bash
    make test-e2e
   ```
- #### with go:
  ```bash
  go test ./... -tags=unit -v
  ```
    ```bash
  go test ./... -tags=e2e -v
  ```

## Benchmark
I've added benchmark functions to `MaximTotalProfits` and `ProfitPerNight`. This last one included an n variable to create  
that number of bookings, I left it with 1M ( can be change on your own ). ( booking/stats_service_test.go:115 )
- #### with docker-compose ( assuming your docker container is the same ):
     ```bash
     docker exec -it s4l_s4l-service-api_1 go test ./... -bench=. -tags=unit
     ```
- #### with make:
   ```bash
    make benchmark
   ```
- #### with go:
  ```bash
  go test ./... -bench=. -tags=unit
  ```

## Fuzzing
I have also added fuzzing test (stress test) to profit logic, with the appropriate time or for production I'd add fuzzing test to HTTP handler.
- #### with docker-compose ( assuming your docker container is the same ):
     ```bash
     docker exec -it s4l_s4l-service-api_1 go test ./booking/ --fuzz=FuzzPerNight -fuzztime=10s -tags=unit
     ```
    ```bash
    docker exec -it s4l_s4l-service-api_1 go test ./booking/ --fuzz=FuzzProfit -fuzztime=10s -tags=unit
     ```
- #### with make:
   ```bash
    make fuzzing-per-night
   ```
     ```bash
    make fuzzing-profit
   ```
- #### with go:
  ```bash
  go test ./booking/ --fuzz=FuzzPerNight -fuzztime=10s -tags=unit
  ```
    ```bash
  go test ./booking/ --fuzz=FuzzProfit -fuzztime=10s -tags=unit
  ```

## Conclusion
Topics that I eager to discuss with stay4long team, to see/know their approach.
### Dealing with floats
I Knew from the beginning that I'd have problems comparing float numbers,  
since comparing the result of floating-point calculations depends on the actual processor.
For that reason a put some comments on unit testing saying that one solution could be Comparing two values using a delta.
### Maximize endpoint
At first, I use a brute force solution, fora nested loop
```
for j := i + 1; j < len(bookings); j++ {
	b := bookings[j]
	if m.noOverlap(booking, b) {
		combinations[booking.RequestID] = append(combinations[booking.RequestID], b)
		checked[b.RequestID] = struct{}{}
	}
}
```
Even the benchmark wasn't that bad, I was not happy with it, so I decided to change it for the current one,
even I know is not perfect it improves around 38% and has better allocations.
With the proper time I'd have implemented a binary search to reduce complexity to 0(n)

