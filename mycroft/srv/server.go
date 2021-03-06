// Package srv implements the main Mycroft-core network listener.
package srv

import (
    "net"
    "fmt"
    "strings"
    "strconv"
    "errors"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/app"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/cmd"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/registry"
    "github.com/robmcl4/Mycroft-Core-Go/mycroft/logging"
)


// Starts listening for client connections.
// When a new application connects, launches listeners in a goroutine.
// Returns an error when error occurs.
func StartListen(port int, useTls bool, crtPath string, keyPath string, sname string) (error) {
    // Create a listening address
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return err
    }

    // start a new server and listen on the address
    var l net.Listener
    l, err = net.ListenTCP("tcp", addr)
    if err != nil {
        return err
    }

    // wrap with TLS if required
    if useTls {
        cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
        if err != nil {
            return err
        }
        conf := tls.Config{}

        certs := make([]tls.Certificate, 1)
        certs[0] = cert
        conf.Certificates = certs

        cp := x509.NewCertPool()
        caCert, err := ioutil.ReadFile(crtPath)
        if err != nil {
            return err
        }
        if !cp.AppendCertsFromPEM(caCert) {
            return errors.New("Could not append PEM cert")
        }
        conf.RootCAs = cp

        conf.ServerName = sname

        conf.ClientAuth = tls.RequireAndVerifyClientCert

        conf.ClientCAs = cp

        l = tls.NewListener(l, &conf)
    }

    // at the end of this function close the server connection
    defer l.Close()

    logging.Debug("Starting listen loop")
    for {
        a, err := acceptApp(l)
        if err != nil {
            return err
        } else {
            logging.Debug("Got connection")
            go ListenForCommands(a)
        }
    }
    return nil
}


// Listens for and accepts a new application connection
// Returns a reference to the App which was accepted
func acceptApp(lnr net.Listener) (*app.App, error) {
    conn, err := lnr.Accept()
    if err != nil {
        return nil, err
    }
    ret := app.NewApp()
    ret.Connection = conn
    return ret, nil
}


// Starts listening for commands through the given app's connection.
// Since this is a blocking function, it should likely be called in
// a goroutine.
// At the end of the function, closes the application's network resources.
func ListenForCommands(a *app.App) {
    defer closeApp(a)

    // loop forever consuming messages
    for {
        // get the next command
        strategy, err := getCommand(a)
        if err != nil {
            id := "NO_ID_FOUND"
            if a.Manifest != nil {
                id = a.Manifest.InstanceId
            }
            logging.Error("Application %s encountered fatal error: %s",
                          id,
                          err.Error())
            return
        }

        // do the command
        if strategy.GetVerb() == "APP_MANIFEST" {
            if !strategy.Execute() {
                logging.Error("Application's manifest did not parse correctly")
                return
            }
        } else {
            // HACK to get around golang compiler error
            // originally:
            //   go strategy.Execute()
            // message:
            //   "go requires function call not conversion"
            go performStrategy(strategy)
        }
    }
}


// Gets the next command from the application.
// Returns the command and an error, if one occured.
func getCommand(a *app.App) (cmd.CommandStrategy, error) {
    // get the message length
    msgLen, err := getMsgLen(a)
    if err != nil {
        return nil, err
    }

    // get the message body
    msgBuff := make([]byte, msgLen)
    totalRead := int64(0)
    // loop until we've read enough bytes
    for totalRead < msgLen {
        n, err := a.Connection.Read(msgBuff[totalRead:])
        if err != nil {
            return nil, err
        }
        totalRead += int64(n)
    }

    // we have the body, parse the command
    command, ok := cmd.ParseCommand(a, string(msgBuff))
    if !ok {
        return nil, errors.New("Command did not parse correctly")
    }
    return command, nil
}


// Gets the message length of the next message to be received by this application.
// Returns the message length and an error, if any occured.
func getMsgLen(a *app.App) (int64, error) {
    // create a small buffer to store the bytes read
    smallBuff := make([]byte, 200)
    smallBuffI := 0
    // loop while the buffer is not full
    for smallBuffI < len(smallBuff) {
        // read one byte
        innerBuff := make([]byte, 1)
        _, err := a.Connection.Read(innerBuff)
        if err != nil {
            return 0, err
        }
        // store that byte
        smallBuff[smallBuffI] = innerBuff[0]
        smallBuffI += 1
        // convert to string, see if it ends in newline
        str := string(smallBuff[:smallBuffI])
        if len(str) > 0 && strings.HasSuffix(str, "\n") {
            // this may be a valid message length
            var msgLen int64
            // parses using base 10, 64 bits
            msgLen, err = strconv.ParseInt(str[:len(str)-1], 10, 64)
            if err != nil {
                return 0, err
            }
            // it parsed, return
            return msgLen, nil
        }
    }
    return 0, errors.New("Message length exceeded 200 byte buffer.")
}


// Performs all operations required to close this app.
// Closes the network resource, queues a new STATUS_DOWN,
// removes from the registry, and logs the close.
func closeApp(a *app.App) {
    // this really should be somewhere else in the code, but i can't figure out where
    // since most places would lead to circular references
    a.Connection.Close()
    if a.Manifest != nil {
        cmd.ChangeAppStatus(a, app.STATUS_DOWN, 0)
        registry.Remove(a)
        logging.Info("Closing application %s", a.Manifest.InstanceId)
    } else {
        logging.Info("Closing application")
    }
}

// HACK to get around go compiler error
func performStrategy(s cmd.CommandStrategy) {
    s.Execute()
}
