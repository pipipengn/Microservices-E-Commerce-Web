package forms

type AddressForm struct {
	State        string `form:"state" json:"state" binding:"required"`
	City         string `form:"city" json:"city" binding:"required"`
	Postcode     string `form:"postcode" json:"postcode" binding:"required"`
	Address      string `form:"address" json:"address" binding:"required"`
	SignerName   string `form:"signer_name" json:"signer_name" binding:"required"`
	SignerMobile string `form:"signer_mobile" json:"signer_mobile" binding:"required"`
}
