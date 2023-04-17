package rungcal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

//go:embed sql/insert.sql
var fs embed.FS

const (
	sqlFile  = "sql/insert.sql"
	timeZone = "Asia/Tokyo"
)

type InsertOption struct {
	Option
	ReCreate bool
}

type Execution struct {
	id            int
	uuid          string
	project       string
	job           string
	dateStarted   time.Time
	dateCompleted time.Time
	schedule      string
	executionTime int
}

func Insert(insertOption InsertOption) int {
	db := dbInit()
	defer db.Close()

	date, _ := time.ParseInLocation(time.DateOnly, insertOption.TargetDate, getLocation())
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
	log.Print("[INFO] === Targate Date(start ~ end) ===")
	log.Printf("[INFO] %s ~ %s", start.String(), end.String())

	sql, err := fs.ReadFile(sqlFile)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(string(sql), start.String(), end.String())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	log.Println("[INFO] === Select Database ===")
	events := make(map[int]*calendar.Event)
	for rows.Next() {
		var execution Execution
		err := rows.Scan(&execution.id, &execution.uuid, &execution.project, &execution.job, &execution.dateStarted, &execution.dateCompleted, &execution.schedule, &execution.executionTime)
		if err != nil {
			panic(err.Error())
		}

		if insertOption.Project != "" && insertOption.Project != execution.project {
			continue
		}

		event := &calendar.Event{
			Summary:     execution.job,
			Description: fmt.Sprintf("uuid: %s\nproject: %s\njob: %s", execution.uuid, execution.project, execution.job),
			Start:       &calendar.EventDateTime{DateTime: execution.dateStarted.Format(time.RFC3339), TimeZone: timeZone},
			End:         &calendar.EventDateTime{DateTime: execution.dateCompleted.Format(time.RFC3339), TimeZone: timeZone},
		}

		if check := checkGcal(execution); !check {
			continue
		}
		log.Printf("[DEBUG] %s %s %5d(s) %s", event.Start.DateTime, event.End.DateTime, execution.executionTime, event.Summary)

		events[execution.id] = event
	}

	credentials := []byte(getEnv("GOOGLE_SERVICE_ACCOUNT_CREDENTIALS_JSON"))
	config, err := google.JWTConfigFromJSON(credentials, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	wts := option.WithTokenSource(config.TokenSource(ctx))
	srv, err := calendar.NewService(ctx, wts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[INFO] === Insert Google Calendar ===")
	if insertOption.DryRun {
		for _, event := range events {
			log.Printf("[DEBUG] *** DRY RUN *** %s %s %s", event.Start.DateTime, event.End.DateTime, event.Summary)
		}
	} else {
		for _, event := range events {
			_, err = srv.Events.Insert(getEnv("GOOGLE_CALENDAR_ID"), event).Do()
			if err != nil {
				log.Fatalf("[ERROR] %v\n", err)
			}

			log.Printf("[DEBUG] %s %s %s", event.Start.DateTime, event.End.DateTime, event.Summary)
			// APIの叩き過ぎないように制限する
			time.Sleep(300 * time.Millisecond)
		}
	}

	return 0
}

func dbInit() *sql.DB {
	c := mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:3306", getEnv("DATABASE_HOST", "127.0.0.1")),
		DBName:               getEnv("DATABASE_NAME"),
		User:                 getEnv("DATABASE_USERNAME"),
		Passwd:               getEnv("DATABASE_PASSWORD"),
		Collation:            getEnv("DATABASE_COLLATION", "utf8mb4_0900_ai_ci"),
		Loc:                  getLocation(),
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getLocation() *time.Location {
	tz, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal(err)
	}

	return tz
}

// 実行間隔からGカレンダーに登録するか判断する
func checkGcal(row Execution) bool {
	if row.executionTime <= 0 {
		return false
	}

	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	s, err := p.Parse(row.schedule)
	if err != nil {
		log.Printf("[ERROR] Cron expression invalid syntax '%s' (%v)", row.schedule, err)
		return false
	}

	t1 := s.Next(time.Now())
	t2 := s.Next(t1)

	// ジョブ実行間隔が60分未満はGカレンダーに登録しない
	return t2.Sub(t1).Minutes() > 60
}
