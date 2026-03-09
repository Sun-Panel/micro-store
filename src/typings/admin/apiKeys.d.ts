
declare namespace Admin.ApiKeys {
    interface ApiKey extends Common.InfoBase{
        key ?:      string
        platform ?: number   
        note     ?: string
        name     ?: string
        status   ?: number   
        exception?: string
    }
}
