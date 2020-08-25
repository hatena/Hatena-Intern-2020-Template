#!/bin/bash
set -eux
mysqladmin -uroot create intern_2020_account
mysqladmin -uroot create intern_2020_account_test
mysql -uroot intern_2020_account < /config/schema.sql
mysql -uroot intern_2020_account_test < /config/schema.sql
