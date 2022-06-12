package forms

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" binding:"required,mobile"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaId string `json:"captcha_id" binding:"required"`
}

type SendSmsForm struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
	Type   string `json:"type" binding:"required,oneof=register login"`
}

type RegisterForm struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Password string `json:"password" binding:"required"`
}
