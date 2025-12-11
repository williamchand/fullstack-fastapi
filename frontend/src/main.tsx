import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { RouterProvider, createRouter } from "@tanstack/react-router"
import { StrictMode } from "react"
import ReactDOM from "react-dom/client"
import { routeTree } from "./routeTree.gen"

import { ApiError, OpenAPI } from "./client/user"
import axios from "axios"
import { CustomProvider } from "./components/ui/provider"

OpenAPI.BASE = import.meta.env.VITE_API_URL
// Debug: confirm API base at startup
console.debug("OpenAPI.BASE:", OpenAPI.BASE)
OpenAPI.TOKEN = async () => {
  return localStorage.getItem("access_token") || ""
}

// Response interceptor: attempt refresh on 401 and retry original request
{
  let refreshing: Promise<boolean> | null = null

  OpenAPI.interceptors.response.use(async (response) => {
    // If not unauthorized, just return
    if (response.status !== 401) return response

    // If no refresh token, force logout
    const refreshToken = localStorage.getItem("refresh_token")
    if (!refreshToken) {
      localStorage.removeItem("access_token")
      localStorage.removeItem("refresh_token")
      window.location.href = "/login"
      return response
    }

    // If a refresh is already in progress, wait for it
    if (refreshing) {
      const ok = await refreshing
      if (!ok) {
        localStorage.removeItem("access_token")
        localStorage.removeItem("refresh_token")
        window.location.href = "/login"
        return response
      }
      // retry original request
      try {
        // Update Authorization header with the new access token before retrying
        const newToken = localStorage.getItem("access_token") || ""
        // @ts-ignore
        const retryConfig = { ...response.config, headers: { ...(response.config?.headers || {}), Authorization: newToken ? `Bearer ${newToken}` : undefined } }
        // @ts-ignore
        return await axios.request(retryConfig)
      } catch (e) {
        localStorage.removeItem("access_token")
        localStorage.removeItem("refresh_token")
        window.location.href = "/login"
        return response
      }
    }

    // start refresh
    refreshing = (async () => {
      try {
        const url = `${OpenAPI.BASE}/v1/login/refresh-token`
        const res = await fetch(url, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refreshToken }),
        })
        if (!res.ok) return false
        const body = await res.json()
        if (body.accessToken) {
          localStorage.setItem("access_token", body.accessToken)
          return true
        }
        return false
      } catch (e) {
        return false
      }
    })()

    const ok = await refreshing
    refreshing = null

    if (!ok) {
      localStorage.removeItem("access_token")
      localStorage.removeItem("refresh_token")
      window.location.href = "/login"
      return response
    }

    // retry original request
    try {
      const newToken = localStorage.getItem("access_token") || ""
      // @ts-ignore
      const retryConfig = { ...response.config, headers: { ...(response.config?.headers || {}), Authorization: newToken ? `Bearer ${newToken}` : undefined } }
      // @ts-ignore
      return await axios.request(retryConfig)
    } catch (e) {
      localStorage.removeItem("access_token")
      localStorage.removeItem("refresh_token")
      window.location.href = "/login"
      return response
    }
  })
}

const handleApiError = (error: Error) => {
  if (error instanceof ApiError && [401, 403].includes(error.status)) {
    localStorage.removeItem("access_token")
    window.location.href = "/login"
  }
}
const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: handleApiError,
  }),
  mutationCache: new MutationCache({
    onError: handleApiError,
  }),
})

const router = createRouter({ routeTree })
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <CustomProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </CustomProvider>
  </StrictMode>,
)
