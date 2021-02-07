package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	API_URL    = "/api/v1"
	PROXY_ADDR = "209.127.191.180:9279"
)

const (
	_cqlshrc_host = "df73af29-9ac6-4c7a-8854-aee16d36ae39-europe-west1.db.astra.datastax.com"
	_cqlshrc_port = "32539"
	_username     = "leaktrack_admin"
	_password     = "YL@C5Xn4ayBflRNxmW9yIO@ktpzLqU2iM7m"
)

type Error struct {
	Message string `json:"message"`
}

func main() {
	cluster := gocql.NewCluster(_cqlshrc_host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: _username,
		Password: _password,
	}
	cluster.Hosts = []string{_cqlshrc_host + ":" + _cqlshrc_port}

	certPath, _ := filepath.Abs("C:\\Users\\andre\\GolandProjects\\The_New_414s\\web_backend\\datastax_certs\\cert")
	keyPath, _ := filepath.Abs("C:\\Users\\andre\\GolandProjects\\The_New_414s\\web_backend\\datastax_certs\\key")
	caPath, _ := filepath.Abs("C:\\Users\\andre\\GolandProjects\\The_New_414s\\web_backend\\datastax_certs\\ca.crt")
	cert, _ := tls.LoadX509KeyPair(certPath, keyPath)
	caCert, _ := ioutil.ReadFile(caPath)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	cluster.SslOpts = &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: false,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	e := echo.New()

	// Create Tracker
	e.POST(API_URL+"/tracker", func(c echo.Context) error {
		// Generate an RFC 4122 UUID.
		_uuid := uuid.New().String()

		err := session.Query("INSERT INTO triggers.trackers ( uuid , created_at ) VALUES ( ? , toUnixTimestamp(now()) ) IF NOT EXISTS", uuid).WithContext(ctx).Exec()
		if err != nil {
			e.Logger.Error(err)

			errMsg := &Error{err.Error()}
			return c.JSON(http.StatusInternalServerError, errMsg)
		}
		return c.JSON(http.StatusCreated, _uuid)
	})

	type Event struct {
		UUID           string    `json:"uuid" cql:"uuid"`
		Data           string    `json:"data" cql:"data"`
		FromType       int16     `json:"from_type" cql:"from_type"`
		MatchedAt      time.Time `json:"matched_at" cql:"matched_at"`
		ServiceFrom    string    `json:"service_from" cql:"service_from"`
		TriggerMatched string    `json:"trigger_matched" cql:"trigger_matched"`
	}

	// View Tracker Events
	e.GET(API_URL+"/tracker/:id/events", func(c echo.Context) error {
		Qry := session.Query("SELECT uuid, data, from_type, matched_at, service_from, trigger_matched FROM triggers.events WHERE uuid = ?",
			c.Param("id")).WithContext(ctx).Iter()
		defer Qry.Close()

		var _events []Event
		_ev := Event{}

		for Qry.Scan(&_ev.UUID, &_ev.Data, &_ev.FromType, &_ev.MatchedAt, &_ev.ServiceFrom, &_ev.TriggerMatched) {
			_events = append(_events, _ev)
			_ev = Event{}
		}

		err := Qry.Close()
		if err != nil {
			e.Logger.Error(err)

			errMsg := &Error{err.Error()}
			return c.JSON(http.StatusInternalServerError, errMsg)
		}

		return c.JSON(http.StatusOK, _events)
	})

	// Delete Tracker
	e.DELETE(API_URL+"/tracker/:id", func(c echo.Context) error {
		err := session.Query("DELETE FROM triggers.trackers WHERE uuid = ? IF EXISTS", c.Param("id")).WithContext(ctx).Exec()
		if err != nil {
			e.Logger.Error(err)

			errMsg := &Error{err.Error()}
			return c.JSON(http.StatusInternalServerError, errMsg)
		}
		return c.NoContent(http.StatusNoContent)
	})

	// Create Trigger
	e.POST(API_URL+"/tracker/:id/trigger", func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	})

	// View Trigger
	e.GET(API_URL+"/tracker/:id/trigger/:str", func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	})

	// Update Trigger
	e.PUT(API_URL+"/tracker/:id/trigger/:str", func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	})

	// Delete Trigger
	e.DELETE(API_URL+"/tracker/:id/trigger/:str", func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	})

	// Process Content
	e.POST(API_URL+"/process", func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	})

	defer session.Close()
	e.Logger.Fatal(e.Start(":1323"))

}
