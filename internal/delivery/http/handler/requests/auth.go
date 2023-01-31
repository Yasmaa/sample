package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type RegisterRequest struct {
	Firstname string `json:"first_name,omitempty"`
	Lastname  string `json:"last_name,omitempty"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type ResetLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=8,max=20"`
	Confirm  string `json:"confirm" validate:"required,min=8,max=20"`
	Token    string `json:"token" validate:"required"`
}

type UpdateUserInfoRequest struct {
	Firstname      string `json:"first_name,omitempty"`
	Lastname       string `json:"last_name,omitempty"`
	CustomerId     string `json:"customer_id,omitempty"`
	SubscriptionId string `json:"subscription_id,omitempty"`
	Country        string `json:"country,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	TwoFa     bool `json:"two_fa"`
	Secret     string `json:"secret"`


}

type UpdateUserPasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=20"`
}

type TwoFaToken struct {
	Token string `json:"token" validate:"required"`
	Email string `json:"email"`
}
