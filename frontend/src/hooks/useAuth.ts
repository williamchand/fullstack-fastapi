import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { useState } from "react"

import type { v1LoginUserRequest as AccessToken } from "@/client/user"
import type { v1User as UserPublic } from "@/client/user"
import type { ApiError } from "@/client/user"
import {
  userServiceCreateUser,
  userServiceGetUser,
  userServiceLoginUser,
} from "@/client/user"
import type { UserRegister } from "@/types/user"
import { getAuthErrorInfo, handleError } from "@/utils"

const isLoggedIn = () => {
  return localStorage.getItem("access_token") !== null
}

const useAuth = () => {
  const [error, setError] = useState<string | null>(null)
  const [authErrorInfo, setAuthErrorInfo] = useState<{
    code: string
    message: string
  } | null>(null)
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { data: user } = useQuery<UserPublic | null, Error>({
    queryKey: ["currentUser"],
    queryFn: async () => {
      const res = await userServiceGetUser()
      return res.user ?? null
    },
    enabled: isLoggedIn(),
  })

  const signUpMutation = useMutation({
    mutationFn: (data: UserRegister) =>
      userServiceCreateUser({
        requestBody: {
          email: data.email,
          fullName: data.fullName,
          password: data.password,
          roles: data.roles,
          isActive: data.is_active,
        },
      }),

    onSuccess: (_data, variables) => {
      navigate({ to: "/verify-email", search: { email: variables.email } })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] })
    },
  })

  const login = async (data: AccessToken) => {
    console.debug("useAuth.login called with:", data)
    const response = await userServiceLoginUser({ requestBody: data })
    if (response.accessToken) {
      localStorage.setItem("access_token", response.accessToken)
    }
    if (response.refreshToken) {
      localStorage.setItem("refresh_token", response.refreshToken)
    }
    if (response.refreshExpiresAt) {
      localStorage.setItem("refresh_expires_at", response.refreshExpiresAt)
    }
  }

  const loginMutation = useMutation({
    mutationFn: login,
    onError: (err: ApiError) => {
      const info = getAuthErrorInfo(err)
      setError(info.message)
      setAuthErrorInfo(info)
      handleError(err)
    },
  })

  const logout = () => {
    localStorage.removeItem("access_token")
    navigate({ to: "/login" })
  }

  return {
    signUpMutation,
    loginMutation,
    logout,
    user,
    error,
    authErrorInfo,
    resetError: () => {
      setError(null)
      setAuthErrorInfo(null)
    },
  }
}

export { isLoggedIn }
export default useAuth
