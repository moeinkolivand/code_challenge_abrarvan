package provider

import "fmt"

func (s *Provider) SendSMS(phone, message string) error {
	return s.BaseProvider.SendSMS(phone, message)
}

func (s *Provider) SendSMSBulk(numbers []string, message string) error {
	if s.smsBulkSender == nil {
		return fmt.Errorf("SendSMSBulk not supported")
	}
	return s.smsBulkSender.SendSMSBulk(numbers, message)
}

func (s *Provider) Get(endpoint string) (map[string]interface{}, error) {
	if s.requestGet == nil {
		return nil, fmt.Errorf("GET not supported")
	}
	return s.requestGet.Get(endpoint)
}

func (s *Provider) Post(endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	if s.requestPost == nil {
		return nil, fmt.Errorf("POST not supported")
	}
	return s.requestPost.Post(endpoint, data)
}

func (s *Provider) Authenticate(credentials map[string]string) error {
	if s.authenticator == nil {
		return fmt.Errorf("Authenticate not supported")
	}
	return s.authenticator.Authenticate(credentials)
}

func (s *Provider) ReceiveSMSDataByID(id string) (map[string]interface{}, error) {
	if s.smsReceiver == nil {
		return nil, fmt.Errorf("ReceiveSMSDataByID not supported")
	}
	return s.smsReceiver.ReceiveSMSDataByID(id)
}
