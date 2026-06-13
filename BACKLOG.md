# BACKLOG — mistah

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
- [x] Comando `clean` con `--dry-run`, `--yes`, `--include-orphans`, `--include-downloads`
- [x] Bug fix: scanner reutilizado en `TerminalPrompter` (no se descarta el buffer entre prompts)
- [x] Tests unitarios (10 casos en cleaner, parser de tamaños Docker, splitJBVersion, fixtures de orphans)
- [x] **`internal/downloads`** con 7 clasificadores: installer-with-app, archive-extracted, project-folder, db-dump, old-video, old-archive, large-other
- [x] **Comando `downloads`** con tabla agrupada por subcategoría
- [x] Integración downloads → clean vía `--include-downloads`
- [x] Tests de downloads con fixtures temporales (8 casos)
- [x] **`internal/i18n`** sin libs externas, autodetect $LANG/$LC_ALL, soporte es+en, fallback a inglés
- [x] **Catálogos es/en** con paridad garantizada por tests, variantes `.simple` y `.advanced` por entrada
- [x] **`Item.HumanName()`, `HumanDetail(simple)`, `HumanRisk()`** para presentación localizada
- [x] **Flags globales `--lang` y `--advanced`** en root.go
- [x] **Migración de todos los detectores** a NameKey/DetailKey/DetailArgs
- [x] **TerminalPrompter localizado**: prompts, summary, confirmación dangerous
- [x] **`internal/spinner`** con detección de TTY y fallback para non-TTY
- [x] **`internal/wizard`** con tres niveles preset (Light/Standard/Deep) y filtros por nivel
- [x] **Wizard UI completa**: intro, scan animado, menú numerado, confirmación, resumen
- [x] **Cobra integration**: `mistah` sin args ejecuta wizard. Subcomandos siguen igual.
- [x] **Tests del wizard**: 6 casos cubriendo niveles, totales, monotonicidad, level.String()

---

## 🔥 Próxima sesión — completar MVP

### 1. Comando `downloads` (P0) — ✅ COMPLETADO

- [x] `internal/downloads/downloads.go` con detector
- [x] Smart rules: installer-with-app, archive-extracted, project-folder, db-dump, old-video, old-archive
- [x] Catch-all `large-other` para archivos >100 MB sin clasificar
- [x] Comando `cmd/downloads.go` con tabla agrupada
- [x] Integrar items en `clean` (categoría `CategoryDownload`, `RiskAskBefore`)
- [x] Tests con fixtures temporales

### 2. Flags globales + `report` (P0)

- [ ] Flag global `--path` (ya existe parcial en `projects`, normalizar)
- [ ] Flag global `--json` — toda salida puede emitirse como JSON
- [ ] Flag global `--no-color` — desactiva ANSI (cuando se agregue color)
- [ ] Flag global `--yes` — heredable (hoy solo está en `clean`)
- [ ] Comando `report` con `--out <file>` y formato estructurado
- [ ] Refactor: que cada `cmd/*.go` use un renderer común (text vs JSON)

**Estimado:** 2-3 horas. Refactor mediano porque toca todos los comandos.

### 3. Wizard mejoras (post-MVP, opcional)

- [ ] Flags directos para no-tty: `mistah --light`, `--standard`, `--deep` que skipean menu
- [ ] `--no-confirm` para modo CI con wizard
- [ ] Mostrar breakdown por categoría dentro del nivel ("npm: 1.2 GB, brew: 200 MB...")
- [ ] Spinner multi-fase: "Scanning caches..." → "Scanning orphans..." → "Scanning downloads..."

### 4. Color y UI polish (P1)

- [ ] Lib: `github.com/fatih/color` o solo ANSI manual (preferible — menos deps)
- [ ] Verde para "safe", amarillo para "ask", rojo para "danger"
- [ ] Detección de TTY: solo color cuando stdout sea terminal
- [ ] Símbolos: ✓ ✗ → ⚠ (con fallback ASCII en `--no-color`)
- [ ] Tabla con anchos calculados dinámicamente (hoy son fijos y truncan feo)

---

## 🟡 Backlog medio plazo (post-MVP)

### Más detectores de caches

#### P0 — alto impacto, fáciles de implementar
- [ ] **Electron caches genéricos**: `~/Library/Application Support/*/Cache`, `*/Code Cache`, `*/GPUCache` — culpable silencioso (VSCode, Cursor, Trae, Antigravity, Comet, Slack, Discord, Teams, Notion, Linear, todos los que sean Electron)
- [ ] **iOS Simulator devices**: `~/Library/Developer/CoreSimulator/Devices/` — los simuladores en sí mismos pueden ocupar 5-50 GB en devs iOS
- [ ] **Time Machine local snapshots**: `tmutil listlocalsnapshots /` — Apple los crea sin avisar, ocupan "purgeable space"
- [ ] **Mail.app adjuntos**: `~/Library/Mail/V*/MailData/`
- [ ] **Messages adjuntos**: `~/Library/Messages/Attachments/`
- [ ] **Photos library** caches: `~/Pictures/Photos Library.photoslibrary/resources/`
- [ ] **Spotify cache**: `~/Library/Caches/com.spotify.client`
- [ ] **SiriTTS voices**: `~/Library/Caches/SiriTTS` (voces descargadas, opcional desinstalar via Settings)
- [ ] **VSCode/Cursor/Trae específicos**: `~/Library/Application Support/{Code,Cursor,Trae,Antigravity}/Cache*` con desglose

#### P1 — dev tools faltantes
- [ ] **Bun**: `~/.bun/install/cache`
- [ ] **Deno**: `~/Library/Caches/deno`
- [ ] **Gradle**: `~/.gradle/caches`
- [ ] **Maven**: `~/.m2/repository` (cuidado: no es solo caché, también deps activas)
- [ ] **Playwright**: `~/Library/Caches/ms-playwright`
- [ ] **Cypress**: `~/Library/Caches/Cypress`
- [ ] **Puppeteer**: `~/.cache/puppeteer`
- [ ] **rustup toolchains viejas**: parsear `rustup toolchain list`
- [ ] **nvm versions viejas**: `~/.nvm/versions/node/v*`
- [ ] **pnpm node versions**: `~/Library/pnpm/nodejs/*`
- [ ] **Volta**: `~/.volta/cache`
- [ ] **VirtualBox VMs huérfanas**: `~/VirtualBox VMs`
- [ ] **OrbStack VMs**: desglosar Group Container por VM
- [ ] **JDownloader 2**: `~/Library/Application Support/JDownloader 2`
- [ ] **TorBrowser-Data**: `~/Library/Application Support/TorBrowser-Data`

### Más detectores de orphans

- [ ] **Group Containers de apps desinstaladas**: cruzar `~/Library/Group Containers/*` con `/Applications`
- [ ] **Application Support de apps desinstaladas**: cruzar con `/Applications`
- [ ] **Containers de apps desinstaladas**: cruzar `~/Library/Containers/*` con `/Applications`
- [ ] **Wondershare leftover**: `~/Library/Application Support/com.wondershare.Installer`
- [ ] **Mega Limited cache** (la app sigue, pero el cache crece sin tope)
- [ ] **Adobe leftover** (típicamente >5 GB tras desinstalar)

### Detector genérico "big files" (catch-all)

- [ ] Comando `mistah big-files` que escanea `$HOME` y reporta los N archivos más grandes
- [ ] Filtros: `--min-size 100M`, `--ext .mov,.mp4,.dmg`, `--older-than 90d`
- [ ] Clasifica automáticamente:
  - **Videos personales** (.mov, .mp4, .mkv) >500 MB
  - **DMGs/PKGs/ISOs** sueltos en cualquier sitio
  - **ZIPs/RARs/7z** + detectar si están extraídos al lado
  - **Backups** (.bak, .sql, .dump, .tar.gz)
  - **Carpetas con node_modules** fuera de `~/sourcecode`
  - **Binarios de build** (`target/`, `dist/`, `build/`, `.next/`) en cualquier proyecto
  - **Logs viejos** `~/Library/Logs/*` con >90 días
  - **Crash reports** `~/Library/Logs/DiagnosticReports/`
- [ ] Opcional: detección de duplicados por hash (SHA-256 sobre archivos >100 MB)

### `projects` mejorado

- [ ] Detectar `node_modules`, `vendor/`, `target/`, `.next/`, `dist/`, `build/` dentro de cada proyecto
- [ ] Reportar tamaño de artefactos vs código real
- [ ] Detectar proyectos sin commit en >1 año como "abandonados"
- [ ] Sugerir `git clean -xdf` para los abandonados con git
- [ ] Modo `--clean-build-artifacts` que borra solo carpetas regenerables

### `apps` mejorado

- [ ] Tamaño TOTAL (app + Application Support + Containers + Caches asociados)
- [ ] Sugerir desinstalación con `mdls` last-used > 180d
- [ ] Comando `mistah apps uninstall <name>` que llama a `mdfind` para encontrar todos los archivos relacionados

### Mejoras al cleaner

- [ ] Modo `--category cache|orphan|download|all` (filtro por categoría)
- [ ] Modo `--tool docker|jetbrains|...` (filtro por tool)
- [ ] Modo `--min-size 100M` (solo items grandes)
- [ ] Confirmación batch: "borrar todos los marcados con [s]" al final
- [ ] Undo log: registrar qué se borró y cuándo en `~/.mistah/history.log`

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
- [ ] Post en HackerNews ("Show HN: mistah — a transparent disk cleaner for macOS developers")
- [ ] Reddit: r/macOS, r/golang, r/programming
- [ ] Demo en video (asciinema) para README
- [ ] Documentación de cada detector con ejemplo de output

### Ideas más ambiciosas

- [ ] **Modo daemon**: ejecuta `scan` periódicamente y notifica cuando disco <20%
- [ ] **Plugins**: detectores externos cargados desde `~/.mistah/plugins/*.so` o subprocess
- [ ] **Dashboard TUI** con bubbletea (alternativa al CLI puro)
- [ ] **Política de retención**: archivo de config (`~/.mistah/config.toml`) con reglas tipo `npm cache: keep 30 days`
- [ ] **Comparativa pre/post**: muestra diff de uso por carpeta tras `clean`
- [ ] **Sugerencias inteligentes**: "tienes 12 carpetas con node_modules y solo trabajas en 3 proyectos activos"

### Multiplataforma — análisis técnico

**Estado actual:** mistah es macOS-only por diseño. Hay 18 referencias a `~/Library/`
y todos los detectores asumen filesystem Unix con paths macOS específicos.

#### Soporte Linux (medio plazo)
- [ ] Más fácil que Windows: solo cambian paths, no APIs.
- [ ] `~/.cache/`, `~/.config/`, `~/.local/share/` en lugar de `~/Library/*`
- [ ] Apps desde `~/.local/share/applications/` y `/usr/share/applications/`
- [ ] No hay `mdls` — usar `stat -c %Y` para "last modified" como proxy de uso
- [ ] Snap, Flatpak, AppImage como sources de caches
- [ ] **Esfuerzo estimado:** 1-2 semanas con CI Linux runner

#### Soporte Windows (lejos)
- [ ] **NO portar antes de v1.0 macOS estable.** Razones:
  - Target del SPEC es devs macOS; ahí está el dolor real.
  - Competencia fuerte y nativa: WizTree, TreeSize, BleachBit, Storage Sense.
  - Devs Windows usan WSL2 → caches viven en Linux, problema diferente.
  - 3-4 semanas de trabajo que no aporta a usuarios actuales.
  - Cada plataforma multiplica costo de mantenimiento ~2.5x.
- [ ] **Si se hace alguna vez:**
  - [ ] Abstraer `PathProvider` con builds tags (`//go:build darwin`, `windows`, `linux`)
  - [ ] Reemplazar `du -sk` por `filepath.Walk` nativo (también beneficia a macOS)
  - [ ] Detector de apps via Windows registry (`HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`)
  - [ ] Caches: `%LOCALAPPDATA%`, `%APPDATA%`, `%TEMP%`
  - [ ] Docker Desktop en Windows: archivos `.vhdx` en `%LOCALAPPDATA%\Docker\wsl\`
  - [ ] Manejar archivos bloqueados (Windows: handle locking, antivirus)
  - [ ] `SafeRoots` portable: `%USERPROFILE%`, `%TEMP%`
  - [ ] Reemplazar prompts ANSI con detección de Windows Terminal vs cmd.exe
  - [ ] CI con `windows-latest` runner en GitHub Actions
  - [ ] Tests funcionales en máquina Windows real
  - [ ] **Esfuerzo estimado:** 3-4 semanas con acceso a Windows real

### Métricas de salud

- [ ] `mistah doctor`: detecta problemas comunes (Docker sin volumes, brew obsoleto, etc.)
- [ ] `mistah baseline`: guarda snapshot del estado actual para comparar después
- [ ] `mistah diff <baseline>`: muestra qué creció/decreció desde un baseline

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
- **`clean` por defecto**: solo caches; orphans requieren `--include-orphans`, downloads requieren `--include-downloads`
- **UX de confirmación**: ítem por ítem con `[s/N/v/q]`, default seguro en NO
- **Docker**: `system prune -f` (sin `--volumes`); volumes solo con flag explícito futuro `--include-volumes`
- **SafeRoots**: solo `$HOME`, `/var/folders`, `/tmp` y `/private/*`. Cualquier otro path es rechazado.
- **Multiplataforma**: macOS only hasta v1.0. Linux después. Windows solo si hay demanda real (no especulativa).
- **Cobertura objetivo v0.1.0**: caches dev + orphans básicos + downloads. P0 de Application Support (Electron, iOS Simulator, Time Machine, Mail, Messages) queda para v0.2.0.

### Democratización (sesión 4 backlog)

- **Audiencia objetivo expandida**: técnicos no-dev y power users macOS, no consumer mainstream.
- **NO hacer GUI nativa Swift**: convertiría mistah en otro producto, mercado saturado por CleanMyMac, perdería ventaja de auditabilidad open-source.
- **NO hacer TUI bubbletea**: gain marginal sobre Camino A para audiencia adicional.
- **Modo por default**: simple (lenguaje humano, niveles preset). Modo técnico opt-in con `--advanced` o `--verbose`.
- **Wizard sin args**: `mistah` sin subcomando ejecuta wizard con tres niveles preset (Ligera / Estándar / Profunda).
- **Distribución**: Apple notarization + curl one-liner + Homebrew. Requiere Apple Developer Account ($99/año).
- **i18n**: detección automática via `LANG`/`LC_ALL`. Soporte es + en. Sin libs externas — mapas Go nativos.
- **Plan de implementación**: 4 sesiones — (1) i18n+lenguaje humano, (2) wizard+niveles, (3) report+flags globales, (4) distribución+notarización.

### User research / feedback (estrategia)

- **No A/B testing tradicional hasta tener 1000+ usuarios reales.** Requiere telemetría que traiciona la promesa de privacidad.
- **Etapa 1 (primeros 5 usuarios)**: sesiones guiadas 1-on-1 de 30 minutos vía Zoom + share screen. Cualitativo profundo. Jakob Nielsen: 5 usuarios descubren 80% de problemas de UX.
- **Etapa 2 (50-100 usuarios)**: encuesta opt-in al final de `clean` con link manual a Tally/Typeform. URL solo lleva versión + bytes liberados, nunca paths o contenido. Sin recolección automática.
- **Etapa 3 (1000+ usuarios)**: solo entonces considerar telemetría opt-in granular con default OFF, comando explícito (`mistah telemetry enable`), documentación de payload exacto y endpoint open-source para auditar.
- **Reclutamiento de testers iniciales**:
  - [ ] Build in public: tweet/post en X y LinkedIn buscando 5 beta testers macOS
  - [ ] r/macapps, r/golang
  - [ ] Slack/Discord de comunidades técnicas locales
  - [ ] Compañeros devs directos
- **NUNCA**:
  - Telemetría sin consentimiento explícito
  - Recolectar paths, nombres de archivo, contenido
  - Servidor de datos sin código fuente abierto
  - Cookie tracking en landing page

---

## 🎯 Definition of Done para v0.1.0

Del SPEC original:

- [x] `mistah scan` muestra resumen de disco y top categorías
- [x] `mistah apps` lista apps con último uso y tamaño
- [x] `mistah caches` detecta y totaliza caches de dev
- [x] `mistah projects --path ~/sourcecode` reporta estado git de cada repo
- [x] `mistah clean --dry-run` lista candidatos sin borrar nada
- [x] `mistah clean` pide confirmación por ítem
- [x] `mistah downloads` clasifica Downloads por tipo y antigüedad
- [ ] `mistah report --json` emite reporte estructurado
- [x] Binario único, sin dependencias externas en runtime
- [x] `make build` compila en < 10 segundos
- [x] Tests pasan con `make test`

**Faltan:** `report --json` y flags globales. Después MVP completo (10/11 ✓).
