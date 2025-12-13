// scripts/generate-openapi.js
import { execSync } from "node:child_process"
import { readdirSync } from "node:fs"
import { resolve } from "node:path"

const OPENAPI_DIR = resolve("openapi")
const files = readdirSync(OPENAPI_DIR).filter((f) => f.endsWith(".json"))

if (files.length === 0) {
  console.error("❌ No OpenAPI files found in ./openapi")
  process.exit(1)
}

console.log(`🔍 Found OpenAPI specs:\n - ${files.join("\n - ")}\n`)

for (const file of files) {
  const input = resolve("openapi", file)
  const name = file.replace("_service.json", "").replace(".json", "")
  const output = resolve("src/client", name)

  console.log(`🚀 Generating client for ${file} → ${output}`)

  // Run openapi-ts via CLI
  execSync(
    `openapi-ts --client legacy/axios --input "${input}" --output "${output}"`,
    {
      stdio: "inherit",
    },
  )
}

console.log("\n✨ All clients generated!")
