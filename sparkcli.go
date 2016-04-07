package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tdeckers/sparkcli/api"
	"github.com/tdeckers/sparkcli/util"
	"log" // TODO: change to https://github.com/Sirupsen/logrus
	"os"
	"strings"
)

func main() {
	config := util.GetConfiguration()
	config.Load()
	client := util.NewClient(config)
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
				login := util.NewLogin(config, client)
				login.Authorize()
			},
		},
		{
			Name:    "rooms",
			Aliases: []string{"r"},
			Usage:   "operations on rooms",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all rooms",
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
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "create a new room",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli room create <name>")
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
					Name:    "get",
					Aliases: []string{"g"},
					Usage:   "get room details",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli room get <id>")
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
					Name:    "update",
					Aliases: []string{"u"},
					Usage:   "update room details",
					Action: func(c *cli.Context) {
						if c.NArg() < 2 {
							log.Fatal("Usage: sparkcli room update <id> <name>")
						}
						id := c.Args().Get(0)
						name := strings.Join(c.Args().Tail(), " ")
						roomService := api.RoomService{Client: client}
						room, err := roomService.Update(id, name)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							fmt.Printf("%v", room.Id)
						}
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete a room",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli room delete <id>")
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
				// Secondary actions (not part of native Spark API)
				{
					Name:  "default",
					Usage: "save default room to config",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli room default <id>")
						}
						id := c.Args().Get(0)
						config.DefaultRoomId = id
						config.Save()
						fmt.Printf("Default room set to %v", id)
					},
				},
			},
		},
		{
			Name:    "messages",
			Aliases: []string{"m"},
			Usage:   "operations on messages",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all messages",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli messages list <roomId>")
						}
						roomId := c.Args().Get(0)
						msgService := api.MessageService{Client: client}
						msgs, err := msgService.List(roomId)
						if err != nil {
							fmt.Println(err)
						} else {
							for _, msg := range *msgs {
								fmt.Printf("[%v] %v: %v\n", msg.Created, msg.PersonEmail, msg.Text)
							}
						}
					},
				},
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "create a new message",
					Action: func(c *cli.Context) {
						// TODO: change this to take all args after the second as additional text.
						if c.NArg() < 1 {
							log.Fatal("Usage: sparkcli messages create <room> <msg>")
						}
						room := c.Args().Get(0)
						msgTxt := strings.Join(c.Args().Tail(), " ")
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
