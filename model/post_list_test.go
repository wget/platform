// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostListJson(t *testing.T) {

	pl := PostList{}
	p1 := &Post{Id: NewId(), Message: NewId()}
	pl.AddPost(p1)
	p2 := &Post{Id: NewId(), Message: NewId()}
	pl.AddPost(p2)

	pl.AddOrder(p1.Id)
	pl.AddOrder(p2.Id)

	json := pl.ToJson()
	rpl := PostListFromJson(strings.NewReader(json))

	if rpl.Posts[p1.Id].Message != p1.Message {
		t.Fatal("failed to serialize")
	}

	if rpl.Posts[p2.Id].Message != p2.Message {
		t.Fatal("failed to serialize")
	}

	if rpl.Order[1] != p2.Id {
		t.Fatal("failed to serialize")
	}
}

func TestPostListExtend(t *testing.T) {
	p1 := &Post{Id: NewId(), Message: NewId()}
	p2 := &Post{Id: NewId(), Message: NewId()}
	p3 := &Post{Id: NewId(), Message: NewId()}

	l1 := PostList{}
	l1.AddPost(p1)
	l1.AddOrder(p1.Id)
	l1.AddPost(p2)
	l1.AddOrder(p2.Id)

	l2 := PostList{}
	l2.AddPost(p3)
	l2.AddOrder(p3.Id)

	l2.Extend(&l1)

	// should not changed l1
	assert.Len(t, l1.Posts, 2)
	assert.Len(t, l1.Order, 2)

	// should extend l2
	assert.Len(t, l2.Posts, 3)
	assert.Len(t, l2.Order, 3)

	// should extend the Order of l2 correctly
	assert.Equal(t, l2.Order[0], p3.Id)
	assert.Equal(t, l2.Order[1], p1.Id)
	assert.Equal(t, l2.Order[2], p2.Id)

	// extend l2 again
	l2.Extend(&l1)
	// extending l2 again should not changed l1
	assert.Len(t, l1.Posts, 2)
	assert.Len(t, l1.Order, 2)

	// extending l2 again should extend l2
	assert.Len(t, l2.Posts, 3)
	assert.Len(t, l2.Order, 3)

	// p3 could be last unread
	p4 := &Post{Id: NewId(), Message: NewId()}
	p5 := &Post{Id: NewId(), RootId: p1.Id, Message: NewId()}
	p6 := &Post{Id: NewId(), RootId: p1.Id, Message: NewId()}

	// Create before and after post list where p3 could be last unread

	// Order has 2 but Posts are 4 which includes additional 2 comments (p5 & p6) to parent post (p1)
	beforePostList := PostList{
		Order: []string{p1.Id, p2.Id},
		Posts: map[string]*Post{p1.Id: p1, p2.Id: p2, p5.Id: p5, p6.Id: p6},
	}

	// Order has 3 but Posts are 4 which includes 1 parent post (p1) of comments (p5 & p6)
	afterPostList := PostList{
		Order: []string{p4.Id, p5.Id, p6.Id},
		Posts: map[string]*Post{p1.Id: p1, p4.Id: p4, p5.Id: p5, p6.Id: p6},
	}

	beforePostList.Extend(&afterPostList)

	// should not changed afterPostList
	assert.Len(t, afterPostList.Posts, 4)
	assert.Len(t, afterPostList.Order, 3)

	// should extend beforePostList
	assert.Len(t, beforePostList.Posts, 5)
	assert.Len(t, beforePostList.Order, 5)

	// should extend the Order of beforePostList correctly
	assert.Equal(t, beforePostList.Order[0], p1.Id)
	assert.Equal(t, beforePostList.Order[1], p2.Id)
	assert.Equal(t, beforePostList.Order[2], p4.Id)
	assert.Equal(t, beforePostList.Order[3], p5.Id)
	assert.Equal(t, beforePostList.Order[4], p6.Id)
}

func TestPostListSortByCreateAt(t *testing.T) {
	pl := PostList{}
	p1 := &Post{Id: NewId(), Message: NewId(), CreateAt: 2}
	pl.AddPost(p1)
	p2 := &Post{Id: NewId(), Message: NewId(), CreateAt: 1}
	pl.AddPost(p2)
	p3 := &Post{Id: NewId(), Message: NewId(), CreateAt: 3}
	pl.AddPost(p3)

	pl.AddOrder(p1.Id)
	pl.AddOrder(p2.Id)
	pl.AddOrder(p3.Id)

	pl.SortByCreateAt()

	assert.EqualValues(t, pl.Order[0], p3.Id)
	assert.EqualValues(t, pl.Order[1], p1.Id)
	assert.EqualValues(t, pl.Order[2], p2.Id)
}

func TestPostListToSlice(t *testing.T) {
	pl := PostList{}
	p1 := &Post{Id: NewId(), Message: NewId(), CreateAt: 2}
	pl.AddPost(p1)
	p2 := &Post{Id: NewId(), Message: NewId(), CreateAt: 1}
	pl.AddPost(p2)
	p3 := &Post{Id: NewId(), Message: NewId(), CreateAt: 3}
	pl.AddPost(p3)

	pl.AddOrder(p1.Id)
	pl.AddOrder(p2.Id)
	pl.AddOrder(p3.Id)

	want := []*Post{p1, p2, p3}

	assert.Equal(t, want, pl.ToSlice())
}
