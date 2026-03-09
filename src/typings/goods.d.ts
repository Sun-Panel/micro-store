declare namespace Goods{

    interface CommonGoodsStruct {
        id:number
        title: string;              // 商品标题
        price: number;              // 价格
        originalPrice: number;      // 原价
        discount: string;           // 优惠活动
        description: string;        // 描述
        param: string;              // 商品参数，JSON 字符串
    }


    interface GoodsInfo extends CommonGoodsStruct {
        // flag: string;                      // 商品唯一标识
        status: number;                    // 状态：参考上方全局变量GOODS_STATUS_XXXX
        sort: number;                      // 排序 数字越大排序越高
        num: number;                       // 库存商品数量
        lastSnapshotId: number;            // 最后快照 ID
        userId: number;                    // 用户 ID
    }
    
    interface GoodsSnapshotInfo extends CommonGoodsStruct {
      
    }
    
}