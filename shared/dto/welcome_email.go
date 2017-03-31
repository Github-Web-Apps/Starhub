package dto

// WelcomeEmailData is the DTO passed to the welcome email template
type WelcomeEmailData struct {
	Login                 string
	Email                 string
	Followers             int
	Stars                 int
	Repos                 int
	ChangeSubscriptionURL string
}
