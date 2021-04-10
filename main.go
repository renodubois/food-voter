package main

import ( "fmt"
	 "log"
	 "io"
	 "net/http"
	 "html/template"
	 "strings"
	 _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	 "math/rand"
	 "time"
)

func isValidSlug(slug string) bool {
	log.Printf("Slug: %s", slug)
	slugSectionLength := 6
	numberOfSections := 6
	numberOfDashes := numberOfSections - 1
	if (len(slug) != (slugSectionLength * numberOfSections + numberOfDashes)) {
		log.Println("Overall length doesn't match")
		return false
	}
	splitSlug := strings.Split(slug, "-")
	if (len(splitSlug) != 6) {
		log.Println("Number of sections doesn't match")
		return false
	}
	for _, section := range(splitSlug) {
		// TODO: if I want to do more, check to see if each section contains the char set
		if (len(section) != 6) {
			log.Println("Length of section doesn't match")
			return false
		}
	}
	return true
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if (r.Method == "" || r.Method == "GET") {
		if (len(path) == 0) {
			// index
			t, _ := template.ParseFiles("index.html")
			t.Execute(w, nil)
			return
		} else if (isValidSlug(path)) {
			// get info about page
			// load
			getPoll(path)
			fmt.Fprintf(w, "Other Pages")
			return
		}
		http.Error(w, "404: Not Found", 404)
		return
	} else if (r.Method == "POST") {
		// Handle submits
		if (strings.HasPrefix(path, "submit-poll")) {
			// TODO: poll submitting
			return
		} else if (strings.HasPrefix(path, "create-poll")) {
			body, _ := io.ReadAll(r.Body)
			log.Println(string(body))
			options := parseBody(string(body))
			log.Println(options)
			// create poll in DB
			// create slug, redirect to slug
			makeNewPoll(options)
			// Make a poll
			t, _ := template.ParseFiles("poll.html")
			t.Execute(w, options)
			return
		} else {
			http.Error(w, "Bad Request", 400)
		}
	} else {
		http.Error(w, "Bad Request", 400)
		return
	}
}

func makeSlug () string {
	// reno's propreitary slug format (rpsf)
	// sections are of length 6, can be a-Z,0-9
	// there are 6 sections
	result := ""
	charset := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	charset_len := len(charset)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			result += string(charset[rand.Intn(charset_len)])
		}
		if (i == 5) {
			continue
		}
		result += "-"
	}
	return result
}

func getPoll(slug string) {
	query := "SELECT * FROM poll WHERE voter_slug='?';"
	db, err := sql.Open("sqlite3", "food_voter.db")
	if err != nil {
		log.Println("error when connecting to db")
		log.Fatal(err)
		return
	}
	defer db.Close()
	rows, db_err := db.Query(query, slug)
	if db_err != nil {
		log.Fatal(db_err)
	}
	for rows.Next() {
		// variables for the row
		var (
			id int
			options string
			creatorSlug string
			voterSlug string
		)
		if err := rows.Scan(&id, &options, &creatorSlug, &voterSlug); err != nil {
			log.Fatal(err)
		}
		log.Println(options)
	}
	log.Println(query)
}

func makeNewPoll(options []string) bool {
	// make slug
	creator_slug := makeSlug()
	voter_slug := makeSlug()
	options_string := strings.Join(options, ",")
	// insert into DB
	query := "INSERT INTO poll (options, creator_slug, voter_slug) VALUES ('" + options_string + "', '" + creator_slug + "', '" + voter_slug + "');"
	db, err := sql.Open("sqlite3", "food_voter.db")
	if err != nil {
		log.Println("error when connecting to db")
		log.Fatal(err)
		return false
	}
	defer db.Close()

	_, db_err := db.Exec(query)
	if (db_err != nil) {
		log.Println("error when executing query")
		log.Fatal(db_err)
		return false
	}
	return true
}

func parseBody(body string) []string {
	splits := strings.Split(body, "&")
	if (len(splits) == 5) {
		// get each value
		var options []string
		for _, val := range splits {
			value := val[8:]
			if (len(value) == 0) {
				continue
			}
			value = strings.ReplaceAll(value, "+", " ")
			// TODO: Replace escaped characters
			// TODO: Clean values for S A F E T Y
			options = append(options, value)
		}
		return options
	}
	return make([]string, 0)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
