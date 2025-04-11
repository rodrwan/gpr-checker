# gpr-checker

Una herramienta CLI que revisa Pull Requests en GitHub con inteligencia artificial (GPT-4) y deja comentarios automáticos.

## Uso

Compila:

```bash
go build -o gpr-checker main.go
```

Ejecuta:

```bash
./gpr-checker --repo usuario/repo --pr 123 --openai-key sk-... --github-token ghp_...
```

O usando variables de entorno:

```bash
export GITHUB_REPOSITORY=usuario/repo
export PR_NUMBER=123
export OPENAI_API_KEY=sk-...
export GITHUB_TOKEN=ghp_...

./gpr-checker
```

## GitHub Actions

Incluye un flujo automático en `.github/workflows/pr-review.yml`
