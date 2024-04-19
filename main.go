package main

import (
	"KubernetesGo/Controllers"
	"fmt"
	"time"
)

func main()  {
	startTime := time.Now()
	fmt.Println("Start Time: ",startTime)
	Controllers.Router()
	fmt.Println("End Time: ", time.Since(startTime))
}
