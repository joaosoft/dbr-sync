Validator
================

[![Build Status](https://travis-ci.org/joaosoft/validator.svg?branch=master)](https://travis-ci.org/joaosoft/validator) | [![codecov](https://codecov.io/gh/joaosoft/validator/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/validator) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/validator)](https://goreportcard.com/report/github.com/joaosoft/validator) | [![GoDoc](https://godoc.org/github.com/joaosoft/validator?status.svg)](https://godoc.org/github.com/joaosoft/validator)

A simple struct validator by tags (exported fields only).

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for validations
###### << command >>={id_field} can be used on all commands and will be replaced with the value of the field with id=id_field; if not exists by the manual sent args, if not exists json:"id_field"
* value (equal to)
* not (not equal to)
* options (one of the options)
* not-options (none of the options)
* size (size equal to)
* min 
* max 
* not-empty (also supports uuid empty validation)
* is-empty (also supports uuid empty validation)
* not-null 
* is-null 
* regex
* url
* email
* uuid
* base64
* ip, ipv4, ipv6
* callback (add handler validations)
* error (simple and multi error handling `validate:"value=1, error={errorValue1}, max=10, error={errorMax10}"`)
* if (conditional validation between fields with operators ("and", "or") [define id=xpto])
* alpha (the value needs to be alphanumeric)
* numeric (the value needs to be numeric)
* bool (the value needs to be boolean [true or false])
* item:<< command >>> (allows you to validate array or map items individually, [example: "item:size=10", means that the array items need to have the size of 10])
* key:<< command >>> (allows you to validate a map key's individually, [example: "key:size=10", means that the map key's need to have the size of 10])
* prefix
* suffix
* contains
* hex
* file

* args (arguments that will be available on callbacks ValidationData struct)

## With support for changing values
###### << command >>={id_field} can be used on all commands and will be replaced with the value of the field with id=id_field; if not exists by the manual sent args, if not exists json:"id_field"
###### to use this you need to use the variable address, like this `validator.Validate(&example)`
* set (allows to set native values) 
* set-empty
* set-md5
* set-random
* set-sanitize (clean characters)
* set-key (converts the value to a url valid key. You can also do key=xpto or key={id} where the id is other field id [example "This is a test" to "this-is-a-test"])
* set-trim
* set-title
* set-upper
* set-lower
* set-distinct (remove duplicated values from slices of primitive types)

## With methods for
* AddBefore (add a before-validation)
* AddMiddle (add a middle-validation [by default has all validations])
* AddAfter (add a after-validation [by default has error validation])
* SetErrorCodeHandler (function to get the error when defined with error={xpto:arg1;arg2})
* SetValidateAll (when activated, validates all object instead of stopping on the first error)
* SetTag (set validation tag to other that you define)
* SetSanitize (set sanitize strings)
* AddCallback (set a specific callback validation)
* Validate (object to validate, arguments...)

## Dependecy Management
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/validator
```

## Usage 
This examples are available in the project at [validator/examples](https://github.com/joaosoft/validator/tree/master/examples)

### Code
```go
const (
	regexForMissingParms = `%\+?[a-z]`
	constRegexForEmail     = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

)

type Data string

type NextSet struct {
	Set int `validate:"set=321, id=next_set"`
}

type Items struct {
	A string
	B int
}

type Example struct {
	Interface               interface{}       `validate:"not-null, not-empty"`
	Interfaces              []interface{}     `validate:"not-null, not-empty"`
	Array                   []string          `validate:"item:size=5"`
	Array2                  []string          `validate:"item:set-distinct"`
	Array3                  Items             `validate:"item:size=5"`
	Map4                    map[string]string `validate:"item:size=5, key:size=5"`
	Name                    string            `validate:"value=joao, dummy_middle, error={{ErrorTag1:a;b}}, max=10"`
	Age                     int               `validate:"value=30, error={{ErrorTag99}}"`
	Street                  int               `validate:"max=10, error={{ErrorTag3}}"`
	Brothers                []Example2        `validate:"size=1, error={{ErrorTag4}}"`
	Id                      uuid.UUID         `validate:"not-empty, error={{ErrorTag5}}"`
	Option1                 string            `validate:"options=aa;bb;cc, error={{ErrorTag6}}"`
	Option2                 int               `validate:"options=11;22;33, error={{ErrorTag7}}"`
	Option3                 []string          `validate:"options=aa;bb;cc, error={{ErrorTag8}}"`
	Option4                 []int             `validate:"options=11;22;33, error={{ErrorTag9}}"`
	NotOption               []int             `validate:"not-options=11;22;33"`
	Map1                    map[string]int    `validate:"options=aa:11;bb:22;cc:33, error={{ErrorTag10}}"`
	Map2                    map[int]string    `validate:"options=11:aa;22:bb;33:cc, error={{ErrorTag11}}"`
	Url                     string            `validate:"url"`
	Email                   string            `validate:"email"`
	unexported              string
	IsNill                  *string `validate:"not-empty, error={{ErrorTag17}}"`
	Sanitize                string  `validate:"set-sanitize=a;b;teste, error={{ErrorTag17}}"`
	Callback                string  `validate:"callback=dummy_callback;dummy_callback_2, error={{ErrorTag19}}"`
	Password                string  `json:"password" validate:"id=password"`
	PasswordConfirm         string  `validate:"value={password}"`
	MyName                  string  `validate:"id=name"`
	MyAge                   int     `validate:"id=age"`
	MyValidate              int     `validate:"if=(id=age value=30) or (id=age value=31) and (id=name value=joao), value=10"`
	DoubleValidation        int     `validate:"not-empty, error=20, min=5, error={{ErrorTag21}}"`
	Set                     int     `validate:"set=321, id=set"`
	NextSet                 NextSet
	DistinctIntPointer      []*int    `validate:"set-distinct"`
	DistinctInt             []int     `validate:"set-distinct"`
	DistinctString          []string  `validate:"set-distinct"`
	DistinctBool            []bool    `validate:"set-distinct"`
	DistinctFloat           []float32 `validate:"set-distinct"`
	IsZero                  int       `validate:"is-empty"`
	Trim                    string    `validate:"set-trim"`
	Lower                   string    `validate:"set-lower"`
	Upper                   string    `validate:"set-upper"`
	Key                     string    `validate:"set-key"`
	KeyValue                string    `validate:"id=my_value"`
	KeyFromValue            string    `validate:"set-key={my_value}"`
	NotMatch1               string    `validate:"id=not_match"`
	NotMatch2               string    `validate:"not={not_match}"`
	TypeAlpha               string    `validate:"alpha"`
	TypeNumeric             string    `validate:"numeric"`
	TypeBool                string    `validate:"bool"`
	ShouldBeNull            *string   `validate:"is-null"`
	ShouldNotBeNull         *string   `validate:"not-null"`
	FirstMd5                string    `validate:"set-md5"`
	SecondMd5               string    `validate:"set-md5=ola"`
	EnableEncodeRandom      bool      `validate:"id=random_enable"`
	EnableEncodeRandomTitle bool      `validate:"id=random_title_enable"`
	Random                  string    `cleanup:"if=(id=random_enable value=true), set-random, if=(id=random_title_enable value=true), set-title"`
	RandomArg               string    `cleanup:"if=(arg=random_enable value=true), set-random, if=(arg=random_title_enable value=true), set-title"`
	RandomClean             string    `cleanup:"if=(id=random_enable value=true), set-random, if=(id=random_title_enable value=true), set="`
	StringPrefix            string    `validate:"prefix=ola"`
	StringSuffix            string    `validate:"suffix=mundo"`
	StringContains          string    `validate:"contains=a m"`
	Hex                     string    `validate:"hex"`
	File                    string    `validate:"file"`
	EmptyText               string    `validate:"set-empty"`
	EmptyInt                int       `validate:"set-empty"`
	EmptyArrayString        []string  `validate:"set-empty"`
	EmptyArrayInt           []int     `validate:"set-empty"`
}

type Example2 struct {
	Name            string         `validate:"value=joao, dummy_middle, error={{ErrorTag1:a;b}}, max=10"`
	Age             int            `validate:"value=30, error={{ErrorTag99}}"`
	Street          int            `validate:"max=10, error={{ErrorTag3}}"`
	Id              uuid.UUID      `validate:"not-empty, error={{ErrorTag5}}"`
	Option1         string         `validate:"options=aa;bb;cc, error={{ErrorTag6}}"`
	Option2         int            `validate:"options=11;22;33, error={{ErrorTag7}}"`
	Option3         []string       `validate:"options=aa;bb;cc, error={{ErrorTag8}}"`
	Option4         []int          `validate:"options=11;22;33, error={{ErrorTag9}}"`
	NotOption       []int          `validate:"not-options=11;22;33"`
	Map1            map[string]int `validate:"options=aa:11;bb:22;cc:33, error={{ErrorTag10}}"`
	Map2            map[int]string `validate:"options=11:aa;22:bb;33:cc, error={{ErrorTag11}}"`
	Url             string         `validate:"url"`
	Email           string         `validate:"email"`
	unexported      string
	IsNill          *string   `validate:"not-empty, error={{ErrorTag17}}"`
	Sanitize        string    `validate:"set-sanitize=a;b;teste, error={{ErrorTag17}}"`
	Callback        string    `validate:"callback=dummy_callback, error={{ErrorTag19}}"`
	CallbackArgs    string    `validate:"args=a;b;c, callback=dummy_args_callback"`
	Password        string    `json:"password" validate:"id=password"`
	PasswordConfirm string    `validate:"value={password}"`
	UUID            string    `validate:"uuid"`
	UUIDStruct      uuid.UUID `validate:"uuid"`
}

type Example3 struct {
	Name     string `validate:"value=joao"`
	LastName string `validate:"set=ribeiro"`
}

var dummy_middle_handler = func(context *validator.ValidatorContext, validationData *validator.ValidationData) []error {
	var rtnErrs []error

	err := errors.New("dummy middle responding...")
	rtnErrs = append(rtnErrs, err)

	return rtnErrs
}

func init() {
	validator.
		AddMiddle("dummy_middle", dummy_middle_handler).
		SetValidateAll(true).
		SetErrorCodeHandler(dummy_error_handler).
		AddCallback("dummy_callback", dummy_callback).
		AddCallback("dummy_args_callback", dummy_args_callback).
		AddCallback("dummy_callback_2", dummy_callback)
}

var errs = map[string]error{
	"ErrorTag1":  errors.New("error 1: a:%s, b:%s"),
	"ErrorTag2":  errors.New("error 2"),
	"ErrorTag3":  errors.New("error 3"),
	"ErrorTag4":  errors.New("error 4"),
	"ErrorTag5":  errors.New("error 5"),
	"ErrorTag6":  errors.New("error 6"),
	"ErrorTag7":  errors.New("error 7"),
	"ErrorTag8":  errors.New("error 8"),
	"ErrorTag9":  errors.New("error 9"),
	"ErrorTag10": errors.New("error 10"),
	"ErrorTag11": errors.New("error 11"),
	"ErrorTag12": errors.New("error 12"),
	"ErrorTag13": errors.New("error 13"),
	"ErrorTag14": errors.New("error 14"),
	"ErrorTag15": errors.New("error 15"),
	"ErrorTag16": errors.New("error 16"),
	"ErrorTag17": errors.New("error 17"),
	"ErrorTag18": errors.New("error 18"),
	"ErrorTag19": errors.New("error 19"),
	"ErrorTag20": errors.New("error 20"),
	"ErrorTag21": errors.New("error 21"),
}
var dummy_error_handler = func(context *validator.ValidatorContext, validationData *validator.ValidationData) error {
	if err, ok := errs[validationData.ErrorData.Code]; ok {
		var regx = regexp.MustCompile(regexForMissingParms)
		matches := regx.FindAllStringIndex(err.Error(), -1)

		if len(matches) > 0 {

			if len(validationData.ErrorData.Arguments) < len(matches) {
				validationData.ErrorData.Arguments = append(validationData.ErrorData.Arguments, validationData.Name)
			}

			err = fmt.Errorf(err.Error(), validationData.ErrorData.Arguments...)
		}

		return err
	}
	return nil
}

var dummy_callback = func(context *validator.ValidatorContext, validationData *validator.ValidationData) []error {
	return []error{errors.New("there's a bug here!")}
}

var dummy_args_callback = func(context *validator.ValidatorContext, validationData *validator.ValidationData) []error {
	fmt.Printf("\nthere are the following arguments %+v!", validationData.Arguments)
	return nil
}

func main() {
	intVal1 := 1
	intVal2 := 2
	id, _ := uuid.NewV4()
	str := "should be null"
	byts := [16]byte{}
	copy(byts[:], "1234567890123456")

	example := Example{
		Interface: &Example3{
			Name:     "JESSICA",
			LastName: "EMPTY",
		},
		Interfaces: []interface{}{
			&Example3{
				Name:     "JESSICA",
				LastName: "EMPTY",
			},
		},
		Array:  []string{"12345", "123456", "12345", "1234567"},
		Array2: []string{"111", "111", "222", "222"},
		Array3: Items{
			A: "123456",
			B: 1234567,
		},
		Map4:             map[string]string{"123456": "1234567", "12345": "12345"},
		Id:               id,
		Name:             "joao",
		Age:              30,
		Street:           10,
		Option1:          "aa",
		Option2:          11,
		Option3:          []string{"aa", "bb", "cc"},
		Option4:          []int{11, 22, 33},
		NotOption:        []int{1, 2, 3},
		Map1:             map[string]int{"aa": 11, "bb": 22, "cc": 33},
		Map2:             map[int]string{11: "aa", 22: "bb", 33: "cc"},
		Url:              "google.com",
		Email:            "joaosoft@gmail.com",
		Password:         "password",
		PasswordConfirm:  "password_errada",
		MyName:           "joao",
		MyAge:            30,
		MyValidate:       30,
		DoubleValidation: 0,
		Set:              123,
		NextSet: NextSet{
			Set: 123,
		},
		DistinctIntPointer:      []*int{&intVal1, &intVal1, &intVal2, &intVal2},
		DistinctInt:             []int{1, 1, 2, 2},
		DistinctString:          []string{"a", "a", "b", "b"},
		DistinctBool:            []bool{true, true, false, false},
		DistinctFloat:           []float32{1.1, 1.1, 1.2, 1.2},
		Trim:                    "     aqui       TEM     espaços    !!   ",
		Upper:                   "     aqui       TEM     espaços    !!   ",
		Lower:                   "     AQUI       TEM     ESPACOS    !!   ",
		Key:                     "     AQUI       TEM     ESPACOS    !!   ",
		KeyValue:                "     aaaaa     3245 79 / ( ) ? =  tem     espaços ...   !!  <<<< ",
		NotMatch1:               "A",
		NotMatch2:               "A",
		TypeAlpha:               "123",
		TypeNumeric:             "ABC",
		TypeBool:                "ERRADO",
		ShouldBeNull:            &str,
		FirstMd5:                "first",
		SecondMd5:               "second",
		EnableEncodeRandom:      true,
		EnableEncodeRandomTitle: true,
		Random:                  "o meu novo teste random 123",
		RandomArg:               "o meu novo teste random 123",
		RandomClean:             "o meu novo teste random 123",
		StringPrefix:            "ola",
		StringSuffix:            "mundo",
		StringContains:          "a m",
		Hex:                     "48656c6c6f20476f7068657221",
		File:                    "./README.md",
		EmptyText:               "text",
		EmptyInt:                111,
		EmptyArrayString:        []string{"text", "text"},
		EmptyArrayInt:           []int{1, 2},
		Brothers: []Example2{
			Example2{
				Name:            "jessica",
				Age:             10,
				Street:          12,
				Option1:         "xx",
				Option2:         99,
				Option3:         []string{"aa", "zz", "cc"},
				Option4:         []int{11, 44, 33},
				NotOption:       []int{11, 22, 33},
				Map1:            map[string]int{"aa": 11, "kk": 22, "cc": 33},
				Map2:            map[int]string{11: "aa", 22: "bb", 99: "cc"},
				Sanitize:        "b teste",
				Url:             "http://www.teste.pt",
				Email:           "joaosoft@gmail.com",
				Password:        "password",
				PasswordConfirm: "password",
				UUID:            "invalid",
				UUIDStruct:      byts,
			},
		},
	}

	fmt.Printf("\nBEFORE SET: %d", example.Set)
	fmt.Printf("\nBEFORE NEXT SET: %d", example.NextSet.Set)
	fmt.Printf("\nBEFORE TRIM: %s", example.Trim)
	fmt.Printf("\nBEFORE KEY: %s", example.Key)
	fmt.Printf("\nBEFORE FROM KEY: %s", example.KeyFromValue)
	fmt.Printf("\nBEFORE UPPER: %s", example.Upper)
	fmt.Printf("\nBEFORE LOWER: %s", example.Lower)

	fmt.Printf("\nBEFORE DISTINCT INT POINTER: %+v", example.DistinctIntPointer)
	fmt.Printf("\nBEFORE DISTINCT INT: %+v", example.DistinctInt)
	fmt.Printf("\nBEFORE DISTINCT STRING: %+v", example.DistinctString)
	fmt.Printf("\nBEFORE DISTINCT BOOL: %+v", example.DistinctBool)
	fmt.Printf("\nBEFORE DISTINCT FLOAT: %+v", example.DistinctFloat)
	fmt.Printf("\nBEFORE DISTINCT ARRAY2: %+v", example.Array2)
	fmt.Printf("\nBEFORE EMPTY TEXT: %+v", example.EmptyText)
	fmt.Printf("\nBEFORE EMPTY INT: %+v", example.EmptyInt)
	fmt.Printf("\nBEFORE EMPTY ARRAY TEXT: %+v", example.EmptyArrayString)
	fmt.Printf("\nBEFORE EMPTY ARRAY INT: %+v", example.EmptyArrayInt)

	// validate
	if errs := validator.Validate(&example,
		validator.NewArgument("random_enable", false),
		validator.NewArgument("random_title_enable", true),
	); len(errs) > 0 {
		fmt.Printf("\n\nERRORS: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("\nERROR: %s", err)
		}
	}

	// cleanup
	if errs := validator.NewValidator().SetTag("cleanup").Validate(&example,
		validator.NewArgument("random_enable", false),
		validator.NewArgument("random_title_enable", true),
	); len(errs) > 0 {
		fmt.Printf("\n\nERRORS: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("\nERROR: %s", err)
		}
	}

	fmt.Printf("\n\nAFTER SET: %d", example.Set)
	fmt.Printf("\nAFTER NEXT SET: %d", example.NextSet.Set)
	fmt.Printf("\nAFTER TRIM: %s", example.Trim)
	fmt.Printf("\nAFTER KEY: %s", example.Key)
	fmt.Printf("\nAFTER FROM KEY: %s", example.KeyFromValue)
	fmt.Printf("\n\nAFTER LOWER: %s", example.Lower)
	fmt.Printf("\n\nAFTER UPPER: %s", example.Upper)

	fmt.Printf("\nAFTER DISTINCT INT POINTER: %+v", example.DistinctIntPointer)
	fmt.Printf("\nAFTER DISTINCT INT: %+v", example.DistinctInt)
	fmt.Printf("\nAFTER DISTINCT STRING: %+v", example.DistinctString)
	fmt.Printf("\nAFTER DISTINCT BOOL: %+v", example.DistinctBool)
	fmt.Printf("\nAFTER DISTINCT FLOAT: %+v", example.DistinctFloat)
	fmt.Printf("\nAFTER DISTINCT ARRAY2: %+v", example.Array2)
	fmt.Printf("\nAFTER EMPTY TEXT: %+v", example.EmptyText)
	fmt.Printf("\nAFTER EMPTY INT: %+v", example.EmptyInt)
	fmt.Printf("\nAFTER EMPTY ARRAY TEXT: %+v", example.EmptyArrayString)
	fmt.Printf("\nAFTER EMPTY ARRAY INT: %+v", example.EmptyArrayInt)

	fmt.Printf("\nFIRST MD5: %+v", example.FirstMd5)
	fmt.Printf("\nSECOND MD5: %+v", example.SecondMd5)
	fmt.Printf("\nRANDOM: %+v", example.Random)
	fmt.Printf("\nRANDOM BY ARG: %+v", example.RandomArg)
	fmt.Printf("\nRANDOM BY ARG CLEAN: %+v", example.RandomClean)
	fmt.Printf("\nLAST NAME: %+v", example.Interface.(*Example3).LastName)
	fmt.Printf("\nLAST NAME 2: %+v", example.Interfaces[0].(*Example3).LastName)

	// validate embed struct
	if errs := validator.Validate(struct {
		Data struct {
			Name  string `validate:"not-null, is-empty"`
			Array []int  `validate:"not-null, not-empty"`
		} `validate:"not-null, not-empty"`
	}{},
		validator.NewArgument("random_enable", false),
		validator.NewArgument("random_title_enable", true),
	); len(errs) > 0 {
		fmt.Printf("\n\nERRORS: %d\n", len(errs))
		for _, err := range errs {
			fmt.Printf("\nERROR: %s", err)
		}
	}

	// benchmark
	timingValidator()
	timingManualValidation()
}

type Example4 struct {
	Name    string `validate:"value=joao, error={{ErrorTag1:a;b}}, max=10"`
	Age     int    `validate:"value=30, error={{ErrorTag99}}"`
	Street  int    `validate:"max=10, error={{ErrorTag3}}"`
	Option1 string `validate:"options=aa;bb;cc, error={{ErrorTag6}}"`
	Email   string `validate:"email"`
}

func timingValidator() {
	fmt.Println("\n-> timing with validator")

	// test
	start := time.Now()

	example := Example4{
		Name:    "joao",
		Age:     30,
		Street:  10,
		Option1: "aa",
		Email:   "joaosoft@gmail.com",
	}

	// validate
	if errs := validator.Validate(&example); len(errs) > 0 {
		fmt.Printf("\n\nERRORS: %+v\n", errs)
	}

	fmt.Printf("Elapsed time: %f", time.Since(start).Seconds())
}

func timingManualValidation() {
	fmt.Println("\n-> timing without validator")

	// test
	start := time.Now()

	example := Example4{
		Name:    "joao",
		Age:     30,
		Street:  10,
		Option1: "aa",
		Email:   "joaosoft@gmail.com",
	}

	// validate
	var errorList []error

	if example.Name != "joao" {
		errorList = append(errorList, fmt.Errorf(errs["ErrorTag1"].Error(), "a", "b"))
	}

	if example.Age != 30 {
		errorList = append(errorList, errs["ErrorTag99"])
	}

	if example.Street > 10 {
		errorList = append(errorList, errs["ErrorTag3"])
	}

	if example.Option1 != "aa" &&
		example.Option1 != "bb" &&
		example.Option1 != "cc" {
		errorList = append(errorList, errs["ErrorTag6"])
	}

	r, err := regexp.Compile(constRegexForEmail)
	if err != nil {
		panic(err)
	}

	if !r.MatchString(example.Email) {
		errorList = append(errorList, errors.New("invalid email"))
	}

	if len(errorList) > 0 {
		fmt.Printf("\n\nERRORS: %+v\n", errorList)
	}

	fmt.Printf("Elapsed time: %f", time.Since(start).Seconds())
}
```

> ##### Response:
```go
BEFORE SET: 123
BEFORE NEXT SET: 123
BEFORE TRIM:      aqui       TEM     espaços    !!   
BEFORE KEY:      AQUI       TEM     ESPACOS    !!   
BEFORE FROM KEY: 
BEFORE UPPER:      aqui       TEM     espaços    !!   
BEFORE LOWER:      AQUI       TEM     ESPACOS    !!   
BEFORE DISTINCT INT POINTER: [0xc0000182e8 0xc0000182e8 0xc000018300 0xc000018300]
BEFORE DISTINCT INT: [1 1 2 2]
BEFORE DISTINCT STRING: [a a b b]
BEFORE DISTINCT BOOL: [true true false false]
BEFORE DISTINCT FLOAT: [1.1 1.1 1.2 1.2]
BEFORE DISTINCT ARRAY2: [111 111 222 222]
BEFORE EMPTY TEXT: text
BEFORE EMPTY INT: 111
BEFORE EMPTY ARRAY TEXT: [text text]
BEFORE EMPTY ARRAY INT: [1 2]
there are the following arguments [a b c]!

ERRORS: 39

ERROR: the value [JESSICA] is different of the expected [joao] on field [Name]
ERROR: the value [JESSICA] is different of the expected [joao] on field [Name]
ERROR: the length [6] is lower then the expected [5] on field [Array] value [123456]
ERROR: the length [7] is lower then the expected [5] on field [Array] value [1234567]
ERROR: the length [6] is lower then the expected [5] on field [Array3] value [123456]
ERROR: the length [7] is lower then the expected [5] on field [Array3] value [1234567]
ERROR: the length [7] is lower then the expected [5] on field [Map4] value [1234567]
ERROR: the length [6] is lower then the expected [5] on field [Map4] value [123456]
ERROR: error 1: a:a, b:b
ERROR: error 1: a:a, b:b
ERROR: the value [10] is different of the expected [30] on field [Age]
ERROR: error 3
ERROR: error 5
ERROR: error 6
ERROR: error 7
ERROR: error 8
ERROR: error 9
ERROR: the value [11] shouldn't be equal to the excluded options [11;22;33] on field [NotOption]
ERROR: the value [22] shouldn't be equal to the excluded options [11;22;33] on field [NotOption]
ERROR: the value [33] shouldn't be equal to the excluded options [11;22;33] on field [NotOption]
ERROR: error 10
ERROR: error 11
ERROR: error 17
ERROR: error 17
ERROR: error 19
ERROR: the value [invalid] on field [UUID] should be a valid UUID
ERROR: the value [google.com] on field [Url] should be a valid Url
ERROR: error 17
ERROR: error 19
ERROR: the value [password_errada] is different of the expected [password] on field [PasswordConfirm]
ERROR: the value [30] is different of the expected [10] on field [MyValidate]
ERROR: 20
ERROR: error 21
ERROR: the expected [A] should be different of the [A] on field [NotMatch2]
ERROR: the value [123] is invalid for type alphanumeric on field [TypeAlpha] value [123]
ERROR: the value [ABC] is invalid for type numeric on field [TypeNumeric] value [ABC]
ERROR: the value [ERRADO] is invalid for type bool on field [TypeBool] value [ERRADO]
ERROR: the value should be null on field [ShouldBeNull] instead of [should be null]
ERROR: the value shouldn't be null on field [ShouldNotBeNull]

AFTER SET: 321
AFTER NEXT SET: 321
AFTER TRIM: aqui TEM espaços !!
AFTER KEY: 
AFTER FROM KEY: aaaaa-3245-79-tem-espacos-

AFTER LOWER:      aqui       tem     espacos    !!   

AFTER UPPER:      AQUI       TEM     ESPAÇOS    !!   
AFTER DISTINCT INT POINTER: [0xc0000182e8 0xc000018300]
AFTER DISTINCT INT: [1 2]
AFTER DISTINCT STRING: [a b]
AFTER DISTINCT BOOL: [true false]
AFTER DISTINCT FLOAT: [1.1 1.2]
AFTER DISTINCT ARRAY2: [111 222]
AFTER EMPTY TEXT: 
AFTER EMPTY INT: 0
AFTER EMPTY ARRAY TEXT: []
AFTER EMPTY ARRAY INT: []
FIRST MD5: d41d8cd98f00b204e9800998ecf8427e
SECOND MD5: 2fe04e524ba40505a82e03a2819429cc
RANDOM: Ç Ówl Òncq Okèxâ Ònsqçd 897
RANDOM BY ARG: I Wôa Émãô Xdâkã Vfùìbp 412
RANDOM BY ARG CLEAN: 
LAST NAME: ribeiro
LAST NAME 2: ribeiro

ERRORS: 2

ERROR: the value shouldn't be empty on field [Data]
ERROR: the value shouldn't be empty on field [Array]
-> timing with validator
Elapsed time: 0.000231
-> timing without validator
Elapsed time: 0.000110
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
