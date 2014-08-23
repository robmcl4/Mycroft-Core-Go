package cmd

func (f *failedCommandStrategy) generalFailure() {
    body := make(jsonData)
    body["recieved"] = recieved
    body["message"] = message
    f.app.Send("MSG_GENERAL_FAILURE", body)
}
