package provider

type Provider struct {
	BaseProvider  IProvider
	smsBulkSender ISmsBulkSender
	requestGet    IRequestMethodGet
	requestPost   IRequestMethodPost
	mapper        IMapper
	errorHandler  IErrorHandler
	authenticator IAuthenticator
	smsReceiver   ISmsReceiver
}

func ProviderServiceFactory(providerName string) *Provider {
	tpA := NewThirdPartyProviderA()
	tpB := NewThirdPartyProviderB()
	switch providerName {
	case tpA.GetName():
		return &Provider{
			BaseProvider:  tpA,
			requestGet:    tpA,
			requestPost:   tpA,
			mapper:        tpA,
			errorHandler:  tpA,
			authenticator: tpA,
			smsReceiver:   tpA,
		}
	case tpB.GetName():
		return &Provider{
			BaseProvider:  tpB,
			requestGet:    tpB,
			requestPost:   tpB,
			mapper:        tpB,
			errorHandler:  tpB,
			authenticator: tpB,
			smsReceiver:   tpB,
			smsBulkSender: tpB,
		}
	default:
		return nil
	}
}
