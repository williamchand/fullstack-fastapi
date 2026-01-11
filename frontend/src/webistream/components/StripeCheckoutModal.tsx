import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogHeader,
  DialogRoot,
} from "@/components/ui/dialog"
import { Box, Button, Flex, Input, Text } from "@chakra-ui/react"
import { ArrowRight, CheckCircle2, CreditCard, ShieldCheck } from "lucide-react"
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
    <DialogRoot open onOpenChange={(e) => !e.open && onClose()}>
      <DialogContent borderRadius="2xl" maxW="lg">
        <DialogHeader>
          <Flex align="center" justify="space-between">
            <Flex align="center" gap={3}>
              <Flex bg="#635BFF" p={2} borderRadius="lg">
                <CreditCard color="white" size={20} />
              </Flex>
              <Text fontWeight="black" fontSize="xl" color="#635BFF">
                Stripe
              </Text>
            </Flex>
          </Flex>
        </DialogHeader>
        <DialogCloseTrigger />
        <DialogBody>
          {step === "details" && (
            <Box
              as="form"
              onSubmit={handlePay}
              display="flex"
              flexDirection="column"
              gap={6}
            >
              <Box textAlign="center" mb={6}>
                <Text
                  color="gray.500"
                  fontWeight="bold"
                  fontSize="xs"
                  textTransform="uppercase"
                  letterSpacing="widest"
                >
                  {purpose === "REGISTRATION"
                    ? "Ticket Purchase"
                    : "Platform Fee"}
                </Text>
                <Text fontSize="3xl" fontWeight="black" color="gray.900">
                  Pay ${amount.toFixed(2)}
                </Text>
                <Text color="gray.400" fontSize="sm" fontWeight="medium">
                  {title}
                </Text>
              </Box>
              <Box display="flex" flexDirection="column" gap={4}>
                <Box>
                  <Text
                    fontSize="xs"
                    fontWeight="black"
                    color="gray.400"
                    letterSpacing="widest"
                    textTransform="uppercase"
                    pl={1}
                  >
                    Card Information
                  </Text>
                  <Flex position="relative" align="center">
                    <Input required placeholder="4242 4242 4242 4242" />
                    <Flex
                      position="absolute"
                      right={4}
                      top="50%"
                      transform="translateY(-50%)"
                      gap={2}
                    >
                      <img
                        src="https://upload.wikimedia.org/wikipedia/commons/5/5e/Visa_Inc._logo.svg"
                        height={16}
                        alt="Visa"
                      />
                      <img
                        src="https://upload.wikimedia.org/wikipedia/commons/2/2a/Mastercard-logo.svg"
                        height={16}
                        alt="Mastercard"
                      />
                    </Flex>
                  </Flex>
                </Box>
                <Flex gap={4}>
                  <Box flex="1">
                    <Text
                      fontSize="xs"
                      fontWeight="black"
                      color="gray.400"
                      letterSpacing="widest"
                      textTransform="uppercase"
                      pl={1}
                    >
                      Expiry
                    </Text>
                    <Input required placeholder="MM / YY" />
                  </Box>
                  <Box flex="1">
                    <Text
                      fontSize="xs"
                      fontWeight="black"
                      color="gray.400"
                      letterSpacing="widest"
                      textTransform="uppercase"
                      pl={1}
                    >
                      CVC
                    </Text>
                    <Input required placeholder="123" />
                  </Box>
                </Flex>
              </Box>
              <Box pt={4}>
                <Button
                  type="submit"
                  w="full"
                  py={5}
                  variant="solid"
                  colorPalette="blue"
                  borderRadius="2xl"
                  fontWeight="black"
                  gap={2}
                >
                  Pay Now <ArrowRight size={20} />
                </Button>
                <Flex
                  align="center"
                  justify="center"
                  gap={2}
                  mt={6}
                  color="gray.400"
                >
                  <ShieldCheck size={16} color="#10b981" />
                  <Text
                    fontSize="xs"
                    fontWeight="black"
                    textTransform="uppercase"
                    letterSpacing="widest"
                  >
                    Secure 256-bit Encrypted Payment
                  </Text>
                </Flex>
              </Box>
            </Box>
          )}

          {step === "processing" && (
            <Box py={20} textAlign="center">
              <Text fontSize="xl" fontWeight="black" color="gray.900">
                Processing Payment...
              </Text>
              <Text color="gray.400" fontWeight="medium" mt={1}>
                Please do not close your browser
              </Text>
            </Box>
          )}

          {step === "success" && (
            <Box
              py={16}
              textAlign="center"
              display="flex"
              flexDirection="column"
              gap={6}
              alignItems="center"
            >
              <Flex
                w={24}
                h={24}
                bg="emerald.100"
                color="emerald.600"
                borderRadius="full"
                align="center"
                justify="center"
                boxShadow="lg"
              >
                <CheckCircle2 size={48} />
              </Flex>
              <Box>
                <Text fontSize="3xl" fontWeight="black" color="gray.900">
                  Payment Successful
                </Text>
                <Text color="gray.500" fontWeight="medium" mt={1}>
                  {purpose === "REGISTRATION"
                    ? "You're all set! Check your learning shelf."
                    : "Webinar published successfully."}
                </Text>
              </Box>
            </Box>
          )}
        </DialogBody>
      </DialogContent>
    </DialogRoot>
  )
}

export default StripeCheckoutModal
