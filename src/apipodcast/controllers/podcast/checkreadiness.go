package podcast

import (
	"net/http"

	// _ "github.com/denisenkom/go-mssqldb"

	// _ "github.com/jinzhu/gorm/dialects/mysql"

	// Mapp "apibafgate/model"

	// Morderdev "apibafgate/model"

	// Conf "apibafgate/Config"

	Conf "apigopodcast/src/config"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

// var dbbaflite, errdb = Conf.Connectbaflite()

func CheckReadiness(c *gin.Context) {
	addr, err := Conf.DetermineListenAddressPodcast()
	if err != nil {
		logger.Fatal(err)
		logger.Println(err)
	}
	host, errhost := Conf.Hostname()
	if errhost != nil {
		logger.Fatal(errhost)
		logger.Println(errhost)
	}
	client := &http.Client{}
	req, errreq := http.NewRequest("GET", host+addr+"/api/dev/v1/podcast/swagger/index.html", nil)
	// req, errreq := http.NewRequest("GET", "http://172.16.1.187:2020/api/dev/v1/apibafcommon/swagger/index.html", nil)
	if errreq != nil {
		// log.Fatal(errreq)
		logger.WithFields(logger.Fields{
			"status": "500",
		}).Error("Can't Get", errreq)
	}
	req.Header.Set("X-Health-Check", "1")
	resp, errresp := client.Do(req)
	if errresp != nil {
		// log.Fatal(errresp)
		c.JSON(http.StatusOK, gin.H{"response": http.StatusInternalServerError, "status": "Service is not ready"})
		logger.WithFields(logger.Fields{
			"status": "500",
		}).Error("Service is not ready because ", errresp)
	}
	// bodyText, errbodytext := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	// log.Fatal(errbodytext)
	// 	logger.Info(errbodytext)
	// }
	c.JSON(http.StatusOK, gin.H{"response": http.StatusOK, "status": "Service is ready"})
	logger.WithFields(logger.Fields{
		"status": resp.StatusCode,
	}).Info("Service is ready")
}
