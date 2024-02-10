CREATE VIEW consistency AS
SELECT t.customer_id,c.amount as customer_amount, t.amount as transactions_amount, t.amount = c.amount as is_consistent
FROM customers c
         LEFT JOIN (SELECT sum(CASE WHEN t."type" = 'c' THEN t.amount else -t.amount end) as amount, t.customer_id
                    FROM transactions t
                    GROUP BY t.customer_id) t ON c.id = t.customer_id