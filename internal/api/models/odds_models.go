package models

type CreateOddsRequest struct {
	League          string  `json:"league"`
	HomeTeam        string  `json:"home_team"`
	AwayTeam        string  `json:"away_team"`
	HomeTeamWinOdds float64 `json:"home_team_win_odds"`
	AwayTeamWinOdds float64 `json:"away_team_win_odds"`
	DrawOdds        float64 `json:"draw_odds"`
	GameDate        string  `json:"game_date"`
}

type CreateOddsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ReadOddsRequest struct {
	League string `json:"league"`
	Date   string `json:"date"`
}

type ReadOddsResponse struct {
	Odds    []CreateOddsRequest `json:"odds"`
	Details string              `json:"details"`
}

type UpdateOddsRequest struct {
	League          string  `json:"league"`
	HomeTeam        string  `json:"home_team"`
	AwayTeam        string  `json:"away_team"`
	HomeTeamWinOdds float64 `json:"home_team_win_odds"`
	AwayTeamWinOdds float64 `json:"away_team_win_odds"`
	DrawOdds        float64 `json:"draw_odds"`
	GameDate        string  `json:"game_date"`
}

type UpdateOddsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type DeleteOddsRequest struct {
	League   string `json:"league"`
	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`
	GameDate string `json:"game_date"`
}

type DeleteOddsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details"`
}
