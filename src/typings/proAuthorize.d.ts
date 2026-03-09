declare namespace ProAuthorize {

    interface GetAuthorizeResp {
        expiredTime: string
    }

    interface GetAuthorizeHistoryRecordResp {
        changeTime: string
        expiredTime: string; // 时间戳字符串
        dayNum: number;
        note: string;
        orderNo: string;
        adminNote: string
    }


}