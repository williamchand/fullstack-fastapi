import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type VerifyPhone = {
  number: string
  region: string
  token: string
}

interface UIState {
  method: 'email' | 'phone'
  verifyPhone: VerifyPhone
  setMethod: (method: 'email' | 'phone') => void
  setVerifyPhone: (data: Partial<VerifyPhone>) => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      method: 'email',
      verifyPhone: {
        number: '',
        region: 'ID',
        token: '',
      },

      setMethod: (method) => set({ method }),

      setVerifyPhone: (data) =>
        set((state) => ({
          verifyPhone: {
            ...state.verifyPhone,
            ...data,
            ...(data.number
              ? { number: data.number.replace(/\D/g, '') }
              : {}),
          },
        })),
    }),
    {
      name: 'ui-storage',
      partialize: (state) => ({
        method: state.method,
        verifyPhone: state.verifyPhone,
      }),
    }
  )
)
