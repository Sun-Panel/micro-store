package scheduler

import (
	"sun-panel/global"
	"sun-panel/scheduler/clientHistoryData"

	"github.com/robfig/cron/v3"
)

// 定义定时任务
func defineCorn(c *cron.Cron) {

	// 优先执行一次24小时统计,每小时执行一次历史客户端数据统计
	go func() {
		clientHistoryData.StartStatisticsHourAgo(24)
		create(c, "0 */1 * * *", clientHistoryData.Start)
	}()

	// // 每天执行一次不活跃的客户端数据清理
	// create(c, "22 5 * * *", func() {
	// 	fmt.Println("清理客户端不活跃数据:", time.Now())
	// })

	// create(c, "*/2 * * * *", func() {
	// 	fmt.Println("执行定时任务/2:", time.Now())
	// })

}

// 开始定时任务
func Start() {
	global.Logger.Infoln("====", "Start scheduler")
	c := cron.New() // 创建一个 cron 实例

	cron.NewChain()

	defineCorn(c)

	// 启动 Cron 调度
	c.Start()
	defer c.Stop()

	// 主线程保持运行
	select {}
	// select {
	// case <-time.After(10 * time.Second):
	// 	global.Logger.Infoln("====", "Start scheduler")
	// }

}

// 创建一个定时任务
func create(c *cron.Cron, spec string, cmd func()) (cron.EntryID, error) {
	return c.AddFunc(spec, cmd)
}

// 创建并执行一次任务
func createAndRun(c *cron.Cron, spec string, cmd func()) (cron.EntryID, error) {
	cmd()
	return create(c, spec, cmd)
}
