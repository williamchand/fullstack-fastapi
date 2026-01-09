import {
  ArrowRight,
  CheckCircle2,
  CreditCard,
  Loader2,
  ShieldCheck,
  X,
} from "lucide-react"
import type React from "react"
import { useState } from "react"

interface StripeCheckoutModalProps {
  amount: number
  title: string
  onClose: () => void
  onSuccess: () => void
  purpose: "REGISTRATION" | "PUBLISHING"
}

const StripeCheckoutModal: React.FC<StripeCheckoutModalProps> = ({
  amount,
  title,
  onClose,
  onSuccess,
  purpose,
}) => {
  const [step, setStep] = useState<"details" | "processing" | "success">(
    "details",
  )

  const handlePay = (e: React.FormEvent) => {
    e.preventDefault()
    setStep("processing")
    setTimeout(() => {
      setStep("success")
      setTimeout(() => {
        onSuccess()
      }, 2000)
    }, 2500)
  }

  return (
    <div className="fixed inset-0 z-[110] flex items-center justify-center p-4 bg-gray-900/80 backdrop-blur-md animate-in fade-in">
      <div className="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden flex flex-col">
        <div className="p-8 border-b border-gray-50 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-[#635BFF] p-2 rounded-lg">
              <CreditCard className="w-5 h-5 text-white" />
            </div>
            <span className="font-black text-xl text-[#635BFF]">Stripe</span>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-full transition-colors"
          >
            <X className="w-6 h-6 text-gray-400" />
          </button>
        </div>

        <div className="p-10">
          {step === "details" && (
            <form onSubmit={handlePay} className="space-y-8">
              <div className="text-center space-y-2 mb-8">
                <p className="text-gray-500 font-bold uppercase tracking-widest text-xs">
                  {purpose === "REGISTRATION"
                    ? "Ticket Purchase"
                    : "Platform Fee"}
                </p>
                <h3 className="text-3xl font-black text-gray-900 leading-tight">
                  Pay ${amount.toFixed(2)}
                </h3>
                <p className="text-gray-400 text-sm font-medium">{title}</p>
              </div>

              <div className="space-y-6">
                <div className="space-y-2">
                  <label className="text-[10px] font-black text-gray-400 uppercase tracking-widest pl-1">
                    Card Information
                  </label>
                  <div className="relative">
                    <input
                      required
                      type="text"
                      placeholder="4242 4242 4242 4242"
                      className="w-full px-5 py-4 bg-gray-50 border-2 border-gray-100 rounded-2xl focus:border-[#635BFF] focus:ring-4 focus:ring-[#635BFF]/10 outline-none transition-all font-mono font-bold"
                    />
                    <div className="absolute right-4 top-1/2 -translate-y-1/2 flex gap-2">
                      <img
                        src="https://upload.wikimedia.org/wikipedia/commons/5/5e/Visa_Inc._logo.svg"
                        className="h-4"
                        alt="Visa"
                      />
                      <img
                        src="https://upload.wikimedia.org/wikipedia/commons/2/2a/Mastercard-logo.svg"
                        className="h-4"
                        alt="Mastercard"
                      />
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-[10px] font-black text-gray-400 uppercase tracking-widest pl-1">
                      Expiry
                    </label>
                    <input
                      required
                      type="text"
                      placeholder="MM / YY"
                      className="w-full px-5 py-4 bg-gray-50 border-2 border-gray-100 rounded-2xl focus:border-[#635BFF] focus:ring-4 focus:ring-[#635BFF]/10 outline-none transition-all font-bold"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-[10px] font-black text-gray-400 uppercase tracking-widest pl-1">
                      CVC
                    </label>
                    <input
                      required
                      type="text"
                      placeholder="123"
                      className="w-full px-5 py-4 bg-gray-50 border-2 border-gray-100 rounded-2xl focus:border-[#635BFF] focus:ring-4 focus:ring-[#635BFF]/10 outline-none transition-all font-bold"
                    />
                  </div>
                </div>
              </div>

              <div className="pt-6">
                <button
                  type="submit"
                  className="w-full py-5 bg-[#635BFF] text-white font-black rounded-2xl hover:bg-[#5249db] shadow-xl shadow-[#635BFF]/20 transition-all transform active:scale-95 flex items-center justify-center gap-2"
                >
                  Pay Now <ArrowRight className="w-5 h-5" />
                </button>
                <div className="flex items-center justify-center gap-2 mt-6 text-gray-400">
                  <ShieldCheck className="w-4 h-4 text-emerald-500" />
                  <span className="text-[10px] font-black uppercase tracking-widest">
                    Secure 256-bit Encrypted Payment
                  </span>
                </div>
              </div>
            </form>
          )}

          {step === "processing" && (
            <div className="py-20 text-center space-y-6 animate-in zoom-in-95">
              <div className="relative w-24 h-24 mx-auto">
                <div className="absolute inset-0 border-4 border-[#635BFF]/10 rounded-full" />
                <Loader2 className="w-full h-full text-[#635BFF] animate-spin" />
              </div>
              <div>
                <h4 className="text-xl font-black text-gray-900">
                  Processing Payment...
                </h4>
                <p className="text-gray-400 font-medium mt-1">
                  Please do not close your browser
                </p>
              </div>
            </div>
          )}

          {step === "success" && (
            <div className="py-20 text-center space-y-6 animate-in bounce-in">
              <div className="w-24 h-24 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto shadow-lg shadow-emerald-100">
                <CheckCircle2 className="w-12 h-12" />
              </div>
              <div>
                <h4 className="text-3xl font-black text-gray-900">
                  Payment Successful
                </h4>
                <p className="text-gray-500 font-medium mt-1">
                  {purpose === "REGISTRATION"
                    ? "You're all set! Check your learning shelf."
                    : "Webinar published successfully."}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default StripeCheckoutModal
