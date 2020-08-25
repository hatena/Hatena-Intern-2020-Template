#!/bin/bash
set -eux
mysqladmin -uroot create intern_2020_blog
mysqladmin -uroot create intern_2020_blog_test
mysql -uroot intern_2020_blog < /config/schema.sql
mysql -uroot intern_2020_blog_test < /config/schema.sql
