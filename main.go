package main

import ( "log"
	 "io"
	 "net/http"
	 "html/template"
	 "strings"
	 _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	 "math/rand"
	 "time"
)

type pollRow struct {
	Id int
	Options []string
	CreatorSlug string
	VoterSlug string
}


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
			t, _ := template.ParseFiles("templates/index.html")
			t.Execute(w, nil)
			return
		} else if (isValidSlug(path)) {
			// get info about page
			// load
			pollInfo := getPoll(path)
			// at this point, poll is a voting poll
			t, err := template.ParseFiles("templates/poll_vote.html")
			if err != nil {
				log.Fatal(err)
			}
			err = t.Execute(w, pollInfo)
			if err != nil {
				log.Fatal(err)
			}
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
			t, _ := template.ParseFiles("templates/poll_view.html")
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


func getPoll(slug string) pollRow {
	query := "SELECT * FROM poll WHERE voter_slug='"+slug+"'"
	db, err := sql.Open("sqlite3", "food_voter.db")
	if err != nil {
		log.Println("error when connecting to db")
		log.Fatal(err)
		return pollRow{}
	}
	defer db.Close()
	rows, db_err := db.Query(query)
	if db_err != nil {
		log.Fatal(db_err)
		return pollRow{}
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
		formattedOptions := optionsStringToSlice(options)
		res := pollRow{id, formattedOptions, creatorSlug, voterSlug}
		return res
	}
	return pollRow{}
}

func makeNewPoll(options []string) bool {
	// make slug
	creatorSlug := makeSlug()
	voterSlug := makeSlug()
	optionsString := strings.Join(options, ",")
	// insert into DB
	query := "INSERT INTO poll (options, creator_slug, voter_slug) VALUES ('" + optionsString + "', '" + creatorSlug + "', '" + voterSlug + "');"
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
	log.Println("Voter: " + voterSlug)
	log.Println("Creator: " + creatorSlug)
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

func optionsStringToSlice (options string) []string {
	return strings.Split(options, ",")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
