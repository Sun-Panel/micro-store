declare namespace User{

	// 角色常量定义
	export const ROLE_USER = 1      // 普通用户
	export const ROLE_DEVELOPER = 2 // 开发者
	export const ROLE_ADMIN = 4     // 管理员
	export const ROLE_AUDITOR = 8   // 审核员
	export const ROLE_OPERATOR = 16 // 运营

	// 角色信息接口
	export interface RoleInfo {
		value: number
		key: string
		name: string
		description: string
	}

	interface Info{
		id?:number
		name ?:string
		createTime?:string
		username?:string
		password?:string
		headImage?:string
		status?:number
		role?:number
		mail?:string
		// userId?:string // id代替
		token?:string
		isAdmin?:number
		isBindSunStore?:bool
		roles?: RoleInfo[] // 角色标签列表
	}

	interface GetReferralCodeResponse{
		referralCode:string
	}

	
}
