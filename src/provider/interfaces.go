package provider



type IProvider interface {
    SendSMS(phoneNumber, message string) error
    GetName() string
}


type ISmsBulkSender interface {
    SendSMSBulk(phoneNumbers []string, message string) error
}


type IRequestMethodPost interface {
    Post(endpoint string, data map[string]interface{}) (map[string]interface{}, error)
}


type IRequestMethodGet interface {
    Get(endpoint string) (map[string]interface{}, error)
}


type IMapper interface {
    MapToProviderFormat(data map[string]interface{}) (map[string]interface{}, error)
    MapFromProviderFormat(data map[string]interface{}) (map[string]interface{}, error)
}


type IErrorHandler interface {
    HandleError(err error) error
}


type IAuthenticator interface {
    Authenticate(credentials map[string]string) error
}


type ISmsReceiver interface {
    ReceiveSMSDataByID(id string) (map[string]interface{}, error)
}