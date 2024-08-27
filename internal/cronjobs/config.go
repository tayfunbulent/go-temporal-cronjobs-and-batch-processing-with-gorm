package cronjobs

import (
	"go.temporal.io/sdk/client"
)

type CronSchedule struct {
	Schedule string
}

var (
	EveryMinute        = CronSchedule{"* * * * *"}    // Every minute
	Every2Minutes      = CronSchedule{"*/2 * * * *"}  // Every 2 minutes
	Every10Minutes     = CronSchedule{"*/10 * * * *"} // Every 10 minutes
	Every30Minutes     = CronSchedule{"*/30 * * * *"} // Every 30 minutes
	EveryHour          = CronSchedule{"0 * * * *"}    // At the beginning of every hour
	EveryDayMidnight   = CronSchedule{"0 0 * * *"}    // Every day at midnight
	EveryDayNoon       = CronSchedule{"0 12 * * *"}   // Every day at noon
	EveryWeekSunday    = CronSchedule{"0 0 * * 0"}    // Every week on Sunday at midnight
	EveryWeekMonday    = CronSchedule{"0 0 * * 1"}    // Every week on Monday at midnight
	EveryMonthStart    = CronSchedule{"0 0 1 * *"}    // On the 1st of every month at midnight
	EveryMonthMidnight = CronSchedule{"0 0 1-15 * *"} // On the 1st through 15th of every month at midnight
	EveryWeekday       = CronSchedule{"0 0 * * 1-5"}  // Every weekday at midnight (Monday to Friday)
	EveryWeekend       = CronSchedule{"0 0 * * 6,0"}  // Every weekend (Saturday and Sunday) at midnight
)

func ConfigureCronJobs(c client.Client) {
	sendEmail(c, EveryMinute.Schedule)
}
