package provider

import (
	"abrarvan_challenge/logging"
	"fmt"
)

type ThirdPartyProviderB struct {
	Provider
}

func NewThirdPartyProviderB() *ThirdPartyProviderB {
	return &ThirdPartyProviderB{}
}

func (tp *ThirdPartyProviderB) GetName() string {
	return "provider_two"
}

func (p *ThirdPartyProviderB) SendSMS(phoneNumber, message string) error {
	logger.Info(logging.SMSProvider, logging.Publish, "Sending SMS via Provider B", map[logging.ExtraKey]interface{}{
		"phone_number": phoneNumber,
		"message":      message,
	})
	return nil
}

func (p *ThirdPartyProviderB) Post(endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	return nil, p.HandleError(fmt.Errorf("POST not supported by Provider B"))
}

func (p *ThirdPartyProviderB) Get(endpoint string) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Request, "Making GET request to Provider B", map[logging.ExtraKey]interface{}{
		"endpoint": endpoint,
	})
	return map[string]interface{}{"status": "success"}, nil
}

func (p *ThirdPartyProviderB) MapToProviderFormat(data map[string]interface{}) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Mapping, "Mapping to Provider B format", nil)
	return data, nil
}

func (p *ThirdPartyProviderB) MapFromProviderFormat(data map[string]interface{}) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Mapping, "Mapping from Provider B format", nil)
	return data, nil
}

func (p *ThirdPartyProviderB) HandleError(err error) error {
	logger.Error(logging.SMSProvider, logging.Request, "Error in Provider B: "+err.Error(), nil)
	return fmt.Errorf("Provider B: %w", err)
}

func (p *ThirdPartyProviderB) Authenticate(credentials map[string]string) error {
	logger.Info(logging.SMSProvider, logging.Authentication, "Authenticating with Provider B", nil)
	return nil
}

func (p *ThirdPartyProviderB) ReceiveSMSDataByID(id string) (map[string]interface{}, error) {
	logger.Info(logging.SMSProvider, logging.Request, "Receiving SMS data by ID from Provider B", map[logging.ExtraKey]interface{}{
		"id": id,
	})
	return map[string]interface{}{"id": id, "data": "sample"}, nil
}

func (p *ThirdPartyProviderB) SendSMSBulk(phoneNumbers []string, message string) error {
	logger.Info(logging.SMSProvider, logging.Request, "Receiving SMS data by ID from Provider B", nil)
	return nil
}
