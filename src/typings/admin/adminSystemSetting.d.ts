declare namespace AdminSystemSetting {
    interface Email{
        host: string
        port: number
        mail: string
        password: string
    }
    interface Website{
        loginCaptcha: boolean
        openRegister: boolean
        webSiteUrl: string
        emailSuffix:string
    }
}