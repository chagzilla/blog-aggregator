package main

import (
	"internal/database"
	"net/http"
	"context"
	"io"
	"encoding/xml"
	"fmt"
	"html"
	"errors"
	"time"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	req.Header.Add("User-Agent", "gator")

	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed 
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return &feed, nil
}

func scrapeFeeds(s * state) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
}

func handlerFeed(s * state, cmd command) error {
	if len (cmd.arguments) > 0 {
		return errors.New("Too many arguements: feed takes 0 arguments")
	}
	f, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", f)
	return scrapeFeeds(s)
}

func handlerAddFeed(s * state, cmd command, u database.User) error {
	if len (cmd.arguments) == 0 {
		return errors.New("Not enough arguments: Add feed takes 1 arguments, url")
	}

	feed := database.CreateFeedParams { Name: cmd.arguments[0], Url: cmd.arguments[1], UserID: u.ID }
	f, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return err
	}
	fmt.Printf("%s, %s, %s\n", f.Name, f.Url, f.UserID)

	feed_follows := database.CreateFeedFollowParams { ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: u.ID, FeedID: f.Url}
	_, err = s.db.CreateFeedFollow(context.Background(), feed_follows)
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s * state, cmd command) error {
	if len (cmd.arguments) > 0 {
		return errors.New("Too many arguements: feed takes 0 arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, val := range feeds {
		fmt.Printf("%s, %s, %s\n", val.Name, val.Url, val.Name_2)
	}
	return nil
}

func handlerFollowFeed(s * state, cmd command, u database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("Wrong number of arguments: follow takes 1 arguments, url")
	}

	feed_follows := database.CreateFeedFollowParams { ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: u.ID, FeedID: cmd.arguments[0]}
	ff, err := s.db.CreateFeedFollow(context.Background(), feed_follows)
	if err != nil {
		return err
	}

	fmt.Printf("%s, %s, %s, %s, %s", ff.ID, ff.CreatedAt, ff.UpdatedAt, ff.Name, ff.Name_2)
	return nil
}

func handlerFollowing(s * state, cmd command, u database.User) error {
	if len (cmd.arguments) > 0 {
		return errors.New("Too many arguements: following takes 0 arguments")
	}
	
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), u.ID)
	if err != nil {
		return err
	}
	for _, ff := range feed_follows {
		fmt.Printf("%s, %s, %s, %s, %s\n", ff.ID, ff.CreatedAt, ff.UpdatedAt, ff.Name, ff.Name_2)
	}
	return nil
}

func handleUnfollow(s * state, cmd command, u database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("Wrong number of arguments: unfollow takes 1 arguments")
	}

	feed_follow, err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams {
		UserID: u.ID,
		FeedID: cmd.arguments[0],
	})
	if err != nil {
		return err
	}
	
	fmt.Printf("Delete %s for %s", feed_follow.FeedID, u.Name)
	return nil
}
