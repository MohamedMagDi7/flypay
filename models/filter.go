package models

//we make sure to add custom filter fields for each provider
type Filter struct {
	CustomAFilter
	CustomBFilter
	Provider  string `schema:"provider"`
	Status    string `schema:"statusCode"`
	AmountMin int    `schema:"amountMin"`
	AmountMax int    `schema:"amountMax"`
	Currency  string `schema:"currency"`
}

//custom filter fields for provider A
type CustomAFilter struct {
	AStatusCode int
}

//custom filter fields for provider BA
type CustomBFilter struct {
	BStatusCode int
}
