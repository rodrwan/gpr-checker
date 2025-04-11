# gpr-checker

`gpr-checker` es una herramienta CLI que revisa automáticamente Pull Requests en GitHub usando inteligencia artificial (GPT-4), dejando comentarios con sugerencias constructivas directamente en los PRs.

Ideal para desarrolladores que quieran automatizar revisiones de código y mejorar la calidad de sus proyectos con ayuda de IA.

---

## 🚀 Instalación

Compila desde el código fuente:

```bash
go build -o gpr-checker main.go
```

---

## 🧠 Uso

### ✅ Desde terminal local

```bash
./gpr-checker \
  --repo tu-usuario/tu-repo \
  --pr 123 \
  --openai-key sk-... \
  --github-token ghp_...
```

O usando variables de entorno (recomendado):

```bash
export GITHUB_REPOSITORY=tu-usuario/tu-repo
export PR_NUMBER=123
export OPENAI_API_KEY=sk-...
export GITHUB_TOKEN=ghp_...

./gpr-checker
```

---

## 🔑 Cómo obtener el GitHub Token

Para ejecutar esta herramienta **fuera de GitHub Actions**, necesitas un **token personal de GitHub (PAT)** con permisos para comentar en PRs.

### 🎯 Instrucciones

1. Ve a: [https://github.com/settings/tokens](https://github.com/settings/tokens)
2. Haz clic en **"Generate new token (classic)"**
3. Selecciona al menos el siguiente scope:
   - ✅ `repo` (acceso completo a repos repositorios)
4. Genera el token y guárdalo en un lugar seguro
5. Usa ese token como `--github-token` o en la variable `GITHUB_TOKEN`

> **Nota**: No compartas este token. Trátalo como una contraseña.

---

## ⚙️ Uso desde GitHub Actions

El flujo automático ya está incluido en `.github/workflows/pr-review.yml` y usa el token `GITHUB_TOKEN` que GitHub proporciona por defecto (no necesitas configurarlo tú mismo).

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
  PR_NUMBER: ${{ github.event.pull_request.number }}
  GITHUB_REPOSITORY: ${{ github.repository }}
```

---

## 🔐 Seguridad

Antes de ejecutar su lógica principal, `gpr-checker` valida que:

- El GitHub Token es válido
- Tiene acceso al repositorio
- Puede comentar en el PR

Esto evita fallos silenciosos y te asegura que todo esté correctamente configurado.

---

## 📦 Releases automáticos

Si marcas un tag como `v0.1.0` y lo haces push:

```bash
git tag -a v0.1.0 -m "primer release"
git push origin v0.1.0
```

Se activará un flujo CI/CD que compilará binarios para Linux, macOS y Windows y los subirá como un release público.

---

## 🛠 Requisitos

- Go 1.24+
- Una clave API válida de OpenAI (`sk-...`)
