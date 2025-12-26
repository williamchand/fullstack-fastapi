import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface UIState {
  method: 'email' | 'phone'
  setMethod: (method: 'email' | 'phone') => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      method: 'email',
      setMethod: (method) => set({ method }),
    }),
    {
      name: 'ui-storage',
      partialize: (state) => ({ method: state.method }),
    }
  )
)