# mistah — spec

> *mistah* (maya yucateco): "él/ella barre", del verbo *mis* (barrer).
> CLI open-source para macOS que analiza el disco, detecta basura de desarrollo y libera espacio con confirmación explícita antes de tocar nada.

**Estado:** borrador v0.1  
**Stack:** Go 1.24+  
**Distribución objetivo:** Homebrew tap (`sistematlan/tools`)

---

## 1. Objetivo

Reemplazar el flujo manual de limpieza que un desarrollador macOS repite cada cierto tiempo:
revisar apps sin uso, limpiar caches de herramientas de desarrollo, identificar proyectos
duplicados o abandonados, y liberar espacio de Docker — todo desde la terminal, sin GUI,
sin suscripción, sin telemetría.

**Usuario objetivo:** desarrollador macOS con herramientas como Docker, Homebrew, Node,
Go, JetBrains, editores Electron (Cursor, VS Code, etc.).

---

## 2. Comandos

```
mistah scan          Escaneo completo: disco, apps, caches, proyectos
mistah apps          Apps instaladas con fecha de último uso y tamaño
mistah caches        Caches de herramientas dev (Docker, npm, Homebrew, etc.)
mistah projects      Análisis de carpetas de código fuente
mistah downloads     Archivos en ~/Downloads agrupados por tipo y antigüedad
mistah clean         Limpieza interactiva con confirmación ítem por ítem
mistah clean --dry-run  Muestra qué se eliminaría sin borrar nada
mistah report        Genera reporte en JSON o texto plano
mistah version       Versión del binario
mistah help          Ayuda general
```

### Flags globales

```
--path <dir>          Directorio raíz para análisis (default: ~)
--json                Output en JSON (para scripting)
--yes                 Confirmar todo automáticamente (modo no-interactivo)
--no-color            Sin colores ANSI
```

---

## 3. Módulos internos

```
mistah/
├── cmd/                  Comandos CLI (cobra)
│   ├── root.go
│   ├── scan.go
│   ├── apps.go
│   ├── caches.go
│   ├── projects.go
│   ├── downloads.go
│   ├── clean.go
│   └── report.go
├── internal/
│   ├── disk/             Uso de disco, df/du equivalente
│   ├── apps/             Lectura de /Applications, mdls last-used
│   ├── caches/           Detectores por herramienta (docker, npm, brew, etc.)
│   ├── projects/         Análisis de repos git (remote, antigüedad, tamaño)
│   ├── downloads/        Clasificación de ~/Downloads
│   ├── cleaner/          Lógica de borrado con confirmación
│   └── report/           Formateo de salida (texto, JSON)
├── pkg/
│   └── units/            Formateo de bytes (MB, GB)
├── SPEC.md
├── README.md
├── go.mod
└── Makefile
```

---

## 4. Comportamiento por módulo

### `scan`
- Muestra resumen de disco (usado / libre / total)
- Llama a los demás módulos en secuencia
- Output agrupado por categoría con totales

### `apps`
- Lee `/Applications/*.app` y `~/Applications/*.app`
- Usa `mdls -name kMDItemLastUsedDate` para fecha de último uso
- Clasifica: activa (< 30 días), inactiva (30–90 días), sin uso (> 90 días o null)
- Muestra tamaño de cada app

### `caches`
Detectores individuales para:
| Herramienta | Ruta |
|---|---|
| Docker images/build cache | `docker system df` |
| npm | `~/.npm` |
| Homebrew | `~/Library/Caches/Homebrew` |
| JetBrains | `~/Library/Caches/JetBrains` |
| Go build | `~/Library/Caches/go-build` |
| pip | `~/Library/Caches/pip` |
| Composer | `~/Library/Caches/composer` |
| yarn | `~/.yarn/cache` |
| node-gyp | `~/Library/Caches/node-gyp` |
| Xcode DerivedData | `~/Library/Developer/Xcode/DerivedData` |
| iOS Simulators | `~/Library/Developer/CoreSimulator/Caches` |
| Browsers (Chrome, Firefox) | `~/Library/Caches/Google`, `Mozilla` |

### `projects`
- Recibe `--path` (default: `~/sourcecode`)
- Por cada subdirectorio:
  - ¿Tiene `.git`? → sí/no
  - Si sí: remote URL, rama activa, fecha último commit, número de commits
  - Si no: fecha de última modificación
  - Tamaño de la carpeta
- Clasifica: activo (< 90 días), inactivo (90–365 días), abandonado (> 1 año o sin git)

### `downloads`
- Agrupa por extensión y antigüedad (< 30d, 30–90d, > 90d)
- Detecta ZIPs que tienen carpeta extraída con mismo nombre al lado

### `clean`
- Presenta lista de ítems candidatos a eliminar con tamaño
- Para cada ítem: `[s/n/ver]` — sí, no, o ver contenido antes
- `--dry-run` no borra nada, solo reporta
- `--yes` confirma todo (para scripts CI/automatización)
- Nunca borra Docker volumes sin flag explícito `--include-volumes`
- Nunca borra archivos fuera de rutas conocidas y seguras

### `report`
- `--json`: output estructurado para pipes o scripts
- `--out <file>`: guarda el reporte en archivo

---

## 5. Estilo de código

- Go idiomático: errores explícitos, sin panic en flujo normal
- Un paquete por responsabilidad, sin dependencias circulares
- Interfaces para los detectores de caché → fácil agregar nuevas herramientas
- Sin globals mutables
- Nombres en inglés en el código, mensajes de UI en español (configurable después)
- Archivos entre 200–400 líneas, máx 800

---

## 6. Testing

- Unit tests para cada detector de caché y clasificador
- Tests de integración con fixtures (directorios temporales simulados)
- `--dry-run` siempre debe tener cobertura de test
- Sin tests que borren archivos reales del sistema
- `make test` corre todo, `make test-unit` solo unitarios

---

## 7. Distribución

```
# Fase 1 — binario manual
make build → ./bin/mistah

# Fase 2 — Homebrew tap
brew tap sistematlan/tools
brew install mistah

# Fase 3 — GitHub Releases con goreleaser
goreleaser release
```

---

## 8. Límites (boundaries)

| Categoría | Regla |
|---|---|
| **Siempre hacer** | Mostrar tamaño antes de borrar. Pedir confirmación. Dry-run disponible. |
| **Preguntar antes** | Docker volumes. Carpetas de proyecto sin git. Apps del sistema. |
| **Nunca hacer** | Borrar fuera de rutas conocidas. Ejecutar como root. Enviar telemetría. Borrar sin confirmación explícita (salvo `--yes`). |

---

## 9. MVP — criterios de aceptación

- [ ] `mistah scan` muestra resumen de disco y top categorías
- [ ] `mistah apps` lista apps con último uso y tamaño
- [ ] `mistah caches` detecta y totaliza caches de dev
- [ ] `mistah projects --path ~/sourcecode` reporta estado git de cada repo
- [ ] `mistah clean --dry-run` lista candidatos sin borrar nada
- [ ] `mistah clean` pide confirmación por ítem
- [ ] Binario único, sin dependencias externas en runtime
- [ ] `make build` compila en < 10 segundos
- [ ] Tests pasan con `make test`
