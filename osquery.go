package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type OsqueryEnrollPostBody struct {
	EnrollSecret   string `json:"enroll_secret" binding:"required"`
	HostIdentifier string `json:"host_identifier" binding:"required"`
}

type OsqueryConfigPostBody struct {
	NodeKey string `json:"node_key" binding:"required"`
}

type OsqueryLogPostBody struct {
	NodeKey string                   `json:"node_key" binding:"required"`
	LogType string                   `json:"log_type" binding:"required"`
	Data    []map[string]interface{} `json:"data" binding:"required"`
}

type OsqueryResultLog struct {
	Name           string            `json:"name"`
	HostIdentifier string            `json:"hostIdentifier"`
	UnixTime       string            `json:"unixTime"`
	CalendarTime   string            `json:"calendarTime"`
	Columns        map[string]string `json:"columns"`
	Action         string            `json:"action"`
}

type OsqueryStatusLog struct {
	Severity string `json:"severity"`
	Filename string `json:"filename"`
	Line     string `json:"line"`
	Message  string `json:"message"`
	Version  string `json:"version"`
}

type OsqueryDistributedReadPostBody struct {
	NodeKey string `json:"node_key" binding:"required"`
}

type OsqueryDistributedWritePostBody struct {
	NodeKey string                         `json:"node_key" binding:"required"`
	Queries map[string][]map[string]string `json:"queries" binding:"required"`
}

func genNodeKey() string {
	return generateRandomText(12)
}

func setError(c *gin.Context, err error) {
	logrus.WithError(err).Error("Returning 500")
	c.AbortWithError(500, err)
	// c.JSON(http.StatusInternalServerError,
	// 	gin.H{
	// 		"error": err.Error(),
	// 	})
}

func OsqueryEnroll(c *gin.Context) {
	var body OsqueryEnrollPostBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Debugf("Error parsing OsqueryEnroll POST body: %s", err.Error())
		return
	}
	logrus.Debugf("OsqueryEnroll: %s %s", body.EnrollSecret, body.HostIdentifier)

	if body.EnrollSecret != config.Osquery.EnrollSecret {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"node_invalid": true,
			})
		return

	}

	db := mustOpenDB(config.MySQL.Username, config.MySQL.Password, config.MySQL.Address, config.MySQL.Database)
	defer db.Close()

	var host Host
	err = db.Debug().
		Where("uuid = ? OR host_name = ?",
			body.HostIdentifier, body.HostIdentifier).
		First(&host).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		DatabaseError(c)
		return

	} else if err == gorm.ErrRecordNotFound {
		// Create new Host
		host = Host{HostName: body.HostIdentifier, UUID: body.HostIdentifier}
	}

	host.NodeKey = genNodeKey()

	if err = db.Debug().Save(&host).Error; err != nil {
		setError(c, err)
		return
	}

	logrus.Debugf("Host: %+v", host)

	c.JSON(http.StatusOK,
		gin.H{
			"node_key":     host.NodeKey,
			"node_invalid": false,
		})
}

func OsqueryConfig(c *gin.Context) {
	var body OsqueryConfigPostBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Debugf("Error parsing OsqueryConfig POST body: %s", err.Error())
		return
	}
	logrus.Debugf("OsqueryConfig: %s", body.NodeKey)

	c.JSON(http.StatusOK,
		gin.H{
			"schedule": map[string]map[string]interface{}{
				"time": {
					"query":    "select * from time;",
					"interval": 1,
				},
			},
			"node_invalid": false,
		})
}

func OsqueryLog(c *gin.Context) {
	var body OsqueryLogPostBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Debugf("Error parsing OsqueryLog POST body: %s", err.Error())
		return
	}
	logrus.Debugf("OsqueryLog: %s", body.LogType)

	if body.LogType == "status" {
		for _, data := range body.Data {
			var log OsqueryStatusLog

			severity, ok := data["severity"].(string)
			if ok {
				log.Severity = severity
			} else {
				logrus.Error("Error asserting the type of status log severity")
			}

			filename, ok := data["filename"].(string)
			if ok {
				log.Filename = filename
			} else {
				logrus.Error("Error asserting the type of status log filename")
			}

			line, ok := data["line"].(string)
			if ok {
				log.Line = line
			} else {
				logrus.Error("Error asserting the type of status log line")
			}

			message, ok := data["message"].(string)
			if ok {
				log.Message = message
			} else {
				logrus.Error("Error asserting the type of status log message")
			}

			version, ok := data["version"].(string)
			if ok {
				log.Version = version
			} else {
				logrus.Error("Error asserting the type of status log version")
			}

			logrus.WithFields(logrus.Fields{
				"node_key": body.NodeKey,
				"severity": log.Severity,
				"filename": log.Filename,
				"line":     log.Line,
				"version":  log.Version,
			}).Info(log.Message)
		}
	} else if body.LogType == "result" {
		// TODO: handle all of the different kinds of results logs
	}

	c.JSON(http.StatusOK,
		gin.H{
			"node_invalid": false,
		})
}

func OsqueryDistributedRead(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"node_invalid": true,
		})
	return
	logrus.Debug("Distributed read")
	var body OsqueryDistributedReadPostBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Debugf("Error parsing OsqueryDistributedRead POST body: %s", err.Error())
		return
	}
	logrus.Debugf("OsqueryDistributedRead: %s", body.NodeKey)

	c.JSON(http.StatusOK,
		gin.H{
			"queries": map[string]string{
				"id1": "select * from osquery_info",
			},
			"node_invalid": false,
		})
}

func OsqueryDistributedWrite(c *gin.Context) {
	var body OsqueryDistributedWritePostBody
	err := c.BindJSON(&body)
	if err != nil {
		logrus.Debugf("Error parsing OsqueryDistributedWrite POST body: %s", err.Error())
		return
	}
	logrus.Debugf("OsqueryDistributedWrite: %s", body.NodeKey)
	c.JSON(http.StatusOK,
		gin.H{
			"node_invalid": false,
		})
}
