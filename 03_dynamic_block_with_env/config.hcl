application "gitlab" "http" {
    hostname = "gitlab.com"
    port     = 443
    tls      = true
    authtype = "http-token"
    auth {
        token = "${env.GITLAB_HTTP_TOKEN}"
    }
}

application "github" "ssh" {
    hostname  = "github.com"
    port      = 22
    tls       = false
    authtype  = "userpass"
    auth {
        username = "${env.GITHUB_SSH_USER}"
        password = "*************"
    }
}

application "bitbucket" "http" {
    hostname  = "bitbucket.org"
    port      = 443
    tls       = true
}
