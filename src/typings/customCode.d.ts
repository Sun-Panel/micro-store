declare namespace CustomCode {
    // 代码类型（块内单选）
    const enum CodeType {
        JS = 1,
        CSS = 2,
        Footer = 3,
    }

    // 适用的客户端版本
    const enum ClientVersion {
        V1 = 1,
        V2 = 2,
    }

    // 主表线上状态
    const enum Status {
        Draft = -1,
        Pending = 0,
        Online = 1,
        Rejected = 2,
        Offline = 3,
    }

    // 审核快照状态
    const enum ReviewStatus {
        Draft = -1,
        Pending = 0,
        Approved = 1,
        Rejected = 2,
    }

    // 机器预审结果（提交时由系统生成，供审核员参考未通过原因）
    interface MachineAudit {
        /** 是否自动审核通过 */
        auto: boolean
        /** 作者是否受信任（近一年内有审核通过记录） */
        trusted: boolean
        /** 静态安全扫描是否通过 */
        scanPass: boolean
        /** 未通过/降级原因（人工审核时展示给审核员） */
        reasons: string[]
        /** 预审时间 */
        time?: string
    }

    // 限制
    const enum Limit {
        MaxBlocks = 10,
        MaxImagesPerBlock = 5,
        MaxImageSize = 524288, // 512K
    }

    // 代码片块
    interface Block {
        id?: number
        /** 1-JS 2-CSS 3-页脚 */
        codeType: number
        title: string
        note: string
        code: string
        images: string[]
        sort?: number
        /** 块唯一标识（复制粘贴去重用，开发者可改） */
        onlyId?: string
    }

    interface ListItem {
        id: number
        title: string
        description: string
        keywords: string[]
        isOriginal: boolean
        sourceNote: string
        versions: number[]
        codeTypes: number[]
        authorId: number
        authorName?: string
        status: number
        readCount: number
        publishedAt?: string
        createTime?: string
        reviewStatus?: number | null
        reviewNote?: string
    }

    interface EditReq {
        /** 0 或不传表示新增 */
        id?: number
        title: string
        description: string
        keywords: string[]
        isOriginal: boolean
        sourceNote: string
        versions: number[]
        blocks: Block[]
        /** true 提交审核，false 保存草稿 */
        submit: boolean
    }

    interface Info {
        id: number
        title: string
        description: string
        keywords: string[]
        isOriginal: boolean
        sourceNote: string
        versions: number[]
        codeTypes: number[]
        blocks: Block[]
        status: number
        readCount: number
        publishedAt?: string
        authorId: number
        authorName?: string
        /** 待审核时不可编辑，需先撤回 */
        canEdit: boolean
        reviewStatus?: number | null
        reviewNote?: string
        pendingReviewId?: number | null
    }

    interface Detail {
        id: number
        title: string
        description: string
        keywords: string[]
        isOriginal: boolean
        sourceNote: string
        versions: number[]
        codeTypes: number[]
        blocks: Block[]
        authorId: number
        authorName?: string
        readCount: number
        publishedAt?: string
    }

    interface ListRequest {
        page: number
        limit: number
        keyword?: string
        sortBy?: string
        sortOrder?: string
        /** 版本过滤：0-全部 */
        version?: number
        /** 代码类型过滤：0-全部 */
        codeType?: number
    }

    interface ReviewListItem {
        id: number
        customCodeId: number
        title: string
        description: string
        keywords: string[]
        isOriginal: boolean
        sourceNote: string
        codeTypes: number[]
        status: number
        authorName?: string
        createTime?: string
        onlineStatus: number
        /** 机器预审结果（含未通过原因） */
        machineAudit?: CustomCode.MachineAudit | null
    }

    interface ReviewDetail {
        review: {
            id: number
            customCodeId: number
            title: string
            description: string
            keywords: string[]
            isOriginal: boolean
            sourceNote: string
            versions: number[]
            codeTypes: number[]
            blocks: Block[]
            status: number
            reviewNote?: string
            createTime?: string
            /** 机器预审结果（含未通过原因） */
            machineAudit?: CustomCode.MachineAudit | null
        }
        online?: {
            id: number
            title: string
            description: string
            keywords: string[]
            isOriginal: boolean
            sourceNote: string
            versions: number[]
            codeTypes: number[]
            blocks: Block[]
            status: number
            publishedAt?: string
            authorName?: string
        }
    }
}
