declare namespace Openness.open {

    interface LoginConfigRegister {
        emailSuffix  :string   // 注册邮箱后缀
        openRegister :boolean    // 开放注册
    }
    
    interface LoginVcodeResponse{
        loginCaptcha: boolean
        register:LoginConfigRegister
    }

    interface HomeBase {
        logo_click_to_link?: string
        logo_text?: string
        logo_url?: string
        home_url?:string
    }
      
}