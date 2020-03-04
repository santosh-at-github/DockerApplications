package main

import (
        "database/sql"
        "fmt"
        "html/template"
        "io/ioutil"
        "log"
        "net"
        "net/http"
        "os"
        "time"

        "github.com/go-redis/redis"
        _ "github.com/go-redis/redis"
        _ "github.com/go-sql-driver/mysql"
)

var tpl = template.Must(template.ParseFiles("index.gohtml"))

func GetOutboundIP() net.IP {
        conn, err := net.Dial("udp", "8.8.8.8:80")
        if err != nil {
                log.Fatal(err)
        }
        defer conn.Close()

        localAddr := conn.LocalAddr().(*net.UDPAddr)

        return localAddr.IP
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
        dt := time.Now().String()
        //fmt.Println("Current date and time is: ", dt)
        hostname, err := os.Hostname()
        if err != nil {
                hostname = "error getting hostname"
        }
        //instance if from meta data
        instanceid := "nil"

        iresp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
        if err != nil {
                instanceid = "error getting instance IP"
        }
        defer iresp.Body.Close()

        instanceidresp, err := ioutil.ReadAll(iresp.Body)
        if err != nil {
                instanceid = "error getting instance IP"
        }
        instanceid = string(instanceidresp)

        //private IP form metadata
        instanceprivateip := "nil"

        resp, err := http.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
        if err != nil {
                instanceprivateip = "error getting instance IP"
        }
        defer resp.Body.Close()

        instanceipresp, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                instanceprivateip = "error getting instance IP"
        }
        instanceprivateip = string(instanceipresp)

        //public ip form meta data
        instancepublicip := "nil"
        nresp, err := http.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
        if err != nil {
                instancepublicip = "error getting instance IP"
        }
        defer nresp.Body.Close()
        instancepublicipresp, err := ioutil.ReadAll(nresp.Body)
        if err != nil {
                instancepublicip = "error getting instance IP"
        }
        instancepublicip = string(instancepublicipresp)

        //increment count in ecache
        IncrementTotalViews()

        //store instance count in mysql
        incrementinstancecount(instanceid)

        // construct data to pass to html template
        data := DisplayPageData{
                TheTime:           dt,
                ClientIP:          r.RemoteAddr,
                HostName:          hostname,
                ContainerIP:       GetOutboundIP().String(),
                InstancePrivateIP: instanceprivateip,
                InstancePublicIP:  instancepublicip,
                InstanceId:        instanceid,
                InstanceCount:     getinstancecount(instanceid),
                TotalViewCount:    GetTotalViews(),
        }
        //load or execute html template with passed data
        tpl.Execute(w, data)
}

type DisplayPageData struct {
        TheTime           string
        ClientIP          string
        HostName          string
        ContainerIP       string
        InstancePrivateIP string
        InstancePublicIP  string
        InstanceId        string
        InstanceCount     string
        TotalViewCount    string
}

func main() {
        getEnv("DB_NAME", "")
        getEnv("DB_USER", "")
        getEnv("DB_PASSWORD", "")
        getEnv("EC_ENDPOINT", "")
        //ExampleNewClient()

        //getinstancecount("i-something")

        port := os.Getenv("PORT")
        if port == "" {
                port = "3000"
        }

        mux := http.NewServeMux()

        mux.HandleFunc("/", indexHandler)
        http.ListenAndServe(":"+port, mux)
}

func getinstancecount(instanceid string) string {
        dbname := os.Getenv("DB_NAME")
        dbuser := os.Getenv("DB_USER")
        dbpass := os.Getenv("DB_PASSWORD")
        //dbhost := os.Getenv("DB_HOST")

        db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp"+"(mysql"+":3306)"+"/"+dbname)
        //db, err := sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname)
        if err != nil {
                panic(err.Error())
        }
        // defer the close till after the main function has finished
        // executing
        defer db.Close()
        // Execute the query
        results, err := db.Query("SELECT count FROM instance WHERE instance_id = '" + instanceid + "'")
        if err != nil {
                panic(err.Error()) // proper error handling instead of panic in your app
        }
        var icount string

        for results.Next() {
                // var tag Tag
                // // for each row, scan the result into our tag composite object
                // err = results.Scan(&tag.ID, &tag.Name)
                if err != nil {
                        icount = "error getting count"
                        fmt.Println("error: ", err.Error())
                        return icount
                        //panic(err.Error()) // proper error handling instead of panic in your app
                }
                results.Scan(&icount)
        }
        //fmt.Println("count is: ", icount)
        return icount
}

func incrementinstancecount(instanceid string) {
        dbname := os.Getenv("DB_NAME")
        dbuser := os.Getenv("DB_USER")
        dbpass := os.Getenv("DB_PASSWORD")
        //dbhost := os.Getenv("DB_HOST")
        //instanceid = "something3"
        db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp"+"(mysql"+":3306)"+"/"+dbname)
        if err != nil {
                fmt.Println("error: ", err.Error())
                //panic(err.Error())
        }
        // defer the close till after the main function has finished
        // executing
        defer db.Close()

        //check if instance id in table, else insert record
        var exists bool
        row := db.QueryRow("SELECT EXISTS(SELECT count FROM instance WHERE instance_id = '" + instanceid + "')")
        if err := row.Scan(&exists); err != nil {
                fmt.Println(" error")
        } else if !exists {
                if _, err := db.Exec("insert into instance values ('" + instanceid + "',0)"); err != nil {
                        fmt.Println("error: ", err.Error())
                }
        }

        // Execute the query
        results, err := db.Query("UPDATE instance SET count = count + 1 WHERE instance_id = '" + instanceid + "'")
        if err != nil {
                panic(err.Error()) // proper error handling instead of panic in your app
        }
        var icount string

        for results.Next() {
                // var tag Tag
                // // for each row, scan the result into our tag composite object
                // err = results.Scan(&tag.ID, &tag.Name)
                if err != nil {
                        //icount = "error getting count"
                        fmt.Println("error: ", err.Error())
                        //panic(err.Error()) // proper error handling instead of panic in your app
                }
                results.Scan(&icount)
        }
        //fmt.Println("count is: ", icount)
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
        if value, exists := os.LookupEnv(key); exists {
                fmt.Println("env defined key: ", key)
                fmt.Println("env defined value: ", value)
                return value
        } else {
                panic("ensure all env variables are set correctly: ")
        }

}

func ExampleNewClient() {
        ec_endpoint := os.Getenv("EC_ENDPOINT")
        client := redis.NewClient(&redis.Options{
                Addr:     ec_endpoint + ":6379",
                Password: "", // no password set
                DB:       0,  // use default DB
        })

        pong, err := client.Ping().Result()
        fmt.Println(pong, err)
        // Output: PONG <nil>
}

func IncrementTotalViews() {
        ec_endpoint := os.Getenv("EC_ENDPOINT")
        client := redis.NewClient(&redis.Options{
                Addr:     ec_endpoint + ":6379",
                Password: "", // no password set
                DB:       0,  // use default DB
        })

        client.Incr("totalviewcount").Result()
        //fmt.Println(pong, err)
        // Output: PONG <nil>
}

func GetTotalViews() string {
        ec_endpoint := os.Getenv("EC_ENDPOINT")
        client := redis.NewClient(&redis.Options{
                Addr:     ec_endpoint + ":6379",
                Password: "", // no password set
                DB:       0,  // use default DB
        })

        val, err := client.Get("totalviewcount").Result()
        if err != nil {
                val = "error getting count"
                fmt.Println(err.Error())
        }
        return val
}


