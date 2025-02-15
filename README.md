# pub-sub

A simple pub-sub communication model implementation from scratch in Go

## How to use

1. Clone the repository

```shell
git clone https://github.com/Kartikeya2710/pub-sub.git
```

OR

```shell
git clone git@github.com:Kartikeya2710/pub-sub.git
```

2. Run the follwing

```shell
make
```

4. Connect to the broker (default `ws://localhost:8080/ws`) via a websocket connection using Postman or any other websocket library

5. Perform the following actions by sending appropriate JSON requests

   - Subscribe:

     ```json
     {
     	"type": "subscribe",
     	"data": { "Id": "sub1", "TopicName": "news", "BufferSize": 10 }
     }
     ```

   - Publish:

     ```json
     {
     	"type": "publish",
     	"data": {
     		"TopicName": "news",
     		"Message": "Welcome to today's news headlines!!"
     	}
     }
     ```

   - Unsubscribe:

     ```json
     {
     	"type": "unsubscribe",
     	"data": { "Id": "sub1", "TopicName": "news" }
     }
     ```
