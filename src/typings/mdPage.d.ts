declare namespace MdPage {
    interface Info {
        isLogin?: boolean;
        content?: string;
        messageTemplateFlag?: string;
        messageTemplatePosition?: string; // top|bottom
    }
    
    interface EditReq extends MdPage.Info {
        mdPageName: string;
        mdPageDescription: string;
    }
    
    interface ListItem extends MdPage.Info {
        mdPageName: string;
        mdPageDescription: string;
    }

}