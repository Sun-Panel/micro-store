declare namespace Admin.Dashboard {
    interface GetStatisticsResp {
        userCount: number
        userToday: number
        clientCount: number
        clientToday: number
        orderCount: number
        orderToday: number
        countCNY: number
        todayCNY: number
        countUSD: number
        todayUSD: number
    }
}