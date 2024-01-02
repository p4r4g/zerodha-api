## Developer Settings

#### Submit package

> Commit changes
> git tag v0.2.0
> git push origin --tags
> GOPROXY=proxy.golang.org go list -m github.com/parag-b/zerodha-api@v0.2.0

#### Use local package

> go mod edit -replace=github.com/parag-b/[zerodha-api@v0.0.0-unpublished](mailto:zerodha-api@v0.0.0-unpublished)\=/repos/zerodha-api/
> go mod tidy
