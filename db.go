package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func newDatabase(connectionString string) *sql.DB {
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	return database
}

type registrationPayload struct {
	email    string `json:"email"`
	password string `json:"password"`
	fullname string `json:"fullname"`
	role     string `json:"role"`
	enabled  bool   `json:"enabled"`
}

type complaintPayload struct {
	email    string `json:"email"`
	fullname string `json:"fullname"`
	comment  string `json:"role"`
	file     string `json:"role"`
}

type complaint struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	Comment   string `json:"role"`
	File      string `json:"file"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

type member struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Enabled  bool   `json:"enabled"`
}

type loginPayload struct {
	email    string `json:"email"`
	password string `json:"password"`
}

func dbAddUser(payload registrationPayload) bool {
	hashedPassword, err := hashPassword(payload.password)
	if err != nil {
		log.Printf("[DB:AddUser] hashing problem %s", err)
		return false
	}
	_, err = db.Query("insert into users (fullname,email,password,role,enabled) values ($1,$2,$3,$4,$5)", payload.fullname, payload.email, hashedPassword, payload.role, true)
	if err != nil {
		log.Printf("[DB:AddUser] insert problem %s", err)
		return false
	}
	return true
}
func dbCreateComplaint(payload complaintPayload) bool {
	_, err := db.Query("insert into complaints (fullname,email,comment,file) values ($1,$2,$3,$4)", payload.fullname, payload.email, payload.comment, payload.file)
	if err != nil {
		log.Printf("[DB:CreateComplaint] insert problem %s", err)
		return false
	}
	return true
}

func dbUpdateComplaint(id string, payload struct {
	comment string
	file    string
}) bool {
	_, err := db.Query("update complaints set comment = $1, file =$2 where id = $3 ", payload.comment, payload.file, id)
	if err != nil {
		log.Printf("[DB:UpdateComplaint] %s", err)
		return false
	}
	return true
}

func dbUpdateComplaintComment(id string, payload struct {
	comment string
}) bool {
	_, err := db.Query("update complaints set comment = $1 where id = $2 ", payload.comment, id)
	if err != nil {
		log.Printf("[DB:UpdateComplaintComment] %s", err)
		return false
	}
	return true
}

func dbGetUser(email string) *registrationPayload {
	result := db.QueryRow("select password,email,role,fullname,enabled from users where email=$1", email)
	user := registrationPayload{}
	err := result.Scan(&user.password, &user.email, &user.role, &user.fullname, &user.enabled)
	if err != nil {
		log.Printf("[DB:GetUser]  %s", err)
		if err == sql.ErrNoRows {
			return nil
		}

		return nil
	}
	return &user
}

func dbUserOwnsFile(email string, file string) bool {
	result := db.QueryRow("select file from complaints where email=$1 and file=$2", email, file)
	var filename string
	err := result.Scan(&filename)
	if err != nil {
		log.Printf("[DB:UserOwnsFile]  %s", err)
		if err == sql.ErrNoRows {
			return false
		}

		return false
	}
	return true
}

func dbComplaintByID(id string) *complaint {
	result := db.QueryRow("select id,email,fullname,comment,file, updated_at, created_at from complaints where id=$1", id)
	cmp := complaint{}
	err := result.Scan(&cmp.ID, &cmp.Email, &cmp.Fullname,
		&cmp.Comment, &cmp.File, &cmp.UpdatedAt, &cmp.CreatedAt)
	if err != nil {
		log.Printf("[DB:ComplaintByID] %s", err)
		if err == sql.ErrNoRows {
			return nil
		}

		return nil
	}
	return &cmp
}

func dbComplaintsByEmail(email string) ([]complaint, error) {
	rows, err := db.Query("select id,email,fullname,comment,file, updated_at, created_at from complaints where email=$1", email)
	if err != nil {
		log.Printf("[DB:ComplaintsByEmail] %s", err)
		return nil, err
	}
	defer rows.Close()

	var complaints []complaint

	for rows.Next() {
		var cmp complaint
		if err := rows.Scan(&cmp.ID, &cmp.Email, &cmp.Fullname,
			&cmp.Comment, &cmp.File, &cmp.UpdatedAt, &cmp.CreatedAt); err != nil {
			log.Printf("[DB:ComplaintsByEmail] %s", err)
			return complaints, err
		}
		complaints = append(complaints, cmp)
	}
	if err = rows.Err(); err != nil {
		log.Printf("[DB:ComplaintsByEmail] %s", err)
		return complaints, err
	}
	return complaints, nil
}

func dbComplaintsAll() ([]complaint, error) {
	rows, err := db.Query("select id,email,fullname,comment,file, updated_at, created_at from complaints")
	if err != nil {
		log.Printf("[DB:ComplaintsAll] %s", err)
		return nil, err
	}
	defer rows.Close()

	var complaints []complaint

	for rows.Next() {
		var cmp complaint
		if err := rows.Scan(&cmp.ID, &cmp.Email, &cmp.Fullname,
			&cmp.Comment, &cmp.File, &cmp.UpdatedAt, &cmp.CreatedAt); err != nil {
			log.Printf("[DB:ComplaintsAll] %s", err)
			return complaints, err
		}
		complaints = append(complaints, cmp)
	}
	if err = rows.Err(); err != nil {
		log.Printf("[DB:ComplaintsAll] %s", err)
		return complaints, err
	}
	return complaints, nil
}

func dbMembersAll() ([]member, error) {
	rows, err := db.Query("select email,fullname,enabled from users where role='member'")
	if err != nil {
		log.Printf("[DB:MemberssAll] %s", err)
		return nil, err
	}
	defer rows.Close()

	var members []member

	for rows.Next() {
		var mem member
		if err := rows.Scan(&mem.Email, &mem.Fullname,
			&mem.Enabled); err != nil {
			log.Printf("[DB:MembersAll] %s", err)
			return members, err
		}
		members = append(members, mem)
	}
	if err = rows.Err(); err != nil {
		log.Printf("[DB:MembersAll] %s", err)
		return members, err
	}
	return members, nil
}

func dbUpdateMemberAccount(email string, enabled bool) bool {
	_, err := db.Query("update users set enabled = $1 where email = $2 ", enabled, email)
	if err != nil {
		log.Printf("[DB:UpdateMemberAccount] %s", err)
		return false
	}
	return true
}
