# profile
[![Build Status](https://travis-ci.org/joaosoft/profile.svg?branch=master)](https://travis-ci.org/joaosoft/profile) | [![codecov](https://codecov.io/gh/joaosoft/profile/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/profile) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/profile)](https://goreportcard.com/report/github.com/joaosoft/profile) | [![GoDoc](https://godoc.org/github.com/joaosoft/profile?status.svg)](https://godoc.org/github.com/joaosoft/profile)

A service used on a web site for a person profile (Web Site) [github](https://github.com/joaosoft/vue-profile).


###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Get sections
* Get a section
* Get a section contents

## Endpoints
* **Get sections:** 

    Method: GET

    Route: http://localhost:9001/api/v1/profile/sections
    
    Body: 
    ```
    [
        {
            "id_section": "1",
            "key": "home",
            "name": "Home",
            "description": "Home Section",
            "active": true,
            "created_at": "2019-04-08T18:23:48.830695Z",
            "updated_at": "2019-04-08T18:23:48.830695Z"
        },
        {
            "id_section": "2",
            "key": "projects",
            "name": "Projects",
            "description": "Projects Section",
            "active": true,
            "created_at": "2019-04-08T18:23:48.830695Z",
            "updated_at": "2019-04-08T18:23:48.830695Z"
        },
        {
            "id_section": "3",
            "key": "about",
            "name": "About",
            "description": "About Section",
            "active": true,
            "created_at": "2019-04-08T18:23:48.830695Z",
            "updated_at": "2019-04-08T18:23:48.830695Z"
        }
    ]
    ```

* **Get a section:** 

    Method: PUT
    
    Route: http://localhost:9001/api/v1/profile/sections/home
    
    Body: 
    ```
    {
        "id_section": "1",
        "key": "home",
        "name": "Home",
        "description": "Home Section",
        "active": true,
        "created_at": "2019-04-08T18:23:48.830695Z",
        "updated_at": "2019-04-08T18:23:48.830695Z"
    }
    ```

* **Get a section contents:** 

    Method: PUT
    
    Route: http://localhost:9001/api/v1/profile/sections/projects/contents
    
    Body: 
    ```
    [
        {
            "key": "dbr",
            "content": {
                "title": "dbr"
            },
            "active": true,
            "created_at": "2019-04-08T18:23:48.830695Z",
            "updated_at": "2019-04-08T18:23:48.830695Z"
        },
        {
            "key": "web",
            "content": {
                "title": "web"
            },
            "active": true,
            "created_at": "2019-04-08T18:23:48.830695Z",
            "updated_at": "2019-04-08T18:23:48.830695Z"
        },
        {
            "key": "validator",
            "content": {
                "title": "validator"
            },
            "active": true,
            "created_at": "2019-04-08T18:23:48.830695Z",
            "updated_at": "2019-04-08T18:23:48.830695Z"
        }
    ]
    ```

## Dependecy Management
>### Dependency

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Get dependency manager: `go get github.com/joaosoft/dependency`
* Install dependencies: `dependency get`


>### Go
```
go get github.com/joaosoft/profile
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
