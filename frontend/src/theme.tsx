import { createSystem, defaultConfig } from "@chakra-ui/react"
import { buttonRecipe } from "./theme/button.recipe"

export const system = createSystem(defaultConfig, {
  globalCss: {
    html: {
      fontSize: "16px",
    },
    body: {
      fontSize: "0.875rem",
      margin: 0,
      padding: 0,
      backgroundColor: "#d8ffdd",
      color: "#524632",
    },
    ".main-link": {
      color: "ui.main",
      fontWeight: "bold",
    },
  },
  theme: {
    tokens: {
      colors: {
        ui: {
          main: { value: "#8f7e4f" },
        },
        brand: {
          darkKhaki: { value: "#524632" },
          fadedCopper: { value: "#8f7e4f" },
          drySage: { value: "#c3c49e" },
          frostedMint: { value: "#d8ffdd" },
          dustGrey: { value: "#dedbd8" },
        },
      },
    },
    recipes: {
      button: buttonRecipe,
    },
  },
})
