# Welcome to Call Taxi Project




# Files Struct
├── api  /
│   ├── handler.go
│   ├── server.go
│   └── status.go
├── database 
│   ├── psql.go
│   ├── redis.go
│   ├── status.go
│   └── types  
│       ├── psql.go
│       └── redis.go
├── domain   &emsp;&emsp;&emsp;//Business Logic
│   └── domain.go
├── go.mod
├── main.go
├── test    &emsp;&emsp;&emsp;&emsp;&emsp;// unit test
│   ├── database
│   │   ├── psql_test.go
│   │   └── redis_test.go
│   └── domain
│       └── domain_test.go
├── types  &emsp;&emsp;&emsp;&emsp;&emsp;// Common Struct
│   └── types.go
└── utils &emsp;&emsp;&emsp;&emsp;&emsp;// tool function
    ├── errorHandel.go
    └── utilis.go 



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

    %% 樣式設定
    style Cache fill:#fff,stroke:#ff4500,stroke-width:2px
    style DB fill:#fff,stroke:#0052cc,stroke-width:2px
    style API fill:#f9f9f9,stroke:#333,stroke-width:2px
