package agolclient

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

func UsersCsv(w io.Writer, users []User) {
	cw := csv.NewWriter(w)

	cw.Write([]string{"Username", "Full Name", "Email", "Role"})
	for _, u := range users {
		cw.Write([]string{u.Username, u.FullName, u.Email, u.Role})
	}

	cw.Flush()
}

type MyAGOL struct {
	User         *User
	Org          *Org
	Subscription *Subscription
}

type User struct {
	Username, FullName, Email, Role, Thumbnail string
	Groups                                     []Group
}

func (u *User) RelativeThumbnailUrl() string {
	if u.Thumbnail == "" {
		return ""
	}
	return fmt.Sprintf("/community/users/%s/info/%s", u.Username, u.Thumbnail)
}

func (u *User) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", u.Username, u.FullName, u.Email, u.Role)
}

type GroupMembership struct {
	Username, MemberType string
	Applications         int
}

type Group struct {
	Id, Title, Owner, Access string
	UserMembership           *GroupMembership
}

type Subscription struct {
	Id, Type, State   string
	ExpDate, MaxUsers int64
	AvailableCredits  float64
}

func (s *Subscription) Expires() *time.Time {
	if s.ExpDate == 0 {
		return nil
	}
	t := time.Unix(0, s.ExpDate*int64(time.Millisecond))
	return &t
}

type Org struct {
	Id, Name, UrlKey, Thumbnail string
	AllSSL                      bool
}

func (org *Org) RelativeThumbnailUrl() string {
	if org.Thumbnail == "" {
		return ""
	}
	return fmt.Sprintf("/portals/%s/resources/%s", org.Id, org.Thumbnail)
}

type PortalSelf struct {
	*Org
	SubscriptionInfo *Subscription
	User             *User
	FeaturedGroups   []Group
}

type UsersResponse struct {
	Total     int
	Start     int
	Num       int
	NextStart int
	Users     []User
}

func (ur *UsersResponse) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Total: %d, Start: %d, Num %d, Next Start: %d", ur.Total, ur.Start, ur.Num, ur.NextStart)

	for _, u := range ur.Users {
		fmt.Fprintf(&buf, "\n%s", u)
	}

	return buf.String()
}
