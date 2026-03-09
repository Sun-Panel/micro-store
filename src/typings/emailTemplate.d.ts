declare namespace EmailTemplate{

    interface Info extends Common.InfoBase{
        title:string 
        content:string
        flag:string
        name:string
        note:string
        args?:Arg[]
        templateArg?:{[key:string]:string}
        isSeparateSend:boolean
    }

    interface Arg{
        name: string
        keyword: string
    }


}