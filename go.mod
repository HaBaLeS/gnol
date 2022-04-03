module github.com/HaBaLeS/gnol

go 1.14

require (
	github.com/HaBaLeS/go-logger v1.3.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/duo-labs/webauthn v0.0.0-20210727191636-9f1b88ef44cc
	github.com/gen2brain/go-fitz v0.0.0-20190716092309-62357ab3d4a9
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/mholt/archiver/v3 v3.5.1-0.20201230180942-1ee1dbd58314
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/rs/xid v1.2.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/stretchr/testify v1.7.0
	github.com/teris-io/cli v1.0.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/tools v0.1.9 // indirect
)

//Security fix
replace github.com/dgrijalva/jwt-go v3.2.0+incompatible => github.com/golang-jwt/jwt/v4 v4.1.0
