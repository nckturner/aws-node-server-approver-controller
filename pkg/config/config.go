package config

// AWSNodeServerApproverOptions
type AWSNodeServerApproverConfig struct {
	Options AWSNodeServerApproverOptions

	//SecureServing  *apiserveroptions.SecureServingOptions
	//Authentication *apiserveroptions.DelegatingAuthenticationOptions
	//Authorization  *apiserveroptions.DelegatingAuthorizationOptions
}

type AWSNodeServerApproverOptions struct {
	Master     string
	Kubeconfig string
}
