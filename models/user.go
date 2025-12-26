package models

type User struct {
	Number      int
	TotalPoints int
	Visits      int
	Name        string
}

type Redeem struct {
	Number int
	Amount int
	// cashvalue etc...
}

func (u User) RedeemPoints(redeemAmount int) *User {
	remaining := u.TotalPoints - redeemAmount
	return &User{
		Number:      u.Number,
		TotalPoints: remaining,
		Visits:      u.Visits,
	}
}

func (u User) AddVisit(pointAmount int) *User {
	u.Visits++
	t := u.TotalPoints + pointAmount
	return &User{
		Number:      u.Number,
		TotalPoints: t,
		Visits:      u.Visits,
	}
}
