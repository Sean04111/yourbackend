// Code generated by goctl. DO NOT EDIT.
package types

type Baseinfor struct {
	Status int  `json:"status"`
	Info   Info `json:"info"`
}

type Info struct {
	AvatarLink string `json:"avatarLink"`
	UserName   string `json:"userName"`
	Profession string `json:"profession"`
	Usermail   string `json:"userMail"`
	Type       string `json:"type"`
}

type Loginreq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Loginresp struct {
	Status      int    `json:"int"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Expires     string `json:"expires"`
}

type Pubkeyresp struct {
	Status int    `json:"status"`
	Pubkey string `json:"pubkey"`
}

type RefreshTokenreq struct {
	Email string `json:"email"`
}

type RefreshTokenresp struct {
	Status  int    `json:"status"`
	Name    string `json:"name"`
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

type Registerreq struct {
	Email string `json:"email"`
	Name  string `json:"name,optional"`
	Pass  string `json:"pass"`
	Code  string `json:"code"`
	Check string `json:"check"`
}

type Registerresp struct {
	Status      int    `json:"status"`
	Accesstoken string `json:"accessToken"`
	Expires     string `json:"expires"`
	Name        string `json:"name"`
}

type Settingavareq struct {
}

type Settingavaresp struct {
	Status int `json:"status"`
}

type Settingbasereq struct {
	Name       string `json:"name"`
	Profession string `json:"profession,optional"`
	Type       string `json:"type,optional"`
}

type Settingbaseresp struct {
	Status int `json:"status"`
}

type Updatepwdreq struct {
	Code     string `json:"code"`
	Check    string `json:"check"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Updatepwdresp struct {
	Status int `json:"status"`
}

type Codereq struct {
	Email string `json:"email"`
}

type Coderesp struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
}