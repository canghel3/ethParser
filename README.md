### Ethereum blockchain parser


<b>INSTALLATION</b>
```go
go get github.com/canghel3/ethParser
```

<b>USAGE</b>

```go
package main

import "github.com/canghel3/ethParser/api"

func main() {
	server := api.NewServer()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
```

<b>ENDPOINTS</b>

### /subscribe
Adds the provided address to the list of subscribers for transaction monitoring.
```
curl --location --request POST 'http://localhost:1234/subscribe?address=subscriber_address'
```

### /block
Retrieves the current ethereum block as an integer.
```
curl --location 'http://localhost:1234/block'
```

### /transactions
Retrieve the transactions of the provided address.
```
curl --location 'http://localhost:1234/transactions?address=subcriber_address'
```

