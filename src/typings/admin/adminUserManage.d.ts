declare namespace AdminUserManage {
    interface GetListRequest{
        page:number
        limit:number
        keyWord?:string
    }

    interface UpdatePasswordRequest{
        userId: number
        newPassword: string
    }
}