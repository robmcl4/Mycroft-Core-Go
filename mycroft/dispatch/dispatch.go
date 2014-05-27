package dispatch


import (
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/cmd"
)


var commandChan chan *cmd.Command = make(chan *cmd.Command, 100)
var preemptChan chan *cmd.Command = make(chan *cmd.Command, 100)


// add a command to the dispatcher
func Enqueue(c *cmd.Command) {
    commandChan <- c
}


// skip to the front of the queue
func PreemptQueue(c *cmd.Command) {
    preemptChan <- c
}


// Listen for and execute commands
func Dispatch() {
    for {
        select {
        case toExec := <- preemptChan:
            toExec.Execute()
        default:
            toExec := <- commandChan
            toExec.Execute()
        }
    }
}
