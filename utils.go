package main

import (
	"os/user"
	"strings"

	"github.com/codegangsta/cli"
)

func taskCoordinatorFromContext(c *cli.Context) (*TaskCoordinator, error) {
	location := c.GlobalString("location")
	if location[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		location = strings.Replace(location, "~/", usr.HomeDir+"/", 1)
	}

	storage, err := NewJSONStore(location)
	if err != nil {
		return nil, err
	}

	return &TaskCoordinator{storage}, nil
}
