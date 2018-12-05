package status

import (
	"context"
	"errors"
)

// PublicAPI represents a set of APIs from the `web3.status` namespace.
type PublicAPI struct {
	s *Service
}

// NewAPI creates an instance of the status API.
func NewAPI(s *Service) *PublicAPI {
	return &PublicAPI{s: s}
}

// LoginRequest : json request for status_login.
type LoginRequest struct {
	Addr     string `json:"address"`
	Password string `json:"password"`
}

// LoginResponse : json response returned by status_login.
type LoginResponse struct {
	AddressKeyID string `json:"address_key_id"`
}

// Login is an implementation of `status_login` or `web3.status.login` API
func (api *PublicAPI) Login(context context.Context, req LoginRequest) (res LoginResponse, err error) {
	_, accountKey, err := api.s.am.AddressToDecryptedAccount(req.Addr, req.Password)
	if err != nil {
		return
	}

	if res.AddressKeyID, err = api.s.w.AddKeyPair(accountKey.PrivateKey); err != nil {
		return
	}

	if err = api.s.am.SelectWalletAccount(req.Addr, req.Password); err != nil {
		return
	}

	return
}

// SignupRequest : json request for status_signup.
type SignupRequest struct {
	Password string `json:"password"`
}

// SignupResponse : json response returned by status_signup.
type SignupResponse struct {
	WalletAddress string `json:"address"`
	WalletPubKey  string `json:"pubkey"`
	ChatAddress   string `json:"chatAddress"`
	ChatPubKey    string `json:"chatPubkey"`
	Mnemonic      string `json:"mnemonic"`
}

// Signup is an implementation of `status_signup` or `web3.status.signup` API
func (api *PublicAPI) Signup(context context.Context, req SignupRequest) (*SignupResponse, error) {
	accountInfo, err := api.s.am.CreateAccount(req.Password)
	if err != nil {
		err = errors.New("could not create the specified account : " + err.Error())
		return nil, err
	}

	resp := &SignupResponse{
		WalletAddress: accountInfo.WalletKeyInfo.Address,
		WalletPubKey:  accountInfo.WalletKeyInfo.PubKey,
		ChatAddress:   accountInfo.ChatKeyInfo.Address,
		ChatPubKey:    accountInfo.ChatKeyInfo.PubKey,
		Mnemonic:      accountInfo.Mnemonic,
	}

	return resp, nil
}
