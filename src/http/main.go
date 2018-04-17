package main

var channel chan int = make(chan int)

func main() {
	hosts := getHost()
	saveHost(hosts)
	<- channel
}
