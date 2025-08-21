package provider

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/logging"
	"fmt"
)

var logger logging.Logger = logging.NewLogger(config.GetConfig())

type ThirdPartyProviderA struct{}

func NewThirdPartyProviderA() *ThirdPartyProviderA {
	return &ThirdPartyProviderA{}
}

func (tp *ThirdPartyProviderA) GetName() string {
	return "provider_one"
}

func (p *ThirdPartyProviderA) SendSMS(phoneNumber, message string) error {
	logger.Info(logging.SMSProvider, logging.Publish, "Sending SMS via Provider A", map[logging.ExtraKey]interface{}{
		"phone_number": phoneNumber,
		"message":      message,
	})
	return nil
}

func (p *ThirdPartyProviderA) Post(endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	return nil, p.HandleError(fmt.Errorf("POST not supported by Provider A"))
}

func (p *ThirdPartyProviderA) Get(endpoint string) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Request, "Making GET request to Provider A", map[logging.ExtraKey]interface{}{
		"endpoint": endpoint,
	})
	return map[string]interface{}{"status": "success"}, nil
}

func (p *ThirdPartyProviderA) MapToProviderFormat(data map[string]interface{}) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Mapping, "Mapping to Provider A format", nil)
	return data, nil
}

func (p *ThirdPartyProviderA) MapFromProviderFormat(data map[string]interface{}) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Mapping, "Mapping from Provider A format", nil)
	return data, nil
}

func (p *ThirdPartyProviderA) HandleError(err error) error {
	logger.Error(logging.SMSProvider, logging.Request, "Error in Provider A: "+err.Error(), nil)
	return fmt.Errorf("Provider A: %w", err)
}

func (p *ThirdPartyProviderA) Authenticate(credentials map[string]string) error {
	logger.Info(logging.SMSProvider, logging.Authentication, "Authenticating with Provider A", nil)
	return nil
}

func (p *ThirdPartyProviderA) ReceiveSMSDataByID(id string) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Request, "Receiving SMS data by ID from Provider A", map[logging.ExtraKey]interface{}{
		"id": id,
	})
	return map[string]interface{}{"id": id, "data": "sample"}, nil
}
