declare namespace HtmlPage {
    interface Info {
        isLogin?: boolean;
        content?: string;
        messageTemplateFlag?: string;
        messageTemplatePosition?: string; // top|bottom
    }

    interface EditReq extends HtmlPage.Info {
        pageName: string;
        pageDescription: string;
    }

    interface ListItem extends HtmlPage.Info {
        pageName: string;
        pageDescription: string;
    }

}
