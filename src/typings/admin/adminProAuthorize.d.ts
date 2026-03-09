declare namespace Admin.ProAuthorize{

    interface ProAuthorizeUpdateUserExpiredTimeByDayReq {
        userId: number;
        dayNum: number;
        note: string;
        adminNote?:string
    }

    interface ProAuthorizeGetUserProAuthorizeListItemResp {
        userId: number;
        username: string;
        expiredTime: string; // 时间戳字符串
        name: string;
    }
}