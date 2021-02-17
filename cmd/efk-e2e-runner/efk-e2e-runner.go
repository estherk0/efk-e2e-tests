package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	args := map[string]*string{}
	args["esURL"] = flag.String("es-host", "https://127.0.0.1:9200", "Elasticsearch URL")
	args["kibanaURL"] = flag.String("kibana-host", "https://127.0.0.1:5601", "Kibana URL")
	args["targetList"] = flag.String("targets", "elasticsearch,kibana", "test target list seperated by comma")

	flag.Parse()
	log.Println(*args["kibanaURL"])

	targetList := strings.Split(*args["targetList"], ",")

	for _, target := range targetList {
		switch target {
		case "elasticsearch":
			err := RunElasticsearchE2ETest(*args["esURL"])
			if err != nil {
				log.Fatal("[ERROR] ", err.Error())
				panic(err)
			}
		case "kibana":
			err := RunKibanaE2ETest()
			if err != nil {
				log.Fatal("[ERROR] ", err.Error())
				panic(err)
			}
		default:
			log.Println("[WARN] Not supported target", target)
		}
		log.Println("[INFO] Successfully completed test for", target)
	}
}
