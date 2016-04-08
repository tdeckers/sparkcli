package main

import (
	"encoding/json"
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
							log.Fatal("Usage: sparkcli rooms create <name>")
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
						if c.NArg() > 1 {
							log.Fatal("Usage: sparkcli rooms get <id>")
						}
						id := c.Args().Get(0)
						if id == "" { // try default room
							id = config.DefaultRoomId
							if id == "" {
								log.Fatal("Usage: sparkcli rooms get <id> (no default room configured)")
							}
						}
						roomService := api.RoomService{Client: client}
						room, err := roomService.Get(id)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							jsonMsg, err := json.MarshalIndent(room, "", "  ")
							if err != nil {
								log.Fatal("Failed to convert room.")
							}
							fmt.Print(string(jsonMsg))
						}
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete a room",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli rooms delete <id>")
						}
						id := c.Args().Get(0)
						roomService := api.RoomService{Client: client}
						err := roomService.Delete(id)
						//TODO: if error is '400 Bad Request', try deleting by name?
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("Room deleted.")
						}
					},
				},
				// Convenience actions (not available in Cisco Spark API)
				{
					Name:  "default",
					Usage: "save default room in config",
					Action: func(c *cli.Context) {
						if c.NArg() > 1 {
							log.Fatal("Usage: sparkcli rooms default (<id>)")
						}
						if c.NArg() == 1 {
							id := c.Args().Get(0)
							config.DefaultRoomId = id
							config.Save()
						} else {
							fmt.Print(config.DefaultRoomId)
						}
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
						// If no arg provided, also use default room.
						if c.NArg() > 1 {
							log.Fatal("Usage: sparkcli messages list <roomId>")
						}
						id := c.Args().Get(0)
						if id == "" || id == "-" {
							id = config.DefaultRoomId
							if id == "" {
								log.Println("No default room configured.")
								log.Fatal("Usage: sparkcli messages list <roomId>")
							}
						}
						msgService := api.MessageService{Client: client}
						msgs, err := msgService.List(id)
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
				{
					Name:    "get",
					Aliases: []string{"g"},
					Usage:   "get message details",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli messages get <id>")
						}
						id := c.Args().Get(0)
						msgService := api.MessageService{Client: client}
						msg, err := msgService.Get(id)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							jsonMsg, err := json.MarshalIndent(msg, "", "  ")
							if err != nil {
								log.Fatal("Failed to convert message.")
							}
							fmt.Print(string(jsonMsg))
						}
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete a message",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli messages delete <id>")
						}
						id := c.Args().Get(0)
						msgService := api.MessageService{Client: client}
						err := msgService.Delete(id)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Print("Message deleted.")
						}
					},
				},
			},
		},
		{
			Name:    "people",
			Aliases: []string{"p"},
			Usage:   "operations on people",
			Subcommands: []cli.Command{
				{
					Name:    "get",
					Aliases: []string{"g"},
					Usage:   "get your details",
					Action: func(c *cli.Context) {
						id := "me"
						if c.NArg() == 1 { // if argument, use that as id
							id = c.Args().Get(0)
						}
						peopleService := api.PeopleService{Client: client}
						person, err := peopleService.Get(id)
						if err != nil {
							fmt.Println(err)
						} else {
							jsonPerson, err := json.MarshalIndent(person, "", "  ")
							if err != nil {
								log.Fatal("Faild to convert person")
							}
							fmt.Print(string(jsonPerson))

						}

					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list people",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "email",
						},
						cli.StringFlag{
							Name: "name",
						},
					},
					Action: func(c *cli.Context) {
						email := c.String("email")
						name := c.String("name")
						peopleService := api.PeopleService{Client: client}
						people, err := peopleService.List(email, name)
						if err != nil {
							fmt.Println(err)
						} else {
							for _, person := range *people {
								fmt.Printf("%v (%v)\n", person.DisplayName, person.Emails[0])
							}

						}
					},
				},
			},
		},
		{
			Name:    "memberships",
			Aliases: []string{"ms"},
			Usage:   "operations on memberships",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list memberships",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "room",
						},
						cli.StringFlag{
							Name: "personid",
						},
						cli.StringFlag{
							Name: "email",
						},
					},
					Action: func(c *cli.Context) {
						roomId := c.String("room")
						personId := c.String("personid")
						personEmail := c.String("email")
						memberService := api.MemberService{Client: client}
						mss, err := memberService.List(roomId, personId, personEmail)
						if err != nil {
							fmt.Println(err)
						} else {
							for _, ms := range *mss {
								fmt.Printf("%v: %v\n", ms.PersonEmail, ms.RoomId)
							}

						}
					},
				},
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "create memberships",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "room",
						},
						cli.StringFlag{
							Name: "personid",
						},
						cli.StringFlag{
							Name: "email",
						},
					},
					Action: func(c *cli.Context) {
						roomId := c.String("room")
						personId := c.String("personid")
						personEmail := c.String("email")
						memberService := api.MemberService{Client: client}
						ms, err := memberService.Create(roomId, personId, personEmail)
						if err != nil {
							fmt.Println(err)
						} else {
							jsonMs, err := json.MarshalIndent(ms, "", "  ")
							if err != nil {
								log.Fatal("Faild to convert membership")
							}
							fmt.Print(string(jsonMs))
						}

					},
				},
				{
					Name:    "get",
					Aliases: []string{"g"},
					Usage:   "get membership details",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli memberships get <id>")
						}
						id := c.Args().Get(0)
						msService := api.MemberService{Client: client}
						ms, err := msService.Get(id)
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						} else {
							jsonMs, err := json.MarshalIndent(ms, "", "  ")
							if err != nil {
								log.Fatal("Failed to convert membership.")
							}
							fmt.Print(string(jsonMs))
						}

					},
				},
				{
					Name:    "update",
					Aliases: []string{"u"},
					Usage:   "update membership",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name: "mod",
						},
					},
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli memberships update -mod <id>")
						}
						id := c.Args().Get(0)
						// TODO: avoid doing update if flag is not present.
						moderator := c.Bool("mod")
						msService := api.MemberService{Client: client}
						ms, err := msService.Update(id, moderator)
						if err != nil {
							fmt.Print(err)
							os.Exit(-1)
						} else {
							jsonMs, err := json.MarshalIndent(ms, "", "  ")
							if err != nil {
								log.Fatal("Failed to convert membership.")
							}
							fmt.Print(string(jsonMs))
						}
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete membership",
					Action: func(c *cli.Context) {
						if c.NArg() != 1 {
							log.Fatal("Usage: sparkcli memberships delete <id>")
						}
						id := c.Args().Get(0)
						msService := api.MemberService{Client: client}
						err := msService.Delete(id)
						if err != nil {
							fmt.Print(err)
							os.Exit(-1)
						} else {
							fmt.Println("Membership deleted.")
						}

					},
				},
			},
		},
	}
	app.Run(os.Args)
}
