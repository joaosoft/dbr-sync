# session
[![Build Status](https://travis-ci.org/joaosoft/session.svg?branch=master)](https://travis-ci.org/joaosoft/session) | [![codecov](https://codecov.io/gh/joaosoft/session/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/session) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/session)](https://goreportcard.com/report/github.com/joaosoft/session) | [![GoDoc](https://godoc.org/github.com/joaosoft/session?status.svg)](https://godoc.org/github.com/joaosoft/session)

A service that allows you to get a new session token and when invalid refresh the session token (jwt).


###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Get session
* Refresh session with refresh token

## Endpoints
* **Get Session:** 

    Method: GET

    Route: http://localhost:8001/api/v1/get-session
    
    Body: 
    ```
    {
        "email": "joaosoft@gmail.com",
        "password": "698dc19d489c4e4db73e28a713eab07b"
    }
    ```

* **Refresh session:** 

    Method: PUT
    
    Route: http://localhost:8001/api/v1/get-session
    
    Headers:
    ```
    Authorization: Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzZXNzaW9uIiwiaWRfdXNlciI6IjEiLCJqdGkiOiI5MWY4MDBhZS00OGQ1LTQzMmUtOWYwZC0xYzAzODY1YmMyZjciLCJzdWIiOiJyZWZyZXNoLXRva2VuIn0.Yfvxgkw3NNqkF9nDuMymp-L0dN9j6vozdeU-A3JTQPc86FGfKeQRSI3CZOaGWZ_Q
    ```

## Dependecy Management
>### Dependency

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Get dependency manager: `go get github.com/joaosoft/dependency`
* Install dependencies: `dependency get`


>### Go
```
go get github.com/joaosoft/session
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
