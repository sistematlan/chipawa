package i18n

// spanish catalog. Should mirror english keys, but missing entries
// fall back gracefully to English via the lookup logic in T().
var spanish = map[string]string{
	// ---- Risk labels ----
	"risk.safe":      "seguro (caché regenerable)",
	"risk.ask":       "preguntar antes (puede contener datos del usuario)",
	"risk.dangerous": "peligroso (irreversible)",

	// ---- Caches: npm ----
	"caches.npm.name":             "Caché de npm",
	"caches.npm.detail.advanced":  "paquetes descargados",
	"caches.npm.detail.simple":    "Paquetes de Node.js descargados",

	"caches.npm-npx.name":            "Caché de npx",
	"caches.npm-npx.detail.advanced": "ejecuciones one-shot de npx",
	"caches.npm-npx.detail.simple":   "Herramientas temporales ejecutadas con npx",

	"caches.npm-logs.name":            "Logs de npm",
	"caches.npm-logs.detail.advanced": "logs viejos de instalación",
	"caches.npm-logs.detail.simple":   "Logs viejos de instalaciones de npm",

	// ---- Caches: pnpm ----
	"caches.pnpm.name":            "Almacén de pnpm",
	"caches.pnpm.detail.advanced": "almacén global content-addressable",
	"caches.pnpm.detail.simple":   "Caché compartido de paquetes pnpm",

	// ---- Caches: yarn ----
	"caches.yarn.name":            "Caché de yarn",
	"caches.yarn.detail.advanced": "paquetes descargados",
	"caches.yarn.detail.simple":   "Paquetes de Yarn descargados",

	// ---- Caches: brew ----
	"caches.brew.name":            "Caché de Homebrew",
	"caches.brew.detail.advanced": "bottles y fuentes descargadas",
	"caches.brew.detail.simple":   "Archivos viejos de instalación de Homebrew",

	// ---- Caches: jetbrains ----
	"caches.jetbrains.name":            "Caché de JetBrains",
	"caches.jetbrains.detail.advanced": "índices y logs",
	"caches.jetbrains.detail.simple":   "Índices y cachés de IDEs JetBrains (PhpStorm, IntelliJ…)",

	"caches.jetbrains-old.detail.simple":   "Versión antigua de %s — la actual es %s",
	"caches.jetbrains-old.detail.advanced": "versión antigua de %s (última: %s)",

	// ---- Caches: go ----
	"caches.go.name":            "Caché de compilación Go",
	"caches.go.detail.advanced": "caché de compilación",
	"caches.go.detail.simple":   "Caché de compilación de Go",

	// ---- Caches: pip / uv / composer / node-gyp ----
	"caches.pip.name":            "Caché de pip",
	"caches.pip.detail.advanced": "wheels y caché HTTP",
	"caches.pip.detail.simple":   "Caché de paquetes de Python",

	"caches.uv.name":            "Caché de uv",
	"caches.uv.detail.advanced": "caché de paquetes Python",
	"caches.uv.detail.simple":   "Caché del gestor de paquetes uv",

	"caches.composer.name":            "Caché de Composer",
	"caches.composer.detail.advanced": "paquetes PHP",
	"caches.composer.detail.simple":   "Caché de paquetes PHP",

	"caches.node-gyp.name":            "Caché de node-gyp",
	"caches.node-gyp.detail.advanced": "headers nativos para compilar",
	"caches.node-gyp.detail.simple":   "Archivos para compilar módulos nativos de Node.js",

	// ---- Caches: browsers ----
	"caches.chrome.name":            "Caché de Chrome",
	"caches.chrome.detail.advanced": "caché del navegador",
	"caches.chrome.detail.simple":   "Archivos temporales de Google Chrome",

	"caches.firefox.name":            "Caché de Firefox",
	"caches.firefox.detail.advanced": "caché del navegador",
	"caches.firefox.detail.simple":   "Archivos temporales de Firefox",

	// ---- Caches: cargo ----
	"caches.cargo-cache.name":            "Caché de Cargo (registry)",
	"caches.cargo-cache.detail.advanced": "crates descargados",
	"caches.cargo-cache.detail.simple":   "Paquetes de Rust descargados",

	"caches.cargo-src.name":            "Fuentes de Cargo (registry)",
	"caches.cargo-src.detail.advanced": "fuentes de crates extraídas",
	"caches.cargo-src.detail.simple":   "Fuentes extraídas de paquetes Rust",

	"caches.cargo-git.name":            "Cargo git checkouts",
	"caches.cargo-git.detail.advanced": "dependencias por git",
	"caches.cargo-git.detail.simple":   "Dependencias de Rust desde Git",

	// ---- Caches: xcode ----
	"caches.xcode-derived.name":            "Xcode DerivedData",
	"caches.xcode-derived.detail.advanced": "artefactos de build",
	"caches.xcode-derived.detail.simple":   "Archivos de compilación de Xcode (se regeneran al compilar)",

	"caches.xcode-archives.name":            "Archivos Xcode",
	"caches.xcode-archives.detail.advanced": "archivos de release antiguos",
	"caches.xcode-archives.detail.simple":   "Archivos viejos de apps Xcode",

	"caches.xcode-ios-support.name":            "iOS DeviceSupport",
	"caches.xcode-ios-support.detail.advanced": "símbolos para versiones de iOS antiguas",
	"caches.xcode-ios-support.detail.simple":   "Símbolos para versiones antiguas de iOS",

	"caches.xcode-simulator.name":            "Cachés de simuladores",
	"caches.xcode-simulator.detail.advanced": "cachés de los simuladores",
	"caches.xcode-simulator.detail.simple":   "Caché de los simuladores de iOS",

	// ---- Caches: docker ----
	"caches.docker.name":            "Docker recuperable",
	"caches.docker.detail.advanced": "imágenes, caché de build, contenedores parados (excluye volúmenes)",
	"caches.docker.detail.simple":   "Imágenes Docker viejas y contenedores parados",

	// ---- Orphans ----
	"orphans.docker-leftover.name":            "Restos de Docker Desktop",
	"orphans.docker-leftover.detail.advanced": "Docker.app fue desinstalada pero sus datos siguen ocupando espacio",
	"orphans.docker-leftover.detail.simple":   "Datos de Docker que quedaron tras desinstalar la app",

	"orphans.whatsapp-media.name":            "Media de WhatsApp",
	"orphans.whatsapp-media.detail.advanced": "fotos/videos/audio descargados (los chats no se ven afectados)",
	"orphans.whatsapp-media.detail.simple":   "Fotos, videos y audios descargados por WhatsApp (los chats se conservan)",

	// ---- Downloads subcategories ----
	"downloads.installer.detail.simple":            "Instalador; %s ya está instalada",
	"downloads.installer.detail.advanced":          "instalador; %s ya instalada",
	"downloads.archive-extracted.detail.simple":    "Archivo comprimido ya extraído en ./%[1]s/",
	"downloads.archive-extracted.detail.advanced":  "archivo %[2]s ya extraído en ./%[1]s/",
	"downloads.project-folder.detail.simple":       "Carpeta de proyecto con node_modules (probablemente abandonada)",
	"downloads.project-folder.detail.advanced":     "carpeta de proyecto con node_modules/target/dist (probablemente abandonada)",
	"downloads.db-dump.detail.simple":              "Backup viejo de base de datos (>30 días)",
	"downloads.db-dump.detail.advanced":            "dump de base de datos (>30 días)",
	"downloads.old-video.detail.simple":            "Video viejo (>90 días)",
	"downloads.old-video.detail.advanced":          "video (>90 días)",
	"downloads.old-archive.detail.simple":          "Archivo comprimido viejo (>90 días)",
	"downloads.old-archive.detail.advanced":        "archivo comprimido (>90 días)",
	"downloads.large-other.detail.simple":          "Archivo grande — revisa antes de borrar",
	"downloads.large-other.detail.advanced":        "archivo grande sin clasificar — revisa antes de borrar",

	// ---- Categories ----
	"category.cache.simple":     "Archivos temporales de herramientas de desarrollo",
	"category.cache.advanced":   "Cachés de desarrollo",
	"category.orphan.simple":    "Datos sobrantes de apps desinstaladas",
	"category.orphan.advanced":  "Datos huérfanos",
	"category.download.simple":  "Archivos en tu carpeta de Descargas",
	"category.download.advanced": "Descargas",

	// ---- Common UI strings ----
	"ui.size":         "tamaño",
	"ui.tool":         "herramienta",
	"ui.detail":       "detalle",
	"ui.note":         "nota",
	"ui.path":         "ruta",
	"ui.risk":         "riesgo",
	"ui.age":          "edad",
	"ui.file":         "archivo",
	"ui.today":        "hoy",
	"ui.days-ago":     "hace %d días",
	"ui.never-used":   "nunca",
	"ui.empty":        "—",
	"ui.total":        "Total",
	"ui.recoverable":  "recuperables",
	"ui.requires-confirmation": "requiere confirmación",
	"ui.nothing":      "Nada que limpiar. Disco en orden.",
	"ui.dry-run-mode": "Modo: dry-run (no se borrará nada)",

	// ---- Cleaner prompts ----
	"cleaner.prompt":              "[s/N/v=ver/q=salir] ",
	"cleaner.prompt.dangerous":    "ESCRIBE el nombre exacto (%q) o vacío para cancelar:\n  > ",
	"cleaner.removing":            "borrando %s ... ",
	"cleaner.ok":                  "ok",
	"cleaner.failed":              "FALLÓ: %v",
	"cleaner.dry-prefix":          "[dry-run] se borraría: ",
	"cleaner.summary":             "Resumen: %d eliminados, %d omitidos, %d fallidos",
	"cleaner.freed":               "Espacio liberado: %s (de %s planificados)",
	"cleaner.plan":                "Plan: %d ítems, %s recuperables",
	"cleaner.user-quit":           "el usuario salió",
	"cleaner.user-declined":       "el usuario declinó",

	// ---- Wizard ----
	"wizard.tagline":             "Limpia tu Mac, al estilo open-source.",
	"wizard.scanning":            "Escaneando tu disco...",
	"wizard.menu.header":         "¿Qué tipo de limpieza quieres hacer?",
	"wizard.menu.prompt":         "Opción [1-4]: ",
	"wizard.menu.cancel":         "Cancelar — no borrar nada",
	"wizard.cancelled":           "Cancelado. No se borró nada.",
	"wizard.thanks":              "Listo. Gracias por usar mistah.",
	"wizard.level.light.name":    "Ligera",
	"wizard.level.light.desc":    "Solo cachés seguras de dev (npm, brew, pip…)",
	"wizard.level.standard.name": "Estándar",
	"wizard.level.standard.desc": "Ligera + Docker + JetBrains + restos de Xcode",
	"wizard.level.deep.name":     "Profunda",
	"wizard.level.deep.desc":     "Estándar + datos huérfanos + candidatos en Downloads",
	"wizard.confirm.level":       "Nivel",
	"wizard.confirm.items":       "ítems",
	"wizard.confirm.prompt":      "¿Continuar? [s/N] ",
}
