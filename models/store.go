package models

type Store struct {
	ID         int
	Name       string
	PointStep  int
	PointRatio float32
}

type Totals struct {
	UserID  int
	StoreID int
	Points  int
	Cash    float32
}

type Overview struct {
	Stores []Store
	Totals []Totals
}

type Activity struct {
	UserID  int
	StoreID int
	//VisitDate   date
	PointRedeem int
	CashRedeem  int
	PointsAdded int
}

type Scanned struct {
	UserID  int
	StoreID int
}

type Response struct {
	StoreID     int
	Name        string
	PointStep   int
	PointRatio  float32
	TotalCash   float32
	TotalPoints int
}

type Redeem struct {
	StoreID      int
	UserID       int
	RedeemPoints int
	RedeemCash   float32
}
