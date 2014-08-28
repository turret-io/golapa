package main

import (
        "github.com/garyburd/redigo/redis"
        "encoding/json"
        "fmt"
		"github.com/turret-io/golapa/golapa"
		"github.com/turretIO/turret-io-go"
		"os"
		"log"
)

func main() {
		via := os.Getenv("SEND_VIA")
		api_key := os.Getenv("TURRET_API_KEY")
		api_secret := os.Getenv("TURRET_API_SECRET")

        conn, err := redis.Dial("tcp", ":6379")
        if err != nil {
                panic(err)
        }
        defer conn.Close()

        for {
                var queue string
                var data  []byte
                reply, err := redis.Values(conn.Do("BLPOP", "signup_worker", 0))
                if err != nil {
                        fmt.Println(err.Error())
                }
                if _, err := redis.Scan(reply, &queue, &data); err != nil {
                        panic(err)
                }

				if via == "email" {
					var msg golapa.Message
					err = json.Unmarshal(data, &msg)
					if err != nil {
							panic(err)
					}
					mailer := &golapa.StandardEmailer{""}
					mailer.Send(&msg)
				}

				if via == "turret_io" {
						var payload map[string]string
						err = json.Unmarshal(data, &payload)
						if err != nil {
								log.Print(err.Error())
						}

						attr_map := make(map[string]string)
						prop_map := make(map[string]string)
						attr_map["contact_name"] = payload["name"]
						attr_map["signedup"] = "1"

						turret := turretIO.NewTurretIO(api_key, api_secret)
						inst := turretIO.NewUser(turret)
						resp, err := inst.Set(payload["email"], attr_map, prop_map)
						log.Print(resp)
						if err != nil {
								log.Print(err.Error())
						}
				}

        }

}
