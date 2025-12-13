"use client"

import { toaster } from "@/components/ui/toaster"

let lastToastKey: string | null = null
let lastToastAt = 0

const shouldSkip = (key: string, windowMs: number) => {
  return lastToastKey === key && Date.now() - lastToastAt < windowMs
}

const record = (key: string) => {
  lastToastKey = key
  lastToastAt = Date.now()
}

const useCustomToast = () => {
  const showSuccessToast = (description: string) => {
    const key = `success:${description}`
    if (shouldSkip(key, 15000)) return
    toaster.create({
      title: "Success!",
      description,
      type: "success",
    })
    record(key)
  }

  const showErrorToast = (description: string) => {
    const key = `error:${description}`
    if (shouldSkip(key, 15000)) return
    toaster.create({
      title: "Something went wrong!",
      description,
      type: "error",
      duration: 15000,
    })
    record(key)
  }

  return { showSuccessToast, showErrorToast }
}

export default useCustomToast
