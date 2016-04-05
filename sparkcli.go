package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tdeckers/sparkcli/api"
	"github.com/tdeckers/sparkcli/util"
	"log"
	"os"
)

func main() {
	config := util.Configuration{}
	config.Load()
	client := util.NewClient(&config)
	app := cli.NewApp()
	app.Name = "sparkcli"
	app.Usage = "Command Line Interface for Cisco Spark"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "login to Cisco Spark",
			Action: func(c *cli.Context) {
				log.Println("Logging in")
				login := util.NewLogin(&config, client)
				login.Authorize()
			},
		},
		{
			Name:    "rooms",
			Aliases: []string{"r"},
			Usage:   "operations on rooms",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all rooms",
					Action: func(c *cli.Context) {
						roomService := api.RoomService{Client: client}
						rooms, err := roomService.List()
						if err != nil {
							fmt.Println(err)
						} else {
							for _, room := range *rooms {
								fmt.Printf("%s: %s\n", room.Id, room.Title)
							}
						}
					},
				},
				{
					Name:  "create",
					Usage: "create a new room",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Must provide name for the new room")
						}
						name := c.Args().Get(0)
						roomService := api.RoomService{Client: client}
						room, err := roomService.Create(name)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							// Print just roomId, so can assign to env variable if desired.
							fmt.Print(room.Id)
						}
					},
				},
				{
					Name:  "get",
					Usage: "get room details",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Must provide room id")
						}
						id := c.Args().Get(0)
						roomService := api.RoomService{Client: client}
						room, err := roomService.Get(id)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							fmt.Printf("%s - ...", room.Title)
						}
					},
				},
				{
					Name:  "delete",
					Usage: "delete a room",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Must provide id for room to delete")
						}
						id := c.Args().Get(0)
						roomService := api.RoomService{Client: client}
						err := roomService.Delete(id)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("Room deleted.")
						}
					},
				},
				// TODO: secondary actions: exists? limit list, ...
			},
		},
		{
			Name:    "messages",
			Aliases: []string{"m"},
			Usage:   "operations on messages",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all messages",
					Action: func(c *cli.Context) {
						// msgService := MessageService{client: client}
						log.Fatal("Not implemented")
					},
				},
				{
					Name:  "create",
					Usage: "create a new message",
					Action: func(c *cli.Context) {
						if c.NArg() != 2 {
							log.Fatal("Usage: ... messages create <room> <msg>")
						}
						room := c.Args().Get(0)
						msgTxt := c.Args().Get(1)
						msgService := api.MessageService{Client: client}
						msg, err := msgService.Create(room, msgTxt)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							fmt.Print(msg.Id)
						}

					},
				},
			},
		},
	}
	//	app.Action = func(c *cli.Context) {
	//		log.Println("Greetings")
	//	}
	app.Run(os.Args)
}
