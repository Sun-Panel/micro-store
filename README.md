# Sun-Pay 授权平台


sun-store钩子

baseUrl= https://test_auth.sun-panel.top

订单相关：
{{baseUrl}}/sunStore/api/goodsOrder

用户相关
{{baseUrl}}/sunStore/api/user

单点退出钩子
{{baseUrl}}/api/ssoLogoutWebhook


sun-store授权端点：/oauth2/v1/authorize

迁移的流程

1. 每个新版程序优先运行一下退出
2. 根据脚本执行数据迁移
3. 创建三方应用
    - 名称：Sun-Panel 服务
    - 备注信息：授权、检查更新等
    - 查看详情获取到 应用id 和 key
    - 单点退出地址
4. 创建钩子地址参照上方两钩子，记录钩子的秘钥
5. 在网站设置中设置一下网站地址、并增加菜单
    菜单数据：
    ```
    [{"key":"1717554425923-0.5262344938783594","lable":"社区","title":"社区","url":"https://doc.sun-panel.top/zh_cn/introduce/author_groups.html"},{"key":"1717554451431-0.5376372518270507","lable":"演示站","title":"演示站","url":"https://sunpaneldemo.enianteam.com/#/"},{"key":"1717554397483-0.7320219115042685","lable":"升级到PRO","title":"升级到PRO","url":"#/pro"},{"children":[{"key":"1717554513270-0.19885156960259276","lable":"TG","title":"TG","url":"https://t.me/+bwOFXt6zXf43Njk1"},{"key":"1717554541399-0.05994032447784159","lable":"QQ","title":"QQ","url":"http://qm.qq.com/cgi-bin/qm/qr?_wv=1027\u0026k=yWCyKgcs2ybPwx-SyVWRX3bQgSEw9Sll\u0026authKey=yMgOqKG9jao5KHmbrjaccXeLewSTBP%2BBPJBcxymjIMGc6H5dq7H9EMnMXtJXugr4\u0026noverify=0\u0026group_code=831615449"}],"key":"1717554472170-0.21550926343116172","lable":"交流群","title":"交流群","url":""}]
    ```
    自建一个系统菜单：我的授权信息，http://127.0.0.1:1003/#/if?u=http://127.0.0.1:1004/platform/proAuthorize
6. 将应用信息和钩子秘钥写在授权服务的配置文件中
6. 重启服务,进入平台（使用账号密码登录）设置自站的地址，关闭新用户注册

创建各种页面 并调节好页面中的链接地址

商城平台设置首页 系统变量：http://127.0.0.1:1003/#/admin/systemVariable 首页地址：/#/hpage/home










商品限制数量


商品删除快照，
快照部分数据没有迁移，默认的货币代码
单笔订单多个数量没有通过钩子传输