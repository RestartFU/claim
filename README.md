# claim
Claim library for Dragonfly.

## Getting started
```
go get github.com/dragonfly-on-steroids/claim
```

## Usage
Usage of the Claim library:
```go
// When accepting a player, you need to give them a claim.ClaimHandler
    for{
	p, err := server.Accept(){
	if err !=nil{
            return
	}
	// You may use a library, so you can have multiple handlers.
        p.Handle(claim.NewClaimHandler(p, loader))
	}
}
// Let's say our claim area is in between 0,0 and 10,10

```

## Creating a loader
```go
// You may use a loader provided by the library.
// For Example:
db, _ := sql.Open("sqlite3", "./claims.db")
// The SQL loader requires a *sql.DB and a claim.Handler (may be nil)
loader := loaders.NewSQL(db, nil)
```
