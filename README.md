# parking-lot-system
## Run aplication
```
go run cmd/server/main.go
```
## 1. Park Vehicle
URL: ``` http://localhost:8080/park ```
Request Body:
```json
{
  "vehicleType": "Bicycle",
  "vehicleNumber": "BC001"
}
```
or using cURL:
```curl
curl -X POST http://localhost:8080/park \
     -H "Content-Type: application/json" \
     -d '{"vehicleType": "Bicycle", "vehicleNumber": "BC001"}'
```

## 2. Unpark Vehicle
URL: ``` http://localhost:8080/unpark ```
Request Body:
```json
{
  "spotId": "0-0-1",
  "vehicleNumber": "BC001"
}
```
or using cURL:
```curl
curl -X POST http://localhost:8080/unpark \
     -H "Content-Type: application/json" \
     -d '{"spotId": "0-0-1", "vehicleNumber": "BC001"}'
```

## 3. Available Spot
cURL:
```curl
curl -X GET "http://localhost:8080/available?vehicleType=Bicycle"
```

## 4. Search Vehicle
cURL:
```curl
curl -X GET "http://localhost:8080/search?vehicleNumber=BC001"
```
