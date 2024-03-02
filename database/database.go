package database

import "github.com/gocql/gocql"

func FetchAllProfiles(session *gocql.Session) ([]Profile, error) {
	var profiles []Profile
	iter := session.Query("SELECT * FROM website.profiles").Iter()

	profile := Profile{}
	for iter.Scan(&profile.Profile, &profile.ProfileImage, &profile.Code, &profile.MainSummary) {
		profiles = append(profiles, profile)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return profiles, nil
}

