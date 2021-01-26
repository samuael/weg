package service

// MainService struct representing the main service class and 
// having continuously running "run()"  method as a handler of messages 
type MainService struct {
	ClientMap map[string]string 
}

// Run method handling sending messages to clients and Coordinating the accessing of 
// Client Map and Group Map ... 
func (mainservice *MainService) Run() {

}
