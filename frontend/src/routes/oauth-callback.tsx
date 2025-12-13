import { createFileRoute, useNavigate, useSearch } from "@tanstack/react-router"
import { oauthServiceHandleOauthCallback } from "@/client/oauth"

export const Route = createFileRoute("/oauth-callback")({
  component: OAuthCallback,
})

function OAuthCallback() {
  const navigate = useNavigate()
  const search = useSearch({ from: "/oauth-callback" }) as {
    code?: string
    state?: string
    provider?: string
  }

  if (window.opener && search.code) {
    try {
      window.opener.postMessage(
        {
          type: "oauth-callback",
          provider: search.provider || "google",
          code: search.code,
          state: search.state,
        },
        window.location.origin,
      )
    } catch {}
    window.close()
  }

  // Mobile/full-page fallback: if opened without an opener, exchange code here
  if (!window.opener && search.code) {
    ;(async () => {
      try {
        const processedKey = search.state
          ? `oauth_processed:${search.state}`
          : null
        try {
          if (processedKey && sessionStorage.getItem(processedKey)) {
            let target = "/"
            try {
              const k = search.state
                ? `oauth_redirect_tab:${search.state}`
                : "oauth_redirect_tab"
              target = sessionStorage.getItem(k) || "/"
            } catch {}
            navigate({ to: target })
            return
          }
        } catch {}
        const res = await oauthServiceHandleOauthCallback({
          provider: search.provider || "google",
          requestBody: { code: search.code },
        })
        if (res.accessToken) {
          localStorage.setItem("access_token", res.accessToken)
        }
        if (res.refreshToken) {
          localStorage.setItem("refresh_token", res.refreshToken)
        }
        if (res.refreshExpiresAt) {
          localStorage.setItem("refresh_expires_at", res.refreshExpiresAt)
        }
        let target = "/"
        try {
          const k = search.state
            ? `oauth_redirect_tab:${search.state}`
            : "oauth_redirect_tab"
          target = sessionStorage.getItem(k) || "/"
          sessionStorage.removeItem(k)
          if (processedKey) sessionStorage.setItem(processedKey, "1")
        } catch {}
        navigate({ to: target })
      } catch {
        // minimal fallback: return to login so user can retry
        try {
          const k = search.state
            ? `oauth_redirect_tab:${search.state}`
            : "oauth_redirect_tab"
          sessionStorage.removeItem(k)
        } catch {}
        navigate({ to: "/login" })
      }
    })()
  }

  return null
}

export default OAuthCallback
