import { create } from "zustand"
import { persist } from "zustand/middleware"

type VerifyPhone = {
  number: string
  region: string
  token: string
}

interface UIState {
  method: "email" | "phone"
  verifyPhone: VerifyPhone
  selectedRole: "customer" | "salon_owner" | null
  setMethod: (method: "email" | "phone") => void
  setVerifyPhone: (data: Partial<VerifyPhone>) => void
  setSelectedRole: (role: "customer" | "salon_owner" | null) => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      method: "email",
      verifyPhone: {
        number: "",
        region: "ID",
        token: "",
      },
      selectedRole: null,

      setMethod: (method) => set({ method }),

      setVerifyPhone: (data) =>
        set((state) => ({
          verifyPhone: {
            ...state.verifyPhone,
            ...data,
            ...(data.number ? { number: data.number.replace(/\D/g, "") } : {}),
          },
        })),
      setSelectedRole: (role) => set({ selectedRole: role }),
    }),
    {
      name: "ui-storage",
      partialize: (state) => ({
        method: state.method,
        verifyPhone: state.verifyPhone,
        selectedRole: state.selectedRole,
      }),
    },
  ),
)
