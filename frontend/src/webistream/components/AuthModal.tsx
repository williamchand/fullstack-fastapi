import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogHeader,
  DialogRoot,
} from "@/components/ui/dialog"
import { UserRole, type User as UserType } from "@/webistream/types"
import {
  Box,
  Button,
  Flex,
  HStack,
  Input,
  Text,
  VStack,
} from "@chakra-ui/react"
import { ArrowRight, Briefcase, ShieldCheck, User } from "lucide-react"
import type React from "react"
import { useState } from "react"

interface AuthModalProps {
  onClose: () => void
  onAuthSuccess: (user: UserType) => void
  initialRole?: UserRole
}

const AuthModal: React.FC<AuthModalProps> = ({
  onClose,
  onAuthSuccess,
  initialRole = UserRole.CUSTOMER,
}) => {
  const [isLogin, setIsLogin] = useState(true)
  const [role, setRole] = useState<UserRole>(initialRole)
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setTimeout(() => {
      const mockUser: UserType = {
        id: Math.random().toString(36).substr(2, 9),
        name: formData.name || (isLogin ? "Alex Rivers" : "New User"),
        email: formData.email,
        role:
          isLogin && formData.email.includes("host") ? UserRole.BUSINESS : role,
        avatar: `https://picsum.photos/seed/${formData.email}/40/40`,
        isStripeConnected: role === UserRole.BUSINESS,
      }
      onAuthSuccess(mockUser)
      setLoading(false)
    }, 1000)
  }

  return (
    <DialogRoot open onOpenChange={(e) => !e.open && onClose()}>
      <DialogContent borderRadius="2xl">
        <DialogHeader>
          <VStack align="start" gap={1}>
            <Text fontSize="xl" fontWeight="bold">
              {isLogin ? "Sign In" : "Join WebiStream"}
            </Text>
            <Text fontSize="sm" color="gray.500">
              Access your{" "}
              {role === UserRole.BUSINESS
                ? "Business Dashboard"
                : "Learning Shelf"}
            </Text>
          </VStack>
        </DialogHeader>
        <DialogCloseTrigger />
        <DialogBody>
          <VStack as="form" onSubmit={handleSubmit} gap={5} align="stretch">
            <HStack gap={2}>
              <Button
                variant={role === UserRole.CUSTOMER ? "solid" : "outline"}
                colorPalette="blue"
                onClick={() => setRole(UserRole.CUSTOMER)}
                flex="1"
              >
                <User size={16} />
                Customer
              </Button>
              <Button
                variant={role === UserRole.BUSINESS ? "solid" : "outline"}
                colorPalette="blue"
                onClick={() => setRole(UserRole.BUSINESS)}
                flex="1"
              >
                <Briefcase size={16} />
                Business
              </Button>
            </HStack>

            {role === UserRole.BUSINESS && (
              <Flex
                gap={3}
                p={3}
                borderRadius="xl"
                bg="blue.50"
                border="1px solid"
                borderColor="blue.100"
                align="start"
              >
                <ShieldCheck size={18} color="#4f46e5" />
                <Text fontSize="xs" color="blue.700">
                  Hosting an event? Sign in to manage your Stripe payments, Zoom
                  links, and attendee lists.
                </Text>
              </Flex>
            )}

            {!isLogin && (
              <Box>
                <Text fontSize="xs" color="gray.500" fontWeight="bold">
                  Full Name
                </Text>
                <Input
                  required
                  placeholder="John Doe"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                />
              </Box>
            )}

            <Box>
              <Text fontSize="xs" color="gray.500" fontWeight="bold">
                Email Address
              </Text>
              <Input
                required
                type="email"
                placeholder={
                  role === UserRole.BUSINESS
                    ? "host@company.com"
                    : "hello@example.com"
                }
                value={formData.email}
                onChange={(e) =>
                  setFormData({ ...formData, email: e.target.value })
                }
              />
            </Box>

            <Box>
              <Text fontSize="xs" color="gray.500" fontWeight="bold">
                Password
              </Text>
              <Input
                required
                type="password"
                placeholder="••••••••"
                value={formData.password}
                onChange={(e) =>
                  setFormData({ ...formData, password: e.target.value })
                }
              />
            </Box>

            <Button
              type="submit"
              disabled={loading}
              loading={loading}
              colorPalette={role === UserRole.BUSINESS ? "gray" : "blue"}
            >
              {!loading && <ArrowRight size={16} />}
              {loading
                ? "Processing..."
                : isLogin
                  ? `Sign In as ${role === UserRole.BUSINESS ? "Host" : "User"}`
                  : "Create Account"}
            </Button>

            <Button
              type="button"
              variant="ghost"
              onClick={() => setIsLogin(!isLogin)}
            >
              {isLogin
                ? "Don't have an account? Sign up"
                : "Already have an account? Sign in"}
            </Button>
          </VStack>
        </DialogBody>
      </DialogContent>
    </DialogRoot>
  )
}

export default AuthModal
