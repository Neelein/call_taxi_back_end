# Call Taxi プロジェクト
[English](README.md) | [日本語](ja_README.md)

## 機能
- ドライバーの位置を追跡する
- 到着予想時間と距離を表示する
- ドライバーの走行経路を表示する

## フローチャート
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

## シーケンス図

```mermaid
sequenceDiagram
    participant Cli as Client
    participant Api as API Server
    participant Cache as Cache
    participant DB as DataBase
    participant Msg as Pub/Sub
    participant Ws as WebSocket Server
    participant FKD as Driver

    Cli ->>Api: セッションキーでユーザーデータを取得
    Api ->>Cache: ユーザーデータを読み取る

    alt "Hit"
    Cache -->> Api: ユーザーデータを返す
    Api -->> Cli: ユーザーデータを返す
    else "Miss"
    Cache -->>Api: null
    Api ->>DB: 新しいユーザーデータを作成
    DB -->>Api: 新しいユーザーデータを返す
    Api ->>Cache: ユーザーデータをキャッシュに保存
    Api -->>Cli: 新しいユーザーデータを返し、セッションキーをCookieに設定
    end

    Cli ->>Api: 注文リクエスト
    Api ->>Cache: 50km以内のドライバーを検索

    alt "Hit"
    Cache -->>Api: 最も近いドライバーを返す
    Api ->> DB: 新しい注文データを作成
    DB  ->> Cache: データをキャッシュに書き込む
    Api -->> Cli: ドライバー情報を返す
    else "Miss"
    Cache-->>Api:null
    Api -->> Cli: 利用可能なドライバーがいないことを返す
    end

    loop "200 millisecond"
        FKD ->>Api: ドライバーの位置を更新
        Api ->>Cache: 最新のドライバー位置を更新
        Api ->>Msg: ドライバー位置を送信
        Msg ->>Ws: ドライバー位置を送信
        Ws -->>Cli: ドライバー位置を返す
        Cli->>Cli: 2秒ごとにドライバー位置を描画
    end
``` 

## 説明

### 1. 出発地点を選択する。
![Select Departure Location](./image/1.png)
![Select Departure Location](./image/4.png)


### 2. 目的地を選択する。
![Select destination Location](./image/5.png)
![Select destination Location](./image/6.png)

### 3. 「タクシーを呼ぶ」ボタンをクリックする。
![Select destination Location](./image/7.png)

### 4. ドライバーのステータスは「pick up」と表示される。
![Select destination Location](./image/8.png)

### 5. 地図はドライバーの位置に追従し、ルートを表示する。
![Select destination Location](./image/9.png)

### 6. ドライバーが出発地点に到着すると、ドライバーのステータスが送信される。
![Select destination Location](./image/10.png)

### 7. ドライバーが目的地に到着すると、注文ステータスは「complete」になる。
![Select destination Location](./image/2.png)

### 8. 「取り消す」ボタンをクリックすると、注文ステータスは「cancel」になる。
![Select destination Location](./image/11.png)


## 注意事項 

### 1.セッションキーはCookieに1時間保存される。1時間何も操作しない場合、Cookieは消える。「タクシーを呼ぶ」ボタンをクリックすると、「Server busy」というアラートが表示される。ページを更新してください。
### 2.「フォローを解除する」をクリックすると、地図はドライバーの追跡を解除するが、ルートは表示され続ける。

## Soruce code
[Frontend](https://github.com/Neelein/call_taxi_front_end) |
[Fake driver server](https://github.com/Neelein/call_taxi_fake_driver) |
[WebSocket server](https://github.com/Neelein/call_taxi_websocket) |
