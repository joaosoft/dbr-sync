package main

import "profile"

func main() {
	p, err := profile.NewProfile()
	if err != nil {
		panic(err)
	}

	if err := p.Start(); err != nil {
		panic(err)
	}
}
