"use client"

import { ChakraProvider } from "@chakra-ui/react"
import React, { type PropsWithChildren, useEffect } from "react"
import useCustomToast from "../../hooks/useCustomToast"
import { system } from "../../theme"
import { ColorModeProvider } from "./color-mode"
import { Toaster } from "./toaster"

export function CustomProvider(props: PropsWithChildren) {
  const { showErrorToast, showSuccessToast } = useCustomToast()

  useEffect(() => {
    try {
      const raw = localStorage.getItem("persisted_toast")
      if (raw) {
        const data = JSON.parse(raw) as { type?: string; description?: string }
        if (data?.description) {
          console.log(data.description)
          if (data.type === "success") showSuccessToast(data.description)
          else showErrorToast(data.description)
        }
        localStorage.removeItem("persisted_toast")
      }
    } catch {}
  }, [showErrorToast, showSuccessToast])

  return (
    <ChakraProvider value={system}>
      <ColorModeProvider defaultTheme="light">
        {props.children}
      </ColorModeProvider>
      <Toaster />
    </ChakraProvider>
  )
}
