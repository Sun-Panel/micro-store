declare namespace Admin.GoodsManage {
    interface GetListReq extends Common.ListRequest{

    }

    interface CommonGoodsStruct {
        title?: string;              // 商品标题
        price?: number;              // 价格
        originalPrice?: number;      // 原价
        discount?: string;           // 优惠活动
        description?: string;        // 描述
        param?: string | object;              // 商品参数，JSON 字符串
    }

    interface GoodsInfo extends CommonGoodsStruct {
        id?:number
        sort?: number;
        status?: number;             // 上下架数据
        num?: number;
        lastSnapshotId?: number;
        userId?: number;
        createTime?:string
        updateTime?:string
    }
    
    interface GoodsGetListItemResp extends GoodsInfo {
    //    继承
    }

    interface UpdateSaleReq  {
        sort?: number;
        status?: boolean;             // 上下架数据
        num?: number;
        id?:number
    }
}