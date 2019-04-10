# profile
[![Build Status](https://travis-ci.org/joaosoft/profile.svg?branch=master)](https://travis-ci.org/joaosoft/profile) | [![codecov](https://codecov.io/gh/joaosoft/profile/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/profile) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/profile)](https://goreportcard.com/report/github.com/joaosoft/profile) | [![GoDoc](https://godoc.org/github.com/joaosoft/profile?status.svg)](https://godoc.org/github.com/joaosoft/profile)

A service used on a web site for a person profile (Web Site) [github](https://github.com/joaosoft/vue-profile).


###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Get sections
* Get sections with contents
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
        },
        {
            "id_section": "2",
            "key": "projects",
            "name": "Projects",
            "description": "Projects Section",
        },
        {
            "id_section": "3",
            "key": "about",
            "name": "About",
            "description": "About Section",
        }
    ]
    ```

* **Get sections with contents:** 

    Method: GET

    Route: http://localhost:9001/api/v1/profile/sections/contents
    
    Body: 
    ```
    [
        {
            "id_section": "1",
            "key": "home",
            "name": "Home",
            "description": "Home Section",
            "contents": []
        },
        {
            "id_section": "2",
            "key": "projects",
            "name": "Projects",
            "description": "Projects Section",
            "contents": [
                {
                    "key": "dbr",
                    "content": {
                        "url": "https://github.com/joaosoft/dbr",
                        "build": "https://travis-ci.org/joaosoft/dbr.svg?branch=master",
                        "title": "Dbr",
                        "description": "A simple database client with support for master/slave databases."
                    }
                },
                {
                    "key": "web",
                    "content": {
                        "url": "https://github.com/joaosoft/web",
                        "build": "https://travis-ci.org/joaosoft/web.svg?branch=master",
                        "title": "Web",
                        "description": "A simple and fast web server and client."
                    }
                },
                {
                    "key": "validator",
                    "content": {
                        "url": "https://github.com/joaosoft/validator",
                        "build": "https://travis-ci.org/joaosoft/validator.svg?branch=master",
                        "title": "Validator",
                        "description": "A simple struct validator by tags."
                    }
                }
            ]
        },
        {
            "id_section": "3",
            "key": "about",
            "name": "About",
            "description": "About Section",
            "contents": []
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
        },
        {
            "key": "web",
            "content": {
                "title": "web"
            },
        },
        {
            "key": "validator",
            "content": {
                "title": "validator"
            },
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
