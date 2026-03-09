declare namespace Admin.RedeemCode {
    interface GetListRequest {
        page: number
        limit: number
        keyWord?: string
        status?: number
    }
}