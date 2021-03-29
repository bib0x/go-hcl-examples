package main

import (
    "fmt"
    "os"

    "github.com/hashicorp/hcl/v2"
    "github.com/hashicorp/hcl/v2/gohcl"
    "github.com/hashicorp/hcl/v2/hclsimple"
    "github.com/zclconf/go-cty/cty"
)

type Config struct {
    Apps []ApplicationHCL `hcl:"application,block"`
}

type AuthHttpToken struct {
    Token string    `hcl:"token"`
}

type AuthUserPass struct {
    Username string  `hcl:"username"`
    Password string  `hcl:"password"`
}

type Auth interface {
    GetCredentials() string
}

func (auth *AuthHttpToken) GetCredentials() string {
    return auth.Token 
}

func (auth *AuthUserPass) GetCredentials() string {
    return fmt.Sprintf("%s:%s", auth.Username, auth.Password)
}

type ApplicationHCL struct {
    Name     string `hcl:"application_name,label"`
    Proto    string `hcl:"proto_name,label"`
    Hostname string `hcl:"hostname"`
    Port     int    `hcl:"port"`
    Tls      bool   `hcl:"tls"`
    AuthType string `hcl:"authtype,optional"`
    AuthHCL *struct {
        HCL     hcl.Body   `hcl:",remain"`
    } `hcl:"auth,block"`
}

type Application struct {
    Name     string 
    Proto    string 
    Hostname string 
    Port     int    
    Tls      bool   
    AuthType string 
    Auth     Auth
}

func NewAuth(aHCL *ApplicationHCL) (Auth, error) {
    switch (aHCL.AuthType) {
        case "http-token":
            return &AuthHttpToken{}, nil
        case "userpass":
            return &AuthUserPass{}, nil
        default:
            return nil, fmt.Errorf(
                "error in AuthFactory invalid authtype: %v", aHCL.AuthType,
            )
    }
}

func AuthFactory(aHCL *ApplicationHCL) (Auth, error) {
    auth, err := NewAuth(aHCL)
    if err != nil {
        return auth, err
    }
    if aHCL.AuthHCL != nil && aHCL.AuthHCL.HCL != nil {
        ctx, err := createContext()
        if err != nil {
            return auth, fmt.Errorf(
                "error in AuthFactory creating HCL eval context: %w", err,
            )
        }
        diag := gohcl.DecodeBody(aHCL.AuthHCL.HCL, ctx, auth)
        if diag.HasErrors() {
            return auth, fmt.Errorf(
                "error in AuthFactory parsing HCL: %w", diag,
            )
        }
    }
    return auth, nil
}

func NewApplication(aHCL *ApplicationHCL) (*Application) {
    auth, _ := AuthFactory(aHCL)
    return &Application{
        Name: aHCL.Name,
        Proto: aHCL.Proto,
        Hostname: aHCL.Hostname,
        Port: aHCL.Port,
        Tls: aHCL.Tls,
        AuthType: aHCL.AuthType,
        Auth: auth,
    }
}

const (
    GITLAB_HTTP_TOKEN = "GITLAB_HTTP_TOKEN"
    GITHUB_SSH_USER   = "GITHUB_SSH_USER"
)

func createContext() (*hcl.EvalContext, error) {
    var gitlabHttpToken string
    if os.Getenv(GITLAB_HTTP_TOKEN) != "" {
        gitlabHttpToken = os.Getenv(GITLAB_HTTP_TOKEN)
    }

    githubSshUser := "git"
    if os.Getenv(GITHUB_SSH_USER) != "" {
        githubSshUser = os.Getenv(GITHUB_SSH_USER)
    }

    variables := map[string]cty.Value{
        "env": cty.ObjectVal(map[string]cty.Value{
            GITLAB_HTTP_TOKEN: cty.StringVal(gitlabHttpToken),
            GITHUB_SSH_USER: cty.StringVal(githubSshUser),
        }),
    }

    return &hcl.EvalContext{
        Variables: variables,
    }, nil
}

func main() {
   var config Config

   if err := hclsimple.DecodeFile("config.hcl", nil, &config); err == nil {
        fmt.Printf("Raw %#v\n\n", config) 

        for _, appHCL := range config.Apps {
            app := NewApplication(&appHCL)
            fmt.Printf("App Name: %v\n", app.Name)
            fmt.Printf("App Proto: %v\n", app.Proto)
            fmt.Printf("App Port: %v\n", app.Port)
            fmt.Printf("App TLS: %v\n", app.Tls)
            if app.AuthType != "" {
                fmt.Printf("App AuthType: %v\n", app.AuthType)
                fmt.Printf("App Credentials: %v\n", app.Auth.GetCredentials())
            }
            fmt.Printf("\n")
        }
   }
}
