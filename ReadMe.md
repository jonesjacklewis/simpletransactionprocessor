# Simple Transaction Processor

This project will read transaction from a CSV file in the below format:

```csv
customer_id,transaction_id,amount,transaction_type,timestamp
acc_123,tx_1,100.00,CREDIT,2024-10-27T10:00:00Z
acc_456,tx_2,50.00,CREDIT,2024-10-27T10:01:00Z
acc_123,tx_3,25.50,DEBIT,2024-10-27T10:02:00Z
acc_456,tx_4,10.00,CREDIT,2024-10-27T10:03:00Z
acc_789,tx_5,200.00,CREDIT,2024-10-27T10:04:00Z
acc_123,tx_6,5.00,DEBIT,2024-10-27T10:05:00Z
```

and process the transactions, such that `CREDIT` represents an increase in balance, and `DEBIT` represents a decrease in balance. It should then calculate the final balance for each `customer_id` in this format:

```csv
acc_123,69.50
acc_456,60.00
acc_789,200.00
```

## To Do

- [ ] Create ReadMe
- [ ] Check for existence of the CSV file
- [ ] Generate test CSV file on load, if one does not exist
- [ ] Stream each line of the CSV file
- [ ] Extract "useful" information (customer_id, amount, transaction_type)
- [ ] Update a map of values with `key=custmomer_id, value=balance`
- [ ] Sort map by `customer_id`
- [ ] Output the map in the requested format