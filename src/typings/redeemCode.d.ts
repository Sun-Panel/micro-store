declare namespace RedeemCode {

    // 扩展数据

    // Pro授权扩展数据
    interface ExtendDataPro {
        days: number;
    }

    // 详情
    interface Info {
        writeOffTime?: string
        redeemType:number
        expireTime: string;
        releaseType: number;
        title: string;
        note: string;
        extendData: ExtendDataPro;
        code: string
        createTime:string
        status:number
        userInfo?: User.Info
    }

    interface CreateReq extends Omit<Info,'code'|'createTime'| 'writeOffTime' | 'releaseType' | 'status' > {
        prefix: string
        number: number
    }

    interface GetListRequest {
        page: number
        limit: number
        keyWord?: string
        status?: number
    }





}