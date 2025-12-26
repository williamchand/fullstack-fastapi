import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface UIState {
  method: 'email' | 'phone'
  verifyPhoneNumber: string
  verifyRegion: string
  setMethod: (method: 'email' | 'phone') => void
  setVerifyData: (phoneNumber: string, region: string) => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      method: 'email',
      verifyPhoneNumber: '',
      verifyRegion: 'ID',
      setMethod: (method) => set({ method }),
      setVerifyData: (phoneNumber, region) => set({ verifyPhoneNumber: phoneNumber.replace(/\D/g, ''), verifyRegion: region }),
    }),
    {
      name: 'ui-storage',
      partialize: (state) => ({ method: state.method, verifyPhoneNumber: state.verifyPhoneNumber, verifyRegion: state.verifyRegion }),
    }
  )
)