# Welcome to Call Taxi Project



## Feature
- Track the driver’s location
- Show estimated arrival time and distance
- Show Driver Path

## Flowchart
``` mermaid
graph LR
    %% 節點定義
    Client([Client Side])
    LB[Load Balancer<br/>Nginx]
    API[ApiServer]
    WSS[WebSocket Server]
    PubSub[(Pub/Sub Redis)]
    Cache[(Cache Redis)]
    DB[(Database Psql)]
    FakeDr[fakeDriverServer]
    GMap{Google Map API}

    %% 接入層
    Client --- LB
    LB --- API

    %% 核心資料雙向流 (Cache-Aside 邏輯)
    API <-->|"1. Read <br>(not found read database)"| Cache
    API <-->|2 Read/Write| DB
    
    %% 強調 DB 到 Cache 的路徑
    DB -->|write| Cache 


    %% 業務邏輯
    FakeDr --> API
    FakeDr --> Cache
    FakeDr <--> GMap

    %% 推送邏輯
    API -.->|push driver location| PubSub
    PubSub -.-> WSS
    WSS -.->|stream| Client

