package repository

import (
	"database/sql"
	"fmt"

	"github.com/43nvy/wb_l0"
	"github.com/jmoiron/sqlx"
)

type OrderPG struct {
	db *sqlx.DB
}

func NewOrderPG(db *sqlx.DB) *OrderPG {
	return &OrderPG{db: db}
}

func (o *OrderPG) CreateOrder(order wb_l0.Order) (int, error) {
	var id int

	tx, err := o.db.Begin()
	if err != nil {
		return 0, err
	}

	deliveryID, err := insertDelivery(tx, order.Delivery)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	paymentID, err := insertPayment(tx, order.Payment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	itemsID, err := insertItems(tx, order.Items)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (order_uid, track_number, entry, delivery_id, payment_id, items_id,
                            locale, internal_signature, customer_id, delivery_service, shardkey,
                            sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id
    `, ordersTable)

	err = tx.QueryRow(query,
		order.OrderUID, order.TrackNumber, order.Entry, deliveryID, paymentID, itemsID,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SMID, order.DateCreated, order.OOFShard,
	).Scan(&id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (o *OrderPG) GetOrder(order_id int) (wb_l0.Order, error) {
	query := `
		SELECT
			o.id,
			o.order_uid,
			o.track_number,
			o.entry,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shardkey,
			o.sm_id,
			o.date_created,
			o.oof_shard,
			d.id AS "delivery.id",
			d.name AS "delivery.name",
			d.phone AS "delivery.phone",
			d.zip AS "delivery.zip",
			d.city AS "delivery.city",
			d.address AS "delivery.address",
			d.region AS "delivery.region",
			d.email AS "delivery.email",
			p.id AS "payment.id",
			p.transaction AS "payment.transaction",
			p.request_id AS "payment.request_id",
			p.currency AS "payment.currency",
			p.provider AS "payment.provider",
			p.amount AS "payment.amount",
			p.payment_dt AS "payment.payment_dt",
			p.bank AS "payment.bank",
			p.delivery_cost AS "payment.delivery_cost",
			p.goods_total AS "payment.goods_total",
			p.custom_fee AS "payment.custom_fee",
			i.id AS "items.id",
			i.chrt_id AS "items.chrt_id",
			i.track_number AS "items.track_number",
			i.price AS "items.price",
			i.rid AS "items.rid",
			i.name AS "items.name",
			i.sale AS "items.sale",
			i.size AS "items.size",
			i.total_price AS "items.total_price",
			i.nm_id AS "items.nm_id",
			i.brand AS "items.brand",
			i.status AS "items.status"
		FROM
			orders o
		LEFT JOIN delivery d ON o.delivery_id = d.id
		LEFT JOIN payments p ON o.payment_id = p.id
		LEFT JOIN items i ON o.items_id = i.id
		WHERE
			o.id = $1
	`

	rows, err := o.db.Queryx(query, order_id)
	if err != nil {
		return wb_l0.Order{}, err
	}
	defer rows.Close()

	var order wb_l0.Order

	for rows.Next() {
		var delivery wb_l0.Delivery
		var payment wb_l0.Payment
		var item wb_l0.Item

		err := rows.Scan(
			&order.ID,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
			&delivery.ID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.ID,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
			&item.ID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return wb_l0.Order{}, err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = []wb_l0.Item{item}
	}

	return order, nil
}

func (o *OrderPG) GetTenOrders() ([]wb_l0.Order, error) {
	query := `
		SELECT
			o.id,
			o.order_uid,
			o.track_number,
			o.entry,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shardkey,
			o.sm_id,
			o.date_created,
			o.oof_shard,
			d.id AS "delivery.id",
			d.name AS "delivery.name",
			d.phone AS "delivery.phone",
			d.zip AS "delivery.zip",
			d.city AS "delivery.city",
			d.address AS "delivery.address",
			d.region AS "delivery.region",
			d.email AS "delivery.email",
			p.id AS "payment.id",
			p.transaction AS "payment.transaction",
			p.request_id AS "payment.request_id",
			p.currency AS "payment.currency",
			p.provider AS "payment.provider",
			p.amount AS "payment.amount",
			p.payment_dt AS "payment.payment_dt",
			p.bank AS "payment.bank",
			p.delivery_cost AS "payment.delivery_cost",
			p.goods_total AS "payment.goods_total",
			p.custom_fee AS "payment.custom_fee",
			i.id AS "items.id",
			i.chrt_id AS "items.chrt_id",
			i.track_number AS "items.track_number",
			i.price AS "items.price",
			i.rid AS "items.rid",
			i.name AS "items.name",
			i.sale AS "items.sale",
			i.size AS "items.size",
			i.total_price AS "items.total_price",
			i.nm_id AS "items.nm_id",
			i.brand AS "items.brand",
			i.status AS "items.status"
		FROM
			orders o
		LEFT JOIN delivery d ON o.delivery_id = d.id
		LEFT JOIN payments p ON o.payment_id = p.id
		LEFT JOIN items i ON o.items_id = i.id
		ORDER BY o.id DESC
		LIMIT 10
	`

	rows, err := o.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []wb_l0.Order

	for rows.Next() {
		var order wb_l0.Order
		var delivery wb_l0.Delivery
		var payment wb_l0.Payment
		var item wb_l0.Item

		err := rows.Scan(
			&order.ID,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
			&delivery.ID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.ID,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
			&item.ID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = []wb_l0.Item{item}

		orders = append(orders, order)
	}

	return orders, nil
}

func insertDelivery(tx *sql.Tx, delivery wb_l0.Delivery) (int, error) {
	var deliveryID int
	query := `
        INSERT INTO delivery (name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
	err := tx.QueryRow(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email).Scan(&deliveryID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return deliveryID, nil
}

func insertPayment(tx *sql.Tx, payment wb_l0.Payment) (int, error) {
	var paymentID int
	query := `
        INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt,
                            bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `
	err := tx.QueryRow(query, payment.Transaction, payment.RequestID, payment.Currency,
		payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank, payment.DeliveryCost,
		payment.GoodsTotal, payment.CustomFee).Scan(&paymentID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return paymentID, nil
}

func insertItems(tx *sql.Tx, items []wb_l0.Item) (int, error) {
	var item_id int
	query := "INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES "

	values := make([]interface{}, 0, len(items)*11)
	for _, item := range items {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d), ",
			len(values)+1, len(values)+2, len(values)+3, len(values)+4, len(values)+5,
			len(values)+6, len(values)+7, len(values)+8, len(values)+9, len(values)+10, len(values)+11)

		values = append(values,
			item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
	}

	query = query[:len(query)-2]

	query += " RETURNING id"

	err := tx.QueryRow(query, values...).Scan(&item_id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return item_id, nil
}
