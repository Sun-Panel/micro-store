declare namespace GoodsOrder{

    interface Details {
        status             :number 
        orderNo            ?:string 
        payPlatform            ?:number 
        payPlatformOrderNo ?:string 
        countPrice         :number 
        goodsId :number
        goods:Shop.Details
        
    
    }

    interface CreateRequest {
        goodsSnapshotId:number 
        number?:number
        goodsId:number
    }

    interface CreateResponse {
        goodsId:number 
        number?:number
    }

    interface PayResponse{
        payUrl:string
    }

    interface Info {
        status: number; 
        orderNo: string;
        payPlatform: number;
        payPlatformOrderNo: string;
        currencyCode:string
        countPrice: number;
        goodsSnapshotId: number;
        createTime: string;
        payTime: string;
        goodsSnapshotInfo:Goods.GoodsSnapshotInfo
        goodsSnapshot:Goods.GoodsSnapshotInfo
        user:User.Info
        adminNote:string
        goods:Goods.GoodsInfo
        payUrl:string
    }
    

	
}