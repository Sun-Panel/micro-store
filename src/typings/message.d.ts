declare namespace Message {


    interface SendReq {
        toUsernames: string[]
        appendixs: string[]
        title?: string
        content: string
        templateArg: object
        templateFlag?: string
    }

    interface TemplateInfo extends Common.InfoBase {
        appendixs?: string[]
        title?: string
        titleShow?: number
        titleEdit?: number
        name?: string
        flag?: string
        content?: string
        contentShow?: number
        contentEdit?: number
        note?: string
        args?: Arg[]
        toUsernames: string[] | string
        toUsernameShow?: number
        toUsernameEdit?: number
        uploadAppendixShow?: number
        uploadAppendixMax?: number
        uploadAppendixMin?: number
        saveAfterUrl?: string
        saveButtonName?: string
    }

    interface Arg {
        name: string
        keyword: string
        required: boolean
        value?: number | string // 前端独有
        placeholder: string
        formType: string
    }


    interface MessageInfo {
        messageId: number;
        title: string;
        content: string;
        fromUserId: number;
        toUserId: number;
        appendix: string[];
        haveRead: number;
        topicId: string;
        templateFlag: string;
        fromUser?: User.Info;
        toUser?: User.Info;
        createTime: string
    }

}