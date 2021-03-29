application "gitlab" "http" {
    hostname = "gitlab.com"
    port     = 443
    tls      = true
    authtype = "http-token"
    auth {
        token = "my-super-token"
    }
}

application "github" "ssh" {
    hostname  = "github.com"
    port      = 22
    tls       = false
    authtype  = "userpass"
    auth {
        username = "git"
        password = "*************"
    }
}
