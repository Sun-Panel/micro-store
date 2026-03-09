declare namespace SystemVariable {

    interface SystemVariableListItem {
        id: number;
        description: string;
        configName: string;
        configValue: string;
    }
    
    interface SystemVariable {
        description: string;
        name: string;
        value: any; // 使用 any 表示接口类型
    }
    
    interface SystemVariableEditReq {
        id?: number;
        description: string;
        name: string;
        value: string;
    }
    
}