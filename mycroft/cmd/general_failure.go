package cmd

func (f *failedCommandStrategy) generalFailure() {
    body := make(jsonData)
    body["recieved"] = f.received
    body["message"] = f.message
    f.app.Send("MSG_GENERAL_FAILURE", body)
}
