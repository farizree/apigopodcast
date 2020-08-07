package Mpodcast

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISODate struct {
	time.Time
}

func (t *ISODate) String() string {
	return t.Format("2006-01-02")
}

type (
	Podcast struct {
		ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Title  string             `json:"title,omitempty" bson:"title,omitempty"`
		Author string             `json:"author,omitempty" bson:"author,omitempty"`
		Tags   []string           `json:"tags,omitempty" bson:"tags,omitempty"`
		Dtmcrt time.Time          `bson:"dtmcrt"`
		Dtmupd time.Time          `bson:"dtmupd"`
	}
	// PodcastJSON struct {
	// 	ID     primitive.ObjectID `json:"_id"`
	// 	Title  string             `json:"title"`
	// 	Author string             `json:"author"`
	// 	Tags   []string           `json:"tags"`
	// 	Dtmcrt string             `json:"dtmcrt"`
	// 	Dtmupd string             `json:"dtmupd"`
	// }

	FindPodcastJSON struct {
		ID string `json:"id"`
	}
	InsertEpisodePodcastJSON struct {
		Title       string   `json:"title"`
		Author      string   `json:"author"`
		Tags        []string `json:"tags"`
		Episode     string   `json:"episode"`
		Description string   `json:"description"`
		Duration    int32    `json:"duration"`
	}

	Episode struct {
		ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Podcast     primitive.ObjectID `json:"podcast,omitempty" bson:"podcast,omitempty"`
		Episode     string             `json:"episode,omitempty" bson:"episode,omitempty"`
		Description string             `json:"description,omitempty" bson:"description,omitempty"`
		Duration    int32              `json:"duration,omitempty" bson:"duration,omitempty"`
	}

	DeletePodcastJSON struct {
		ID string `json:"id"`
	}

	UpdatePodcastJSON struct {
		ID     string   `json:"id" bson:"id"`
		Title  string   `json:"title"`
		Author string   `json:"author"`
		Tags   []string `json:"tags"`
	}
)
