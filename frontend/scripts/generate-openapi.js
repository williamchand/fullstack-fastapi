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
console.log(
  "\uD83D\uDD0D Found OpenAPI specs:\n - ".concat(files.join("\n - "), "\n"),
)
for (let _i = 0, files_1 = files; _i < files_1.length; _i++) {
  const file = files_1[_i]
  const input = resolve("openapi", file)
  const name_1 = file.replace("_service.json", "").replace(".json", "")
  const output = resolve("src/client", name_1)
  console.log(
    "\uD83D\uDE80 Generating client for "
      .concat(file, " \u2192 ")
      .concat(output),
  )
  // Run openapi-ts via CLI
  execSync(
    'openapi-ts --client legacy/axios --input "'
      .concat(input, '" --output "')
      .concat(output, '"'),
    {
      stdio: "inherit",
    },
  )
}
console.log("\n✨ All clients generated!")
