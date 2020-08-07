package podcast

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexcesaro/log/stdlog"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// _ "github.com/jinzhu/gorm/dialects/mysql"

	// Mapp "apibafgate/model"

	// Morderdev "apibafgate/model"

	Mpodcast "apigopodcast/src/apipodcast/model/podcast"

	// Conf "apibafgate/Config"
	Conf "apigopodcast/src/config"

	logger "github.com/sirupsen/logrus"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func GetPodcast(c *gin.Context) {
	logkoe := stdlog.GetFromFlags()

	env, errenv := Conf.Environment()
	if errenv != nil {
		logger.Println(errenv)
		logkoe.Info(errenv)
	} else {
		if env == "production" {
			gin.SetMode(gin.ReleaseMode)
			// router := gin.New()
		} else if env == "development" {
			gin.SetMode(gin.DebugMode)
		}
	}

	var ctx = func() context.Context {
		return context.Background()
	}()
	quickstartDatabase, errdb := Conf.Connectmongo()
	// quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	// episodesCollection := quickstartDatabase.Collection("episodes")
	//fmt.Println(errdb)

	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{"statusload": http.StatusInternalServerError, "statusdb": errdb,
			"result": "Missing Connection"})
		logger.WithFields(logger.Fields{
			"detail": errdb,
		}).Error("Missing Connection")
		logkoe.Info("Missing Connection", "statusdb:", errdb, "statusload :", http.StatusInternalServerError)
		return
	}

	var txt Mpodcast.FindPodcastJSON
	c.BindJSON(&txt)
	podcastID := txt.ID

	if podcastID == "" {
		var podcasts []Mpodcast.Podcast
		podcastCursor, err := podcastsCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		//fmt.Println(podcasts)
		if err = podcastCursor.All(ctx, &podcasts); err != nil {
			panic(err)
		}
		defer podcastCursor.Close(ctx)
		if err != nil {
			c.JSON(404, gin.H{"message": "Podcast Not Found"})
		} else {
			c.JSON(200, gin.H{"data": podcasts})
		}
	} else {
		id, _ := primitive.ObjectIDFromHex(podcastID)
		var podcasts []Mpodcast.Podcast

		podcastCursor, err := podcastsCollection.Find(ctx, bson.M{"_id": id})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if err = podcastCursor.All(ctx, &podcasts); err != nil {
			panic(err)
		}
		fmt.Println(podcasts)
		// c.JSON(podcasts)
		defer podcastCursor.Close(ctx)
		if podcasts == nil {
			c.JSON(404, gin.H{"message": "Podcast Not Found"})
		} else {
			c.JSON(200, gin.H{"data": podcasts})
		}
	}
}

func InsertPodcast(c *gin.Context) {
	logkoe := stdlog.GetFromFlags()

	env, errenv := Conf.Environment()
	if errenv != nil {
		logger.Println(errenv)
		logkoe.Info(errenv)
	} else {
		if env == "production" {
			gin.SetMode(gin.ReleaseMode)
			// router := gin.New()
		} else if env == "development" {
			gin.SetMode(gin.DebugMode)
		}
	}

	var ctx = func() context.Context {
		return context.Background()
	}()
	quickstartDatabase, errdb := Conf.Connectmongo()
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{"statusload": http.StatusInternalServerError, "statusdb": errdb,
			"result": "Missing Connection"})
		logger.WithFields(logger.Fields{
			"detail": errdb,
		}).Error("Missing Connection")
		logkoe.Info("Missing Connection", "statusdb:", errdb, "statusload :", http.StatusInternalServerError)
		return
	}

	var txtEpisodePodcast Mpodcast.InsertEpisodePodcastJSON
	c.BindJSON(&txtEpisodePodcast)

	var Title = txtEpisodePodcast.Title
	var Author = txtEpisodePodcast.Author
	var Tags = txtEpisodePodcast.Tags
	// getTimeZone, _ := time.LoadLocation("Asia/Jakarta")
	// currentTimeCrt := time.Now().In(getTimeZone)
	// currentTimeUpd := time.Now().In(getTimeZone)
	//getTime := time.Now()
	var currentTimeCrt = time.Now()
	var currentTimeUpd = time.Now()
	//fmt.Println(currentTimeCrt)
	//os.Exit(1)
	//fmt.Println(getTime)

	if Title == "" || Author == "" {
		c.JSON(404, gin.H{"Message": "Something went wrong, please check your data podcast!"})
	} else {
		var episode = txtEpisodePodcast.Episode
		var description = txtEpisodePodcast.Description
		var duration = txtEpisodePodcast.Duration
		if episode == "" || description == "" || duration == 0 {
			c.JSON(404, gin.H{"Message": "Something went wrong, please check your data episode!"})
		} else {

			podcastsResult, err := podcastsCollection.InsertOne(ctx, bson.D{
				{Key: "title", Value: Title},
				{Key: "author", Value: Author},
				{Key: "tags", Value: Tags},
				{Key: "dtmcrt", Value: currentTimeCrt},
				{Key: "dtmupd", Value: currentTimeUpd},
			})
			if err != nil {
				log.Fatal(err)
			}
			c.JSON(200, gin.H{"podcast": podcastsResult.InsertedID, "Message": "Your Podcast has been inserted"})

			//fmt.Println(podcastsResult)
			//fmt.Println(podcastsResult.InsertedID)

			var podcast = podcastsResult.InsertedID
			var episode = txtEpisodePodcast.Episode
			var description = txtEpisodePodcast.Description
			var duration = txtEpisodePodcast.Duration
			episodesResult, err := episodesCollection.InsertMany(ctx, []interface{}{
				bson.D{
					{Key: "podcast", Value: podcast},
					{Key: "description", Value: description},
					{Key: "duration", Value: duration},
					{Key: "episode ", Value: episode},
					{Key: "dtmcrt", Value: currentTimeCrt},
					{Key: "dtmupd", Value: currentTimeUpd},
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			// episodeResultString := fmt.Sprintf("%v", episodesResult.InsertedIDs)

			c.JSON(200, gin.H{"episode": episodesResult.InsertedIDs, "Message": "You Episode Has Been Inserted"})
		}
		//fmt.Println(episodesResult.InsertedIDs)
	}
}

func DeletePodcast(c *gin.Context) {
	logkoe := stdlog.GetFromFlags()

	env, errenv := Conf.Environment()
	if errenv != nil {
		logger.Println(errenv)
		logkoe.Info(errenv)
	} else {
		if env == "production" {
			gin.SetMode(gin.ReleaseMode)
			// router := gin.New()
		} else if env == "development" {
			gin.SetMode(gin.DebugMode)
		}
	}

	var ctx = func() context.Context {
		return context.Background()
	}()
	quickstartDatabase, errdb := Conf.Connectmongo()
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	// episodesCollection := quickstartDatabase.Collection("episodes")

	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{"statusload": http.StatusInternalServerError, "statusdb": errdb,
			"result": "Missing Connection"})
		logger.WithFields(logger.Fields{
			"detail": errdb,
		}).Error("Missing Connection")
		logkoe.Info("Missing Connection", "statusdb:", errdb, "statusload :", http.StatusInternalServerError)
		return
	}

	var txtDeletePodcast Mpodcast.DeletePodcastJSON
	c.BindJSON(&txtDeletePodcast)

	podcastID := txtDeletePodcast.ID
	id, _ := primitive.ObjectIDFromHex(podcastID)

	result, err := podcastsCollection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"result": result.DeletedCount, "status": "Failed"})
	} else {
		c.JSON(200, gin.H{"update": "Update %v Documents!\n", "result": result.DeletedCount, "status": "Success"})
	}
	//fmt.Printf("Update %v Documents!\n", result.DeletedCount)
}

func UpdatePodcast(c *gin.Context) {
	logkoe := stdlog.GetFromFlags()

	env, errenv := Conf.Environment()
	if errenv != nil {
		logger.Println(errenv)
		logkoe.Info(errenv)
	} else {
		if env == "production" {
			gin.SetMode(gin.ReleaseMode)
			// router := gin.New()
		} else if env == "development" {
			gin.SetMode(gin.DebugMode)
		}
	}

	var ctx = func() context.Context {
		return context.Background()
	}()
	quickstartDatabase, errdb := Conf.Connectmongo()
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	// episodesCollection := quickstartDatabase.Collection("episodes")

	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{"statusload": http.StatusInternalServerError, "statusdb": errdb,
			"result": "Missing Connection"})
		logger.WithFields(logger.Fields{
			"detail": errdb,
		}).Error("Missing Connection")
		logkoe.Info("Missing Connection", "statusdb:", errdb, "statusload :", http.StatusInternalServerError)
		return
	}

	var txtUpdatePodcast Mpodcast.UpdatePodcastJSON
	c.BindJSON(&txtUpdatePodcast)

	podcastID := txtUpdatePodcast.ID
	title := txtUpdatePodcast.Title
	author := txtUpdatePodcast.Author
	tags := txtUpdatePodcast.Tags
	currentTimeUpd := time.Now()

	id, _ := primitive.ObjectIDFromHex(podcastID)
	result, err := podcastsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{
				"$set", bson.D{
					{Key: "title", Value: title},
					{Key: "author", Value: author},
					{Key: "tags", Value: tags},
					{Key: "dtmupd", Value: currentTimeUpd},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 0 {
		c.JSON(404, gin.H{"result": result.ModifiedCount, "status": "Failed", "Message": "No One Field modified"})
	} else {
		c.JSON(200, gin.H{"update": "Update %v Documents!\n", "result": result.ModifiedCount, "status": "Success"})
	}

}
