package controller

import (
	"github.com/anjolabassey/Rave-go/rave"
	"github.com/yusuf/p-catalogue/model"
	"os"
)

var card = rave.Card{
	Rave: rave.Rave{
		Live:      false,
		PublicKey: os.Getenv("RAVE_PUBKEY"),
		SecretKey: os.Getenv("RAVE_SECKEY"),
	},
}

func (ct *Catalogue) Process(payload *model.PayLoad) (map[string]interface{}, error) {

	details := rave.CardChargeData{
		Cardno:        payload.CardNo,
		Cvv:           payload.Cvv,
		Expirymonth:   payload.ExpiryMonth,
		Expiryyear:    payload.ExpiryYear,
		Pin:           payload.Pin,
		Amount:        payload.Amount,
		Currency:      "NGN",
		CustomerPhone: payload.Phone,
		Firstname:     "",
		Lastname:      "",
		Email:         payload.Email,
		Txref:         payload.TxRef,
		RedirectUrl:   "https://localhost:8000/view-books",
	}
	err, resp := card.ChargeCard(details)
	if err != nil {
		panic(err)
	}
	return resp, nil
}

func (ct *Catalogue) Validate(ref, otp string) (map[string]interface{}, error) {
	validatePayload := rave.CardValidateData{
		Reference: ref,
		Otp:       "1234",
		PublicKey: os.Getenv("RAVE_PUBKEY"),
	}

	err, resp := card.ValidateCard(validatePayload)
	if err != nil {
		panic(err)
	}
	return resp, nil
}
