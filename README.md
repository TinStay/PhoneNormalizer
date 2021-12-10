# Phone Number Normalizer

## Exercise details

This exercise is fairly straight-forward - a program that will iterate through a database and normalize all of the phone numbers in it. After normalizing all of the data one might find that there are duplicates, so then those duplicates must be removed and only one entry should exist in our database.

### What will be covered

- Writing raw SQL and using the [database/sql](https://golang.org/pkg/database/sql/) package in the standard library
- Using the very popular [sqlx](https://github.com/jmoiron/sqlx) third party package, which is basically an extension of Go's sql package.
- Using a relatively minimalistic ORM (I will be using [gorm](https://github.com/jinzhu/gorm))

Inside a table in our SQL database following entries will be added and used for the exercise :

```
1234567890
123 456 7891
(123) 456 7892
(123) 456-7893
123-456-7894
123-456-7890
1234567892
(123)456-7892
```
