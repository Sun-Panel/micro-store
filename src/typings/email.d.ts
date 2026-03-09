declare namespace Email{

    interface Info extends Common.InfoBase{
        title:string 
        content:string
        name:string
        note:string
    }

    interface SendEmailReq {
        mailTo:string[] 
        title:string
        body:string
        emailTemplateId?:number
        templateArg?:{[key:string]:string}
        replaceArg:boolean
        isSeparateSend:boolean
    }

}