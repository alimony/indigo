package schemagen

import (
	"context"

	"github.com/whyrusleeping/gosky/xrpc"
)

// schema: app.bsky.feed.getTimeline

func init() {
}

type FeedGetTimeline_MyState struct {
	Repost   string `json:"repost" cborgen:"repost"`
	Upvote   string `json:"upvote" cborgen:"upvote"`
	Downvote string `json:"downvote" cborgen:"downvote"`
}

type FeedGetTimeline_Output struct {
	Cursor *string                     `json:"cursor" cborgen:"cursor"`
	Feed   []*FeedGetTimeline_FeedItem `json:"feed" cborgen:"feed"`
}

type FeedGetTimeline_FeedItem struct {
	ReplyCount    int64                    `json:"replyCount" cborgen:"replyCount"`
	TrendedBy     *ActorRef_WithInfo       `json:"trendedBy" cborgen:"trendedBy"`
	RepostedBy    *ActorRef_WithInfo       `json:"repostedBy" cborgen:"repostedBy"`
	Embed         *FeedEmbed               `json:"embed" cborgen:"embed"`
	Record        any                      `json:"record" cborgen:"record"`
	RepostCount   int64                    `json:"repostCount" cborgen:"repostCount"`
	UpvoteCount   int64                    `json:"upvoteCount" cborgen:"upvoteCount"`
	DownvoteCount int64                    `json:"downvoteCount" cborgen:"downvoteCount"`
	IndexedAt     string                   `json:"indexedAt" cborgen:"indexedAt"`
	Uri           string                   `json:"uri" cborgen:"uri"`
	Cid           string                   `json:"cid" cborgen:"cid"`
	Author        *ActorRef_WithInfo       `json:"author" cborgen:"author"`
	MyState       *FeedGetTimeline_MyState `json:"myState" cborgen:"myState"`
}

func FeedGetTimeline(ctx context.Context, c *xrpc.Client, algorithm string, before string, limit int64) (*FeedGetTimeline_Output, error) {
	var out FeedGetTimeline_Output

	params := map[string]interface{}{
		"algorithm": algorithm,
		"before":    before,
		"limit":     limit,
	}
	if err := c.Do(ctx, xrpc.Query, "", "app.bsky.feed.getTimeline", params, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}