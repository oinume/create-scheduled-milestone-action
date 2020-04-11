package main

import (
	"errors"
	"os"
	"time"
)

func main() {
	//myInput := os.Getenv("INPUT_MYINPUT")
	//
	//output := fmt.Sprintf("Hello %s", myInput)
	//
	//fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, output))
	a := &app{
		// TODO: github client
	}
	title := os.Getenv("INPUT_TITLE")
	description := os.Getenv("INPUT_DESCRIPTION")
	dueOn := os.Getenv("INPUT_DUE_ON")
	if err := a.run(title, description, dueOn); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

type app struct {}

func (c *app) run(title, description string, dueOn string) error {
	if title == "" {
		return errors.New("'title' is required")
	}

	_, err := c.validateDueOn(dueOn)
	if err != nil {
		return err
	}

	return nil
}

func (c *app) validateDueOn(dueOn string) (time.Time, error) {
	return time.Time{}, nil
}
