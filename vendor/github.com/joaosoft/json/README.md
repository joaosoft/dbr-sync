json
================

[![Build Status](https://travis-ci.org/joaosoft/json.svg?branch=master)](https://travis-ci.org/joaosoft/json) | [![codecov](https://codecov.io/gh/joaosoft/json/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/json) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/json)](https://goreportcard.com/report/github.com/joaosoft/json) | [![GoDoc](https://godoc.org/github.com/joaosoft/json?status.svg)](https://godoc.org/github.com/joaosoft/json)

A simple json marshal and unmarshal by customized tags (exported fields only).

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Customized tags
* Ignore field using tag "-"

## Dependecy Management
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/json
```

## Usage 
This examples are available in the project at [json/examples](https://github.com/joaosoft/json/tree/master/examples)

### Code
```go
type address struct {
	Ports  []int             `db.read:"ports"`
	Street string            `db.read:"street"`
	Number float64           `db.write:"number"`
	Map    map[string]string `db:"map"`
}

type person struct {
	Name      string              `db:"name"`
	Age       int                 `db:"age"`
	Address   *address            `db:"address"`
	Numbers   []int               `db:"numbers"`
	Others    map[string]*address `db:"others"`
	Addresses []*address          `db:"addresses"`
}

type contents []content

type content struct {
	Data *j.RawMessage `db:"data"`
}

func main() {
	marshal()
	unmarshal()
}

func marshal() {
	fmt.Println("\n\n:: MARSHAL")

	marshal_example_1()
	marshal_example_2()

}

func unmarshal() {
	fmt.Println("\n\n:: UNMARSHAL")

	unmarshal_example_1()
	unmarshal_example_2()
	unmarshal_example_3()
	unmarshal_example_4()
	unmarshal_example_5()
	unmarshal_example_6()
	unmarshal_example_7()
}

func marshal_example_1() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func marshal_example_2() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.write"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func unmarshal_example_1() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.read"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// unmarshal
	var newExample person
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n:: Example: %+v", newExample)
	fmt.Printf("\n:: Address: %+v\n\n\n", newExample.Address)
}

func unmarshal_example_2() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	// with tags "db" and "db.write"
	// marshal
	bytes, err := json.Marshal(example, "db", "db.write")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	// unmarshal
	newExample := person{}
	err = json.Unmarshal(bytes, &newExample, "db", "db.write")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: Example: %+v", newExample)
	fmt.Printf("\n:: Address: %+v", newExample.Address)

	for key, value := range newExample.Others {
		fmt.Printf("\n:: Others Key: %+v", key)
		fmt.Printf("\n:: Others Value: %+v", value)
	}

	for _, value := range newExample.Addresses {
		fmt.Printf("\n:: Addresses: %+v", value)
	}
}

func unmarshal_example_3() {
	addr := &address{
		Street: "street one",
		Number: 1.2,
		Map:    map[string]string{`"ola" "joao"`: `"adeus" "joao"`, "c": "d"},
	}

	example := person{
		Name:      "joao",
		Age:       30,
		Address:   addr,
		Numbers:   []int{1, 2, 3},
		Others:    map[string]*address{`"ola" "joao"`: addr, "c": addr},
		Addresses: []*address{addr, addr},
	}

	persons := []*person{&example, &example}
	bytes, err := json.Marshal(persons, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newPersons []*person
	err = json.Unmarshal(bytes, &newPersons, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newPersons))
	fmt.Printf("\n:: Example 1: %+v", newPersons[0])
	fmt.Printf("\n:: Example 1 Address: %+v", newPersons[0].Address)
	fmt.Printf("\n:: Example 2: %+v", newPersons[1])
}

func unmarshal_example_4() {
	example := []int{1, 2, 3}

	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newExample []int
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample)
}

func unmarshal_example_5() {
	example := []int{}

	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newExample []int
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample)
}

func unmarshal_example_6() {
	example := map[string]int{"one": 1, "two": 2, "three": 3}

	bytes, err := json.Marshal(example, "db", "db.read")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %s", string(bytes))

	// unmarshal
	var newExample map[string]int
	err = json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample)
}

func unmarshal_example_7() {
	bytes := []byte(`[{"data":{"test": "one", "test": "two"}}]`)

	// unmarshal
	var newExample contents
	err := json.Unmarshal(bytes, &newExample, "db", "db.read")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n:: LEN: %d", len(newExample))
	fmt.Printf("\n:: Example: %+v", newExample[0].Data)
}
```

> ##### Result:
```go
:: MARSHAL
{"name":"joao","age":30,"address":{"ports":null,"street":"street one","map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}},"addresses":[{"ports":null,"street":"street one","map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}},{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}]}
{"name":"joao","age":30,"address":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"number":1.2,"map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}},"c":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}},"addresses":[{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},{"number":1.2,"map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}}]}


:: UNMARSHAL
{"name":"joao","age":30,"address":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"c":{"ports":null,"street":"street one","map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}},"\"ola\" \"joao\"":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}},"addresses":[{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},{"ports":null,"street":"street one","map":{"c":"d","\"ola\" \"joao\"":"\"adeus\" \"joao\""}}]}

:: Example: {Name:joao Age:30 Address:0xc000086200 Numbers:[1 2 3] Others:map[c:0xc000086300 "ola" "joao":0xc000086340 addresses:0xc000086380] Addresses:[0xc0000863c0 0xc000086400]}
:: Address: &{Ports:[] Street:street one Number:0 Map:map["ola" "joao":"adeus" "joao" c:d]}


{"name":"joao","age":30,"address":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}},"addresses":[{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},{"number":1.2,"map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}]}

:: Example: {Name:joao Age:30 Address:0xc000086500 Numbers:[1 2 3] Others:map[addresses:0xc0000865c0 "ola" "joao":0xc000086600 c:0xc000086640] Addresses:[0xc000086680 0xc0000866c0]}
:: Address: &{Ports:[] Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Others Key: addresses
:: Others Value: &{Ports:[] Street: Number:1.2 Map:map[c:d "ola" "joao":"adeus" "joao"]}
:: Others Key: "ola" "joao"
:: Others Value: &{Ports:[] Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Others Key: c
:: Others Value: &{Ports:[] Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Addresses: &{Ports:[] Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Addresses: &{Ports:[] Street: Number:1.2 Map:map["ola" "joao":"adeus" "joao" c:d]}

 [{"name":"joao","age":30,"address":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}},"addresses":[{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}]},{"name":"joao","age":30,"address":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"numbers":[1,2,3],"others":{"\"ola\" \"joao\"":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},"c":{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}},"addresses":[{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}},{"ports":null,"street":"street one","map":{"\"ola\" \"joao\"":"\"adeus\" \"joao\"","c":"d"}}]}]
:: LEN: 2
:: Example 1: &{Name:joao Age:30 Address:0xc000086940 Numbers:[1 2 3] Others:map[addresses:0xc000086a40 "ola" "joao":0xc000086a80 c:0xc000086ac0] Addresses:[0xc000086b00 0xc000086b40]}
:: Example 1 Address: &{Ports:[] Street:street one Number:0 Map:map["ola" "joao":"adeus" "joao" c:d]}
:: Example 2: &{Name:joao Age:30 Address:0xc000086b80 Numbers:[1 2 3] Others:map["ola" "joao":0xc000086c80 c:0xc000086cc0 addresses:0xc000086d00] Addresses:[0xc000086d40 0xc000086d80]}

 [1,2,3]
:: LEN: 3
:: Example: [1 2 3]

 []
:: LEN: 0
:: Example: []

 {"one":1,"two":2,"three":3}
:: LEN: 3
:: Example: map[three:3 one:1 two:2]
:: LEN: 1
:: Example: &[123 34 116 101 115 116 34 58 32 34 111 110 101 34 44 32 34 116 101 115 116 34 58 32 34 116 119 111 34 125]
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
