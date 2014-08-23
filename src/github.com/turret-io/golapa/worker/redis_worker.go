package main

import (
        "github.com/garyburd/redigo/redis"
        "encoding/json"
        "fmt"
		"github.com/turret-io/golapa/golapa"
)

func main() {
        conn, err := redis.Dial("tcp", ":6379")
        if err != nil {
                panic(err)
        }
        defer conn.Close()

        for {
                var queue string
                var data  []byte
                reply, err := redis.Values(conn.Do("BLPOP", "email_worker", 0))
                if err != nil {
                        fmt.Println(err.Error())
                }
                if _, err := redis.Scan(reply, &queue, &data); err != nil {
                        panic(err)
                }
                var msg golapa.Message
                err = json.Unmarshal(data, &msg)
                if err != nil {
                        panic(err)
                }
                mailer := &golapa.StandardEmailer{""}
				mailer.Send(msg)
        }

}
