# session
[![Build Status](https://travis-ci.org/joaosoft/session.svg?branch=master)](https://travis-ci.org/joaosoft/session) | [![codecov](https://codecov.io/gh/joaosoft/session/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/session) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/session)](https://goreportcard.com/report/github.com/joaosoft/session) | [![GoDoc](https://godoc.org/github.com/joaosoft/session?status.svg)](https://godoc.org/github.com/joaosoft/session)

A service that allows you to get a new session token and when invalid refresh the session token (WST, Web Security Token) [github](https://github.com/joaosoft/auth-types/wst).


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
    Authorization: Bearer 53464673625464434c557442584467316144776f4d574530535842475258466e54444e624c6a4230504352495877.53464673625542424d45386d5755596f53305532516d7841624455764d46786c505439615a30596b52567368543263775a44677562554d7a4e4338324d31737564575642545542505354466a556964455155315163577378597a633554533953586c56755155306a576b74424d6d34745654466f4a6a704f51564a6a627a67734a7935744f5377684a47396e51564d73633239474b4751755355526c6158416c4c436777.4f56737a4e57417350436c71534370695655516a62304a6e546d6f725657413551445976526d676c4a6b52725845684e625678524e6b6b32624664304e4367775a5759744c6a5a796231417662536b6f
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
