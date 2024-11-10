package posts

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
)

type FacetType int

const (
	LinkType FacetType = iota + 1
	MentionType
	TagType
)

func (t FacetType) String() string {
	base := "app.bsky.richtext.facet#"
	switch t {
	case LinkType:
		return fmt.Sprintf("%slink", base)
	case MentionType:
		return fmt.Sprintf("%smention", base)
	case TagType:
		return fmt.Sprintf("%stag", base)
	default:
		return "Unknown"
	}
}

type Facet struct {
	Val       string
	T_facet   string
	FacetType FacetType
}

type Embed struct {
	Link           Link
	Images         []Image
	UploadedImages []lexutil.LexBlob
}

type Link struct {
	Title       string
	Url         url.URL
	Description string
	Thumbnail   lexutil.LexBlob
}

type Image struct {
	Title string
	Url   url.URL
}

type PostBuilder struct {
	Text  string
	Facet []Facet
	Embed Embed
}

func NewPostBuilder(msg string) PostBuilder {
	return PostBuilder{
		Text:  msg,
		Facet: []Facet{},
	}
}

func (p PostBuilder) WithLink(title, description string, link url.URL, thumb lexutil.LexBlob) PostBuilder {
	p.Embed.Link.Title = title
	p.Embed.Link.Description = description
	p.Embed.Link.Url = link
	p.Embed.Link.Thumbnail = thumb

	return p
}

func (p PostBuilder) WithImages(blobs []lexutil.LexBlob, images []Image) PostBuilder {
	p.Embed.Images = images
	p.Embed.UploadedImages = blobs

	return p
}

var FeedPost_Embed bsky.FeedPost_Embed

func (p PostBuilder) CreatePost() (bsky.FeedPost, error) {
	post := bsky.FeedPost{}

	if p.Text == "" {
		return bsky.FeedPost{}, errors.New("post cannot be blank")
	}

	post.Text = p.Text
	post.LexiconTypeID = "app.bsky.feed.post"
	post.CreatedAt = time.Now().Format(time.RFC3339)

	Facets := []*bsky.RichtextFacet{}

	for _, f := range p.Facet {
		facet := &bsky.RichtextFacet{}
		features := []*bsky.RichtextFacet_Features_Elem{}
		feature := &bsky.RichtextFacet_Features_Elem{}

		switch f.FacetType {

		case LinkType:
			{
				feature = &bsky.RichtextFacet_Features_Elem{
					RichtextFacet_Link: &bsky.RichtextFacet_Link{
						LexiconTypeID: f.FacetType.String(),
						Uri:           f.Val,
					},
				}
			}

		case MentionType:
			{
				feature = &bsky.RichtextFacet_Features_Elem{
					RichtextFacet_Mention: &bsky.RichtextFacet_Mention{
						LexiconTypeID: f.FacetType.String(),
						Did:           f.Val,
					},
				}
			}

		case TagType:
			{
				feature = &bsky.RichtextFacet_Features_Elem{
					RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
						LexiconTypeID: f.FacetType.String(),
						Tag:           f.Val,
					},
				}
			}

		}

		features = append(features, feature)
		facet.Features = features

		ByteStart, ByteEnd, err := findSubstring(post.Text, f.T_facet)
		if err != nil {
			return post, fmt.Errorf("unable to find the substring: %v , %v", f.T_facet, err)
		}

		index := &bsky.RichtextFacet_ByteSlice{
			ByteStart: int64(ByteStart),
			ByteEnd:   int64(ByteEnd),
		}
		facet.Index = index

		Facets = append(Facets, facet)
	}

	post.Facets = Facets

	if p.Embed.Link != (Link{}) {
		FeedPost_Embed.EmbedExternal = &bsky.EmbedExternal{
			LexiconTypeID: "app.bsky.embed.external",
			External: &bsky.EmbedExternal_External{
				Title:       p.Embed.Link.Title,
				Uri:         p.Embed.Link.Url.String(),
				Description: p.Embed.Link.Description,
				Thumb:       &p.Embed.Link.Thumbnail,
			},
		}
	} else {
		if len(p.Embed.Images) != 0 && len(p.Embed.Images) == len(p.Embed.UploadedImages) {

			EmbedImages := bsky.EmbedImages{
				LexiconTypeID: "app.bsky.embed.images",
				Images:        make([]*bsky.EmbedImages_Image, len(p.Embed.Images)),
			}

			for i, img := range p.Embed.Images {
				EmbedImages.Images[i] = &bsky.EmbedImages_Image{
					Alt:   img.Title,
					Image: &p.Embed.UploadedImages[i],
				}
			}

			FeedPost_Embed.EmbedImages = &EmbedImages

		}
	}

	// avoid error when trying to marshal empty field (*bsky.FeedPost_Embed)
	if len(p.Embed.Images) != 0 || p.Embed.Link.Title != "" {
		post.Embed = &FeedPost_Embed
	}

	return post, nil
}
