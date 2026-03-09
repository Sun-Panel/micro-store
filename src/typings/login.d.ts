declare namespace Login{

    interface LoginReqest{
        username:string 
        password:string
        vcode?:string
    }

	interface LoginResponse extends User.Info{
		token :string
	}

    interface ResetPasswordByVCodeReqest extends System.Register.SendRegisterVcodeRquest{
    }


    interface OAuth2CodeLoginResq {
        token: string;
        username: string;
        name: string;
        headImage: string;
        role: number;
        mail: string;
    }
    

}