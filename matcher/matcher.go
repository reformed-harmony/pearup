package matcher

import (
	"time"

	"github.com/reformed-harmony/pearup/db"
	"github.com/sirupsen/logrus"
)

// Matcher performs the pearup procedure.
type Matcher struct {
	conn        *db.Conn
	log         *logrus.Entry
	triggerChan chan bool
	stopChan    chan bool
	stoppedChan chan bool
}

func (m *Matcher) run() {
	defer close(m.stoppedChan)
	defer m.log.Info("matcher stopped")
	m.log.Info("matcher started")
	for {
		var (
			timerChan <-chan time.Time
			p         = &db.Pearup{}
			now       = time.Now()
		)
		err := m.conn.Transaction(func(conn *db.Conn) error {
			if db := conn.
				Order("end_date").
				Where("is_complete = ?", false).
				First(p); db.Error != nil {
				if !db.RecordNotFound() {
					return db.Error
				}
				m.log.Info("no pearups scheduled - waiting for trigger")
				return nil
			}
			if p.EndDate.After(now) {
				nextPearupTime := p.EndDate.Sub(now)
				m.log.Infof("waiting %s", nextPearupTime)
				timerChan = time.After(nextPearupTime)
				return nil
			}
			return m.match(conn, p)
		})
		if err != nil {
			m.log.Errorf("%s - pausing for 10 minutes", err.Error())
			timerChan = time.After(10 * time.Minute)
		}
		select {
		case <-timerChan:
		case <-m.triggerChan:
		case <-m.stopChan:
			return
		}
	}
}

// New creates a new matcher with the specified configuration.
func New(cfg *Config) *Matcher {
	m := &Matcher{
		conn:        cfg.Conn,
		log:         logrus.WithField("context", "matcher"),
		triggerChan: make(chan bool, 1),
		stopChan:    make(chan bool),
		stoppedChan: make(chan bool),
	}
	go m.run()
	return m
}

// Trigger indicates that the matcher should check for new pearups that have
// reached their end date and should be handled.
func (m *Matcher) Trigger() {
	select {
	case m.triggerChan <- true:
	default:
	}
}

// Close shuts down the matcher.
func (m *Matcher) Close() {
	close(m.stopChan)
	<-m.stoppedChan
}
