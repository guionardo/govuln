# govulncheck pre-commit

pre-commit hook for vulnerabilites checking in golang projects.

This hook will use the [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) for checking vulnerabilities in the current golang project.

Just add this to your .pre-commit-config.yaml file:

```yaml
repos:
    - repo: https://github.com/guionardo/govuln
      rev: v0.0.9
      hooks:
          - id: go-vulncheck
```

If you don't want to block the commit when some vulnerability is found, use this hook:

```yaml
repos:
    - repo: https://github.com/guionardo/govuln
      rev: v0.0.9
      hooks:
          - id: go-vulncheck-just-warn
```
