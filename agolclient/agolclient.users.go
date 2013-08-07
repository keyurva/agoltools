package agolclient

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

type User struct {
	Username, FullName, Email, Role string
}

func (u *User) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", u.Username, u.FullName, u.Email, u.Role)
}

func UsersCsv(w io.Writer, users []User) {
	cw := csv.NewWriter(w)

	cw.Write([]string{"Username", "Full Name", "Email", "Role"})
	for _, u := range users {
		cw.Write([]string{u.Username, u.FullName, u.Email, u.Role})
	}

	cw.Flush()
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
