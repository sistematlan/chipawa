# BACKLOG — chipawa

> Pendientes priorizados. Ítems marcados con `[x]` están completados.
> Última actualización: sesión 2026-06-08.

---

## ✅ Completado en sesiones previas

- [x] Scaffold inicial: `scan`, `apps`, `caches`, `projects`
- [x] Tipo `item.Item` unificado con `Category` y `Risk`
- [x] Detectores ampliados de caches (uv, pnpm, npm `_npx`, JetBrains old versions, Cargo registry, iOS DeviceSupport, CoreSimulator, Xcode Archives)
- [x] `internal/orphans` con detectores: Docker leftover, WhatsApp media
- [x] Comando `orphans`
- [x] `internal/cleaner` con `Mode` (DryRun/Interactive/Yes), `Remover`, `Plan`, `Summary`
- [x] `PathRemover` con guardia `SafeRoots`; `DockerPruneRemover` (sin `--volumes`)
- [x] `TerminalPrompter` con `[s/N/v=ver/q=salir]` y confirmación por nombre para `RiskDangerous`
- [x] Comando `clean` con `--dry-run`, `--yes`, `--include-orphans`
- [x] Bug fix: scanner reutilizado en `TerminalPrompter` (no se descarta el buffer entre prompts)
- [x] Tests unitarios (10 casos en cleaner, parser de tamaños Docker, splitJBVersion, fixtures de orphans)

---

## 🔥 Próxima sesión — completar MVP

### 1. Comando `downloads` (P0)

Clasifica `~/Downloads` para identificar candidatos a borrar.

- [ ] `internal/downloads/downloads.go` con detector
- [ ] Agrupar por extensión: `.dmg`, `.pkg`, `.zip`, `.sql`, `.cbr`, `.mov`, `.mp4`, `.iso`
- [ ] Agrupar por antigüedad: `<30d`, `30-90d`, `>90d`
- [ ] Detectar **ZIPs ya extraídos**: si existe carpeta con mismo nombre al lado del .zip
- [ ] Detectar **node_modules dentro de Downloads** (proyectos abandonados)
- [ ] Detectar **DMGs cuya app está instalada** (instalador ya usado)
- [ ] Detectar **.sql/.sql.bak grandes** como candidatos obvios
- [ ] Comando `cmd/downloads.go` con tabla agrupada
- [ ] Integrar items en `clean` (categoría `CategoryDownload`, `RiskAskBefore`)
- [ ] Tests con fixtures temporales

**Estimado:** 1-2 horas.

### 2. Flags globales + `report` (P0)

- [ ] Flag global `--path` (ya existe parcial en `projects`, normalizar)
- [ ] Flag global `--json` — toda salida puede emitirse como JSON
- [ ] Flag global `--no-color` — desactiva ANSI (cuando se agregue color)
- [ ] Flag global `--yes` — heredable (hoy solo está en `clean`)
- [ ] Comando `report` con `--out <file>` y formato estructurado
- [ ] Refactor: que cada `cmd/*.go` use un renderer común (text vs JSON)

**Estimado:** 2-3 horas. Refactor mediano porque toca todos los comandos.

### 3. Color y UI polish (P1)

- [ ] Lib: `github.com/fatih/color` o solo ANSI manual (preferible — menos deps)
- [ ] Verde para "safe", amarillo para "ask", rojo para "danger"
- [ ] Detección de TTY: solo color cuando stdout sea terminal
- [ ] Símbolos: ✓ ✗ → ⚠ (con fallback ASCII en `--no-color`)
- [ ] Tabla con anchos calculados dinámicamente (hoy son fijos y truncan feo)

---

## 🟡 Backlog medio plazo (post-MVP)

### Más detectores de caches

- [ ] **VSCode/Cursor/Trae caches**: `~/Library/Application Support/{Code,Cursor,Trae}/Cache*`
- [ ] **Electron caches**: `~/Library/Application Support/*/Cache`, `*/Code Cache`, `*/GPUCache` (los Electron son culpables silenciosos)
- [ ] **Slack/Discord caches**
- [ ] **Spotify cache**: `~/Library/Caches/com.spotify.client`
- [ ] **Time Machine local snapshots**: `tmutil listlocalsnapshots /` y borrado controlado
- [ ] **Mail attachments cache**: `~/Library/Mail/V*/MailData/`
- [ ] **gradle / maven**: `~/.gradle/caches`, `~/.m2/repository`
- [ ] **Bun cache**: `~/.bun/install/cache`
- [ ] **Deno cache**: `~/Library/Caches/deno`
- [ ] **rustup toolchains viejas**: parsear `rustup toolchain list` y detectar las no usadas
- [ ] **playwright browsers**: `~/Library/Caches/ms-playwright`

### Más detectores de orphans

- [ ] **Group Containers de apps desinstaladas**: cruzar `~/Library/Group Containers/*` con `/Applications`
- [ ] **Application Support de apps desinstaladas**: cruzar con `/Applications`
- [ ] **Wondershare leftover**: `~/Library/Application Support/com.wondershare.Installer`
- [ ] **Mega Limited cache** (la app sigue, pero el cache crece sin tope)
- [ ] **Adobe leftover** (típicamente >5 GB tras desinstalar)

### `projects` mejorado

- [ ] Detectar `node_modules`, `vendor/`, `target/`, `.next/`, `dist/`, `build/` dentro de cada proyecto
- [ ] Reportar tamaño de artefactos vs código real
- [ ] Detectar proyectos sin commit en >1 año como "abandonados"
- [ ] Sugerir `git clean -xdf` para los abandonados con git
- [ ] Modo `--clean-build-artifacts` que borra solo carpetas regenerables

### `apps` mejorado

- [ ] Tamaño TOTAL (app + Application Support + Containers + Caches asociados)
- [ ] Sugerir desinstalación con `mdls` last-used > 180d
- [ ] Comando `chipawa apps uninstall <name>` que llama a `mdfind` para encontrar todos los archivos relacionados

### Mejoras al cleaner

- [ ] Modo `--category cache|orphan|download|all` (filtro por categoría)
- [ ] Modo `--tool docker|jetbrains|...` (filtro por tool)
- [ ] Modo `--min-size 100M` (solo items grandes)
- [ ] Confirmación batch: "borrar todos los marcados con [s]" al final
- [ ] Undo log: registrar qué se borró y cuándo en `~/.chipawa/history.log`

---

## 🟢 Backlog largo plazo (visión)

### OSS launch

- [ ] LICENSE (MIT)
- [ ] README.md con: badges, features, screenshots, quickstart
- [ ] CONTRIBUTING.md con: cómo agregar nuevos detectores
- [ ] CODE_OF_CONDUCT.md
- [ ] GitHub Actions: `go test`, `go vet`, `staticcheck`, `golangci-lint`
- [ ] goreleaser config para multi-arch (amd64, arm64)
- [ ] Homebrew tap: `sistematlan/tools`
- [ ] Issue templates (bug, feature, new-detector)
- [ ] PR template
- [ ] Sitio mínimo en GitHub Pages o Vercel (landing + docs)

### Distribución y crecimiento

- [ ] Release v0.1.0 a Homebrew
- [ ] Post en HackerNews ("Show HN: chipawa — a transparent disk cleaner for macOS developers")
- [ ] Reddit: r/macOS, r/golang, r/programming
- [ ] Demo en video (asciinema) para README
- [ ] Documentación de cada detector con ejemplo de output

### Ideas más ambiciosas

- [ ] **Modo daemon**: ejecuta `scan` periódicamente y notifica cuando disco <20%
- [ ] **Plugins**: detectores externos cargados desde `~/.chipawa/plugins/*.so` o subprocess
- [ ] **Dashboard TUI** con bubbletea (alternativa al CLI puro)
- [ ] **Política de retención**: archivo de config (`~/.chipawa/config.toml`) con reglas tipo `npm cache: keep 30 days`
- [ ] **Comparativa pre/post**: muestra diff de uso por carpeta tras `clean`
- [ ] **Sugerencias inteligentes**: "tienes 12 carpetas con node_modules y solo trabajas en 3 proyectos activos"
- [ ] **Soporte Linux** (largo plazo, requiere abstraer detectores macOS-specific)
- [ ] **Soporte Windows** (lejos)

### Métricas de salud

- [ ] `chipawa doctor`: detecta problemas comunes (Docker sin volumes, brew obsoleto, etc.)
- [ ] `chipawa baseline`: guarda snapshot del estado actual para comparar después
- [ ] `chipawa diff <baseline>`: muestra qué creció/decreció desde un baseline

---

## 🐛 Deuda técnica conocida

- [ ] `cmd/orphans.go` y `cmd/scan.go` tienen escapes `\u00xx` literales — no afectan output (las strings con comillas dobles se procesan correctamente) pero quedan feos en el código fuente. Reemplazar por caracteres directos.
- [ ] `disk.DirSize` usa `du -sk` (subprocess). Para detectores que se llaman muchas veces (JetBrains versions) sumar latencia. Considerar implementación nativa con `filepath.Walk`.
- [ ] `caches.Scan()` llama a `du` secuencialmente. Paralelizar con goroutines + `errgroup` (gain ~3-5x).
- [ ] No hay manejo de errores diferenciado: si un detector falla, podría dejar la lista vacía. Cada detector debería retornar `(items, []error)` y `Scan()` agregar errores como warnings.
- [ ] Tests no cubren `Scan()` real (solo helpers); son frágiles ante cambios de paths del sistema. Considerar inyección de filesystem (afero o similar) para mockear.
- [ ] `cmd/scan.go` y `cmd/caches.go` duplican lógica de ordenado por `Bytes`. Mover a `internal/item`.

---

## 📌 Decisiones tomadas

- **Licencia**: MIT (a confirmar al hacer LICENSE)
- **Lenguaje**: Go 1.26
- **Distribución primaria**: Homebrew tap `sistematlan/tools`
- **CLI framework**: cobra
- **Telemetría**: cero, jamás
- **Modelo comercial**: pure OSS por ahora; sponsorship + consulting como vías futuras
- **`clean` por defecto**: solo caches; orphans requieren `--include-orphans`
- **UX de confirmación**: ítem por ítem con `[s/N/v/q]`, default seguro en NO
- **Docker**: `system prune -f` (sin `--volumes`); volumes solo con flag explícito futuro `--include-volumes`
- **SafeRoots**: solo `$HOME`, `/var/folders`, `/tmp` y `/private/*`. Cualquier otro path es rechazado.

---

## 🎯 Definition of Done para v0.1.0

Del SPEC original:

- [x] `chipawa scan` muestra resumen de disco y top categorías
- [x] `chipawa apps` lista apps con último uso y tamaño
- [x] `chipawa caches` detecta y totaliza caches de dev
- [x] `chipawa projects --path ~/sourcecode` reporta estado git de cada repo
- [x] `chipawa clean --dry-run` lista candidatos sin borrar nada
- [x] `chipawa clean` pide confirmación por ítem
- [ ] `chipawa downloads` clasifica Downloads por tipo y antigüedad
- [ ] `chipawa report --json` emite reporte estructurado
- [x] Binario único, sin dependencias externas en runtime
- [x] `make build` compila en < 10 segundos
- [x] Tests pasan con `make test`

**Faltan:** `downloads`, `report --json`, flags globales. Después MVP completo.
