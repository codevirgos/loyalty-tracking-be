package repo

import (
	"Payback_BE/models"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
)

type StoreRepo interface {
	RedeemPoints(ctx context.Context, redeemed models.Redeem) (*models.Overview, error)
	StoreScan(ctx context.Context, scanned models.Scanned) (*models.Overview, error)
	GetOverview(ctx context.Context, userId int) (*models.Overview, error)
	GetStoreInfo(ctx context.Context, storeIds []int) ([]models.Store, error)
}

type pgStoreRepo struct {
	db *sql.DB
}

func NewPgStoreRepo(db *sql.DB) *pgStoreRepo {
	return &pgStoreRepo{
		db: db,
	}
}

func (sr *pgStoreRepo) GetOverview(ctx context.Context, userId int) (*models.Overview, error) {
	query := `Select user_id, store_id, total_points, total_cash FROM track_totals WHERE user_id = $1`

	var totals []models.Totals
	rows, err := sr.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	var storeIds []int

	for rows.Next() {
		var tot models.Totals
		if err := rows.Scan(&tot.UserID, &tot.StoreID, &tot.Points, &tot.Cash); err != nil {
			return nil, err
		}
		storeIds = append(storeIds, tot.StoreID)
		totals = append(totals, tot)
	}

	if len(storeIds) == 0 {
		return nil, errors.New("No rows round")
	}

	stores, err := sr.GetStoreInfo(ctx, storeIds)

	overview := &models.Overview{
		Totals: totals,
		Stores: stores,
	}
	return overview, nil
}

func (sr *pgStoreRepo) GetStoreInfo(ctx context.Context, storeIds []int) ([]models.Store, error) {
	query := `Select id, business_name, point_step, point_to_cash FROM store where id = ANY($1)`

	var stores []models.Store

	rows, err := sr.db.QueryContext(ctx, query, pq.Array(storeIds))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var sInfo models.Store
		if err := rows.Scan(&sInfo.ID, &sInfo.Name, &sInfo.PointStep, &sInfo.PointRatio); err != nil {
			return nil, err
		}
		stores = append(stores, sInfo)
	}

	return stores, rows.Err()
}

func (sr *pgStoreRepo) StoreScan(ctx context.Context, scanned models.Scanned) (*models.Overview, error) {
	query := `Select t.store_id, s.business_name, t.total_cash, t.total_points, s.point_step, s.point_to_cash FROM track_totals t JOIN store s ON t.store_id = s.id WHERE user_id = $1 and store_id = $2`

	var rs models.Response
	err := sr.db.QueryRowContext(ctx, query, scanned.UserID, scanned.StoreID).Scan(
		&rs.StoreID,
		&rs.Name,
		&rs.TotalCash,
		&rs.TotalPoints,
		&rs.PointStep,
		&rs.PointRatio,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	// add the total points
	newPointTotal := rs.PointStep + rs.TotalPoints

	// add the total cash
	newCashTotal := rs.TotalCash + float32(rs.PointStep)/rs.PointRatio

	tt_update := `Update track_totals set total_points = $1, total_cash = $2 WHERE user_id = $3 and store_id = $4`
	_, err = sr.db.Exec(tt_update, newPointTotal, newCashTotal, scanned.UserID, scanned.StoreID)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	// response objects
	store := &models.Store{
		Name:       rs.Name,
		ID:         rs.StoreID,
		PointStep:  rs.PointStep,
		PointRatio: rs.PointRatio,
	}

	totals := &models.Totals{
		Points: newPointTotal,
		Cash:   newCashTotal,
	}

	overView := &models.Overview{
		Stores: []models.Store{*store},
		Totals: []models.Totals{*totals},
	}
	return overView, nil
}

func (sr *pgStoreRepo) RedeemPoints(ctx context.Context, redeemed models.Redeem) (*models.Overview, error) {
	query := `Select s.point_step, s.point_to_cash, t.total_cash, t.total_points FROM store s JOIN track_totals t ON s.id = t.store_id where user_id = $1 and store_id = $2`

	var rs models.Response
	err := sr.db.QueryRowContext(ctx, query, redeemed.UserID, redeemed.StoreID).Scan(
		&rs.PointStep,
		&rs.PointRatio,
		&rs.TotalCash,
		&rs.TotalPoints,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	// subtract the total points
	newPointTotal := rs.TotalPoints - redeemed.RedeemPoints

	// subtract the total cash
	newCashTotal := rs.TotalCash - float32(redeemed.RedeemPoints)/rs.PointRatio

	var k = newCashTotal - redeemed.RedeemCash
	// double check
	if (k) > .005 {
		// somethings wrong
	}

	tt_update := `Update track_totals set total_points = $1, total_cash = $2 WHERE user_id = $3 and store_id = $4`
	_, err = sr.db.Exec(tt_update, newPointTotal, newCashTotal, redeemed.UserID, redeemed.StoreID)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	store := &models.Store{
		ID:         redeemed.StoreID, // why not
		PointStep:  rs.PointStep,
		PointRatio: rs.PointRatio,
	}

	totals := &models.Totals{
		Points: newPointTotal,
		Cash:   newCashTotal,
	}

	overView := &models.Overview{
		Stores: []models.Store{*store},
		Totals: []models.Totals{*totals},
	}
	return overView, nil
}
