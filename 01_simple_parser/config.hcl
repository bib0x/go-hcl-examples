type = "example"
name = "simple-parser"

resource "my-app" "user" {
    name    = "toto"
    state   = "blocked"

    task "verify" {
        connector = "http"
    }
}
