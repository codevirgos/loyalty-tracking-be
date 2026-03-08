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

type Activity struct {
	UserID  int
	StoreID int
	//VisitDate   date
	PointRedeem int
	CashRedeem  int
	PointsAdded int
}
