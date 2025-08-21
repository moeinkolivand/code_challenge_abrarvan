// package provider

// import (
// 	"abrarvan_challenge/config"
// 	"abrarvan_challenge/logging"
// 	"fmt"
// )

// var logger = logging.NewLogger(config.GetConfig())

// func (s *Provider) SendSMS(phoneNumber, message string) error {
//     if s.smsSender == nil {
//         err := fmt.Errorf("SMS sending not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Publish, err.Error(), nil)
//         return s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Publish, "Sending SMS", map[logging.ExtraKey]interface{}{
//         "phone_number": phoneNumber,
//         "message":      message,
//     })
//     if err := s.authenticator.Authenticate(map[string]string{"key": "value"}); err != nil {
//         logger.Error(logging.SMSProvider, logging.Authentication, "Authentication failed: "+err.Error(), nil)
//         return s.errorHandler.HandleError(err)
//     }
//     if err := s.smsSender.SendSMS(phoneNumber, message); err != nil {
//         logger.Error(logging.SMSProvider, logging.Publish, "Failed to send SMS: "+err.Error(), nil)
//         return s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Publish, "SMS sent successfully", map[logging.ExtraKey]interface{}{
//         "phone_number": phoneNumber,
//     })
//     return nil
// }

// func (s *Provider) SendSMSBulk(phoneNumbers []string, message string) error {
//     if s.smsBulkSender == nil {
//         err := fmt.Errorf("bulk SMS not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Publish, err.Error(), nil)
//         return s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Publish, "Sending bulk SMS", map[logging.ExtraKey]interface{}{
//         "phone_numbers": phoneNumbers,
//         "message":       message,
//     })
//     if err := s.authenticator.Authenticate(map[string]string{"key": "value"}); err != nil {
//         logger.Error(logging.SMSProvider, logging.Authentication, "Authentication failed: "+err.Error(), nil)
//         return s.errorHandler.HandleError(err)
//     }
//     if err := s.smsBulkSender.SendSMSBulk(phoneNumbers, message); err != nil {
//         logger.Error(logging.SMSProvider, logging.Publish, "Failed to send bulk SMS: "+err.Error(), nil)
//         return s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Publish, "Bulk SMS sent successfully", map[logging.ExtraKey]interface{}{
//         "phone_numbers": phoneNumbers,
//     })
//     return nil
// }

// func (s *Provider) Get(endpoint string) (map[string]interface{}, error) {
//     if s.requestGet == nil {
//         err := fmt.Errorf("GET not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Request, err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Request, "Making GET request", map[logging.ExtraKey]interface{}{
//         "endpoint": endpoint,
//     })
//     if err := s.authenticator.Authenticate(map[string]string{"key": "value"}); err != nil {
//         logger.Error(logging.SMSProvider, logging.Authentication, "Authentication failed: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     data, err := s.requestGet.Get(endpoint)
//     if err != nil {
//         logger.Error(logging.SMSProvider, logging.Request, "Failed to make GET request: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Request, "GET request successful", map[logging.ExtraKey]interface{}{
//         "endpoint": endpoint,
//     })
//     return data, nil
// }

// func (s *Provider) Post(endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
//     if s.requestPost == nil {
//         err := fmt.Errorf("POST not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Request, err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Request, "Making POST request", map[logging.ExtraKey]interface{}{
//         "endpoint": endpoint,
//     })
//     if err := s.authenticator.Authenticate(map[string]string{"key": "value"}); err != nil {
//         logger.Error(logging.SMSProvider, logging.Authentication, "Authentication failed: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     response, err := s.requestPost.Post(endpoint, data)
//     if err != nil {
//         logger.Error(logging.SMSProvider, logging.Request, "Failed to make POST request: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Request, "POST request successful", map[logging.ExtraKey]interface{}{
//         "endpoint": endpoint,
//     })
//     return response, nil
// }

// func (s *Provider) MapToProviderFormat(data map[string]interface{}) (map[string]interface{}, error) {
//     if s.mapper == nil {
//         err := fmt.Errorf("mapping to provider format not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Mapping, err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Mapping, "Mapping to provider format", nil)
//     mappedData, err := s.mapper.MapToProviderFormat(data)
//     if err != nil {
//         logger.Error(logging.SMSProvider, logging.Mapping, "Failed to map to provider format: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Mapping, "Mapped to provider format successfully", nil)
//     return mappedData, nil
// }

// func (s *Provider) MapFromProviderFormat(data map[string]interface{}) (map[string]interface{}, error) {
//     if s.mapper == nil {
//         err := fmt.Errorf("mapping from provider format not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Mapping, err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Mapping, "Mapping from provider format", nil)
//     mappedData, err := s.mapper.MapFromProviderFormat(data)
//     if err != nil {
//         logger.Error(logging.SMSProvider, logging.Mapping, "Failed to map from provider format: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Mapping, "Mapped from provider format successfully", nil)
//     return mappedData, nil
// }

// func (s *Provider) ReceiveSMSDataByID(id string) (map[string]interface{}, error) {
//     if s.smsReceiver == nil {
//         err := fmt.Errorf("receiving SMS data by ID not supported by provider")
//         logger.Error(logging.SMSProvider, logging.Request, err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Request, "Receiving SMS data by ID", map[logging.ExtraKey]interface{}{
//         "id": id,
//     })
//     if err := s.authenticator.Authenticate(map[string]string{"key": "value"}); err != nil {
//         logger.Error(logging.SMSProvider, logging.Authentication, "Authentication failed: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     data, err := s.smsReceiver.ReceiveSMSDataByID(id)
//     if err != nil {
//         logger.Error(logging.SMSProvider, logging.Request, "Failed to receive SMS data: "+err.Error(), nil)
//         return nil, s.errorHandler.HandleError(err)
//     }
//     logger.Info(logging.SMSProvider, logging.Request, "SMS data received successfully", map[logging.ExtraKey]interface{}{
//         "id": id,
//     })
//     return data, nil
// }

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
