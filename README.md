# Welcome to Call Taxi Project
[English](README.md) | [日本語](ja_README.md)

## Feature
- Track the driver’s location
- Show estimated arrival time and distance
- Show Driver Path

## Flowchart
``` mermaid
graph TD
    Cli(Cleint)
    DB[(DataBase)]
    Cache[(cache)]
    Messagep([Pub/Sub])
    Api(apiServer)
    LB(LoadBalnce)
    FKD(Fake Driver Server)
    Ws(WebSocketServer)
    Gm(Google Map Api)
    
    Cli <--> LB
    LB <--> Api
    Api <--"Write/Read" --> DB
    DB -- "Write" --> Cache
    Api <-- "Read" --> Cache
    FKD <--> Api
    FKD <-->Gm
    Api -- "Write driver position"-->Cache
    Api -- Send driver position--> Messagep -->Ws --> Cli
```

## Sequence Diagram

```mermaid
sequenceDiagram
    participant Cli as Client
    participant Api as API Server
    participant Cache as Cache
    participant DB as DataBase
    participant Msg as Pub/Sub
    participant Ws as WebSocket Server
    participant FKD as Driver

    Cli ->>Api:Get user data by session key
    Api ->>Cache:read user data

    alt "Hit"
    Cache -->> Api:retrun user data
    Api -->> Cli:response user data 
    else "Miss"
    Cache -->>Api:null
    Api ->>DB:create new user data
    DB -->>Api:return new user data
    Api ->>Cache:write user data to cache
    Api -->>Cli:response new user data and set session key in cookie
    end

    Cli ->>Api:order request
    Api ->>Cache:find drivers within 50 km

    alt "Hit"
    Cache -->>Api:return nearest driver
    Api ->> DB:create new order data
    DB  ->> Cache:write data to cahce
    Api -->> Cli:response driver data
    else "Miss"
    Cache-->>Api:null
    Api -->> Cli:response no availiable driver
    end

    loop "200 millisecond"
        FKD ->>Api:update driver position
        Api ->>Cache:update latest driver position
        Api ->>Msg:push driver position
        Msg ->>Ws:send driver position
        Ws -->>Cli:response driver position
        Cli->>Cli:render new driver position every 2 second
    end
``` 

## Instructions

### 1.Select departure location.
![Select Departure Location](./image/1.png)
![Select Departure Location](./image/4.png)

### 2.Select destination location.
![Select destination Location](./image/5.png)
![Select destination Location](./image/6.png)

### 3.Click call texi button.
![Select destination Location](./image/7.png)

### 4.The driver status will  show pick up.
![Select destination Location](./image/8.png)

### 5.The map will follow the driver’s location and show the route.
![Select destination Location](./image/9.png)

### 6.When the driver arrives at the departure location, the driver status will be "send".
![Select destination Location](./image/10.png)

### 7.When the driver arrives at the destination location, the order status will be "complete".
![Select destination Location](./image/2.png)

### 8.Click cancel button,the order status will be "cancel".
![Select destination Location](./image/11.png)

### Notice 
#### 1.The session key will be stored in a cookie for one hour. If you are inactive for one hour, the cookie will expire. If you click the “Call Taxi” button, a “Server busy” alert will be displayed. Please refresh the page.

#### 2.Click “Unfollow Driver” to stop following the driver on the map, while the route will remain displayed.

## Soruce code
[Frontend](https://github.com/Neelein/call_taxi_front_end) |
[Fake driver server](https://github.com/Neelein/call_taxi_fake_driver) |
[WebSocket server](https://github.com/Neelein/call_taxi_websocket) |
