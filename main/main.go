package main

import "session"

func main() {
	m, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}
}
